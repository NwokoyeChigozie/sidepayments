package payment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func PaymentAccountMonnifyListService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.PaymentAccountMonnifyListRequest) (models.PaymentAccount, int, error) {
	var (
		data               models.PaymentAccount
		generatedReference = ""
		configData         = config.GetConfig()
		monnify            = Monnify{ExtReq: extReq}
		paymentChannelD    = config.GetConfig().Slack.PaymentChannelID
		paymentAccount     = models.PaymentAccount{BusinessID: strconv.Itoa(int(req.AccountID))}
		amount             float64
		charge             float64
	)

	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	businessProfile, err := GetBusinessProfileByAccountID(extReq, extReq.Logger, req.AccountID)
	if err != nil {
		return data, http.StatusInternalServerError, fmt.Errorf("business profile not found for user: %v", err)
	}

	accountName := businessProfile.BusinessName
	if accountName == "" {
		accountName = user.Firstname
		if accountName == "" {
			accountName = strconv.Itoa(int(user.AccountID))
		}
	}
	accountEmail := user.EmailAddress

	if req.TransactionID != "" {
		transaction, err := ListTransactionsByID(extReq, req.TransactionID)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}
		amount, charge = transaction.Amount, transaction.EscrowCharge
		paymentAccount = models.PaymentAccount{BusinessID: strconv.Itoa(int(user.AccountID)), TransactionID: req.TransactionID}
		code, err := paymentAccount.GetPaymentAccountByBusinessIDAndTransactionID(db.Payment)
		if err != nil && code == http.StatusInternalServerError {
			return data, code, err
		}
	} else {
		paymentAccount = models.PaymentAccount{BusinessID: strconv.Itoa(int(user.AccountID))}
		code, err := paymentAccount.GetPaymentAccountByBusinessID(db.Payment)
		if err != nil && code == http.StatusInternalServerError {
			return data, code, err
		}
	}

	if paymentAccount.ID == 0 {
		if req.GeneratedReference != "" {
			generatedReference = req.GeneratedReference
		} else {
			uuID, _ := uuid.NewV4()
			generatedReference = "VESICASH_VA_" + uuID.String()
		}
		currencyCode := "NGN"

		paymentInfo := models.PaymentInfo{Reference: generatedReference}
		code, err := paymentInfo.GetPaymentInfoByReference(db.Payment)
		if err != nil {
			if code == http.StatusInternalServerError {
				return data, code, err
			}

			payment := models.Payment{
				PaymentID:     utility.RandomString(10),
				TotalAmount:   amount,
				EscrowCharge:  charge,
				IsPaid:        false,
				AccountID:     int64(req.AccountID),
				BusinessID:    int64(businessProfile.AccountID),
				TransactionID: req.TransactionID,
				Currency:      currencyCode,
			}

			err = payment.CreatePayment(db.Payment)
			if err != nil {
				return data, http.StatusInternalServerError, err
			}

			paymentInfo = models.PaymentInfo{
				PaymentID: payment.PaymentID,
				Reference: generatedReference,
				Status:    "pending",
				Gateway:   req.Gateway,
			}
			err = paymentInfo.CreatePaymentInfo(db.Payment)
			if err != nil {
				return data, http.StatusInternalServerError, err
			}
		}

		paymentAccount := models.PaymentAccount{
			PaymentAccountID: generatedReference,
			TransactionID:    req.TransactionID,
			PaymentID:        paymentInfo.PaymentID,
		}

		if req.Gateway == "rave" {
			paymentAccount.AccountNumber = configData.Rave.MerchantId
			paymentAccount.AccountName = configData.Rave.AccountName
			paymentAccount.BankCode = "flutterwave"
			paymentAccount.BankName = "flutterwave"
			paymentAccount.Status = "ACTIVE"
		} else {
			accountDetails, err := monnify.ReserveAccount(generatedReference, accountName, currencyCode, accountEmail)
			if err != nil {
				return data, http.StatusInternalServerError, err
			}
			paymentAccount.AccountNumber = accountDetails.AccountNumber
			paymentAccount.AccountName = accountDetails.AccountName
			paymentAccount.BankCode = accountDetails.BankCode
			paymentAccount.BankName = accountDetails.BankName
			paymentAccount.ReservationReference = accountDetails.ReservationReference
			paymentAccount.Status = accountDetails.Status
		}

		paymentAccount.IsUsed = true
		paymentAccount.ExpiresAfter = strconv.Itoa(int(time.Now().Add(72 * time.Hour).Unix()))
		paymentAccount.BusinessID = strconv.Itoa(req.AccountID)
		err = paymentAccount.CreatePaymentAccount(db.Payment)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}
		err = SlackNotify(paymentChannelD, `
					Virtual Bank Account Generated
                    Environment: `+config.GetConfig().App.Name+`
                    Account Number: `+paymentAccount.AccountNumber+`
                    Account Name: `+paymentAccount.AccountName+`
                    Bank: `+paymentAccount.BankName+`
                    Status: SUCCESSFUL
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
	} else {
		paymentAccount.TransactionID = req.TransactionID
		err := paymentAccount.UpdateAllFields(db.Payment)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}
	}

	pendingTransferFunding := models.PendingTransferFunding{
		Reference: paymentAccount.PaymentAccountID,
		Status:    "pending",
		Type:      "walletfunding",
	}

	err = pendingTransferFunding.CreatePendingTransferFunding(db.Payment)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return paymentAccount, http.StatusOK, nil
}
func PaymentAccountMonnifyVerifyService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.PaymentAccountMonnifyVerifyRequest) (map[string]interface{}, string, int, error) {
	var (
		data            = map[string]interface{}{"reference": req.Reference, "amount": 0.00, "pdf_link": "", "status": false}
		msg             = ""
		monnify         = Monnify{ExtReq: extReq}
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
		paymentAccount  = models.PaymentAccount{PaymentAccountID: req.Reference}
		transaction     external_models.TransactionByID
	)

	code, err := paymentAccount.GetPaymentAccountByPaymentAccountID(db.Payment)
	if err != nil {
		return data, msg, code, err
	}

	payment := models.Payment{PaymentID: paymentAccount.PaymentID}
	code, err = payment.GetPaymentByPaymentID(db.Payment)
	if err != nil && code == http.StatusInternalServerError {
		return data, msg, code, err
	}

	trans, err := monnify.FetchAccountTrans(req.Reference)
	if err != nil {
		return data, msg, http.StatusInternalServerError, err
	}

	if len(trans) < 1 {
		return data, msg, http.StatusBadRequest, fmt.Errorf("account has not received any payments")
	}

	if trans[0].PaymentReference == paymentAccount.PaymentReference {
		return data, msg, http.StatusBadRequest, fmt.Errorf("account has received this payment")
	}

	if paymentAccount.TransactionID != "" && paymentAccount.PaymentID != "" && paymentAccount.PaymentID != "0" {
		payment.TransactionID = paymentAccount.TransactionID
		paymentAmount, charge := payment.TotalAmount, payment.EscrowCharge
		if paymentAmount > charge {
			paymentAmount = paymentAmount - charge
		} else {
			paymentAmount = 0
		}

		if payment.IsPaid {
			data["amount"] = paymentAmount
			return data, msg, http.StatusOK, nil
		}

		verify, amountPaid, err := monnify.VerifyTrans(req.Reference, paymentAmount)
		if err != nil {
			return data, msg, http.StatusBadRequest, err
		}

		if verify {
			transaction, pdfLink, err := sendTransactionConfirmed(extReq, db, &payment, req.Reference, amountPaid)
			if err != nil {
				return data, msg, http.StatusInternalServerError, err
			}
			chargeBearer := transaction.Parties["charge_bearer"]
			paymentAccount.Delete(db.Payment)

			err = SlackNotify(paymentChannelD, `
					Bank Transfer Payment
					Environment: `+config.GetConfig().App.Name+`
					Reference: `+trans[0].TransactionReference+`
					Transaction ID: `+paymentAccount.TransactionID+`
					Payment ID: `+payment.PaymentID+`
					Account Number: `+paymentAccount.AccountNumber+`
					Account Name: `+paymentAccount.AccountName+`
					Bank: `+paymentAccount.BankName+`
					Amount: `+fmt.Sprintf("%v", payment.TotalAmount)+`
					Escrow Charge: `+fmt.Sprintf("%v", payment.EscrowCharge)+`
					Escrow Charge Bearer: `+fmt.Sprintf("%v", chargeBearer.AccountID)+`
					Status: SUCCESSFUL
			`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}

			businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, transaction.BusinessID)
			if businessProfileData.Webhook_uri != "" {
				InitWebhook(extReq, db, businessProfileData.Webhook_uri, "payment.success", map[string]interface{}{
					"reference": req.Reference,
					"amount":    payment.TotalAmount,
					"status":    "success",
				}, businessProfileData.AccountID)
			}

			return map[string]interface{}{"reference": req.Reference, "amount": payment.TotalAmount, "pdf_link": pdfLink, "status": verify}, "Bank Transfer Verified", http.StatusOK, nil
		} else {
			return map[string]interface{}{"reference": req.Reference, "amount": payment.TotalAmount, "pdf_link": "", "status": verify}, "Bank Transfer Not Verified", http.StatusOK, nil
		}

	}

	fundEscrowWallet := "no"
	if req.TransactionID != "" {
		transaction, err = ListTransactionsByID(extReq, payment.TransactionID)
		if err != nil {
			return data, msg, http.StatusInternalServerError, err
		}
		if transaction.EscrowWallet != "" {
			fundEscrowWallet = transaction.EscrowWallet
		}
	}

	verify, amountPaid, err := monnify.VerifyTrans(req.Reference, trans[0].AmountPaid)
	if err != nil {
		return data, msg, http.StatusBadRequest, err
	}

	if verify {
		paymentAccountBusinessID, _ := strconv.Atoi(paymentAccount.BusinessID)
		user, err := GetUserWithAccountID(extReq, paymentAccountBusinessID)
		if err != nil {
			return data, msg, http.StatusInternalServerError, fmt.Errorf("user with businessID %v not found : %v", paymentAccountBusinessID, err.Error())
		}
		pendingTransferFunding := models.PendingTransferFunding{Reference: req.Reference}
		pendingTransferFunding.GetPendingTransferFundingByReference(db.Payment)
		pendingTransferFunding.Delete(db.Payment)

		paymentAccount.PaymentReference = req.Reference
		payment.UpdateAllFields(db.Payment)

		var (
			fundingCharge float64 = 500
			firstLimit    float64 = 500000
			secondLimit   float64 = 1000000
		)

		if amountPaid <= firstLimit {
			fundingCharge = 500
		} else if amountPaid > firstLimit && amountPaid <= secondLimit {
			fundingCharge = 1000
		} else {
			fundingCharge = 2000
		}
		finalAmount := amountPaid - fundingCharge

		err = CreditWallet(extReq, db, finalAmount, transaction.Currency, int(user.AccountID), false, fundEscrowWallet, transaction.TransactionID)
		if err != nil {
			return data, msg, http.StatusBadRequest, err
		}

		walletFunded := strings.ToUpper(transaction.Currency)
		if walletFunded == "" {
			walletFunded = "NGN"
		}

		if strings.ToLower(fundEscrowWallet) == "yes" {
			walletFunded = "ESCROW_" + walletFunded
		}
		if payment.ID == 0 {
			payment = models.Payment{
				PaymentID:      utility.RandomString(10),
				TotalAmount:    amountPaid,
				EscrowCharge:   fundingCharge,
				IsPaid:         true,
				AccountID:      int64(user.AccountID),
				BusinessID:     int64(user.BusinessId),
				Currency:       "NGN",
				Payment_method: "bank_transfer",
				WalletFunded:   walletFunded,
			}
			err := payment.CreatePayment(db.Payment)
			if err != nil {
				return data, msg, http.StatusBadRequest, err
			}
		} else {
			payment.TotalAmount = amountPaid
			payment.EscrowCharge = fundingCharge
			payment.IsPaid = true
			payment.Payment_method = "bank_transfer"
			payment.WalletFunded = walletFunded
			err := payment.UpdateAllFields(db.Payment)
			if err != nil {
				return data, msg, http.StatusBadRequest, err
			}
		}

		err = SlackNotify(paymentChannelD, `

					Bank Transfer Payment
					Environment: `+config.GetConfig().App.Name+`
					Reference: `+trans[0].TransactionReference+`
					Transaction ID: `+paymentAccount.TransactionID+`
					Payment ID: `+payment.PaymentID+`
					Account Number: `+paymentAccount.AccountNumber+`
					Account Name: `+paymentAccount.AccountName+`
					Bank: `+paymentAccount.BankName+`
					Amount: `+fmt.Sprintf("%v", payment.TotalAmount)+`
					Escrow Charge: `+fmt.Sprintf("%v", payment.EscrowCharge)+`
					Status: SUCCESSFUL
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}

		businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, transaction.BusinessID)
		if businessProfileData.Webhook_uri != "" {
			InitWebhook(extReq, db, businessProfileData.Webhook_uri, "payment.success", map[string]interface{}{
				"reference": req.Reference,
				"amount":    payment.TotalAmount,
				"status":    "success",
			}, businessProfileData.AccountID)
		}

		return map[string]interface{}{"reference": req.Reference, "amount": payment.TotalAmount, "pdf_link": "", "status": verify}, "Transfer Verified", http.StatusOK, nil
	}

	return map[string]interface{}{"reference": req.Reference, "amount": payment.TotalAmount, "pdf_link": "", "status": verify}, "Bank Transfer Not Verified", http.StatusOK, nil
}

func sendTransactionConfirmed(extReq request.ExternalRequest, db postgresql.Databases, payment *models.Payment, reference string, amountPaid float64) (external_models.TransactionByID, string, error) {
	var (
		amount = payment.TotalAmount
	)
	if amount <= 0 && amountPaid > 0 {
		amount = amountPaid
	}
	paymentInfo := models.PaymentInfo{PaymentID: payment.PaymentID}
	code, err := paymentInfo.GetPaymentInfoByPaymentID(db.Payment)
	if err != nil && code == http.StatusInternalServerError {
		return external_models.TransactionByID{}, "", err
	}

	transaction, err := ListTransactionsByID(extReq, payment.TransactionID)
	if err != nil {
		return transaction, "", err
	}
	extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
		AccountID:     transaction.BusinessID,
		TransactionID: transaction.TransactionID,
		MilestoneID:   transaction.MilestoneID,
		Status:        "ip",
	})
	if err != nil {
		extReq.Logger.Error(err.Error())
	}
	chargeBearer := transaction.Parties["charge_bearer"]
	seller, ok := transaction.Parties["seller"]
	if !ok {
		return transaction, "", fmt.Errorf("seller not found for transaction")
	}

	if chargeBearer.AccountID == seller.AccountID {
		amount = amount - payment.EscrowCharge
		if amount < 0 {
			amount = 0
		}

		payment.TotalAmount = amount
		err = payment.UpdateAllFields(db.Payment)
		if err != nil {
			return transaction, "", err
		}
	}

	businessEscrowCharge, err := getBusinessChargeWithBusinessIDAndCurrency(extReq, transaction.BusinessID, transaction.Currency)

	if err == nil {
		businessPerc, _ := strconv.ParseFloat(businessEscrowCharge.BusinessCharge, 64)
		vesicashCharge, _ := strconv.ParseFloat(businessEscrowCharge.VesicashCharge, 64)
		// credit vesicash
		err = CreditWallet(extReq, db, utility.PercentageOf(amount, vesicashCharge), transaction.Currency, 1, false, "no", transaction.TransactionID)
		if err != nil {
			return transaction, "", err
		}

		err = CreditWallet(extReq, db, utility.PercentageOf(amount, businessPerc), transaction.Currency, transaction.BusinessID, false, transaction.EscrowWallet, transaction.TransactionID)
		if err != nil {
			return transaction, "", err
		}
	}

	payment.TotalAmount = amount
	payment.IsPaid = true
	payment.PaymentMadeAt = time.Now()
	payment.PaymentType = "transaction"
	err = payment.UpdateAllFields(db.Payment)
	if err != nil {
		return transaction, "", err
	}

	if paymentInfo.ID != 0 {
		paymentInfo.Status = "paid"
		err = paymentInfo.UpdateAllFields(db.Payment)
		if err != nil {
			return transaction, "", err
		}
	}
	// TODO: generate pdflink
	pdfData := NewPdfData(extReq, transaction, *payment, reference, "")
	pdflink, err := GetPdfLink(extReq, "./templates/invoice_pdf.html", pdfData)
	if err != nil {
		extReq.Logger.Error("error generating pdf ", err.Error())
	}
	return transaction, pdflink, nil
}
