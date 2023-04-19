package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func MonnifyWebhookService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.MonnifyWebhookRequest, requestBody []byte) (int, error) {
	var (
		secret                   = config.GetConfig().Monnify.MonnifySecret
		data                     models.MonnifyWebhookRequestEventData
		monnifySignature         = utility.GetHeader(c, "monnify-signature")
		transactionReference     string
		paymentReference         string
		generatedReference       string
		amountPaid               float64
		customerEmail            string
		paidOn                   time.Time
		currency                 string
		accountNumber            string
		paymentSourceInformation string
		bankCode                 string
		bankName                 string
		paymentChannelD          = config.GetConfig().Slack.PaymentChannelID
	)

	hash := utility.Sha512Hmac(secret, requestBody)
	if hash != monnifySignature {
		extReq.Logger.Error("monnify webhhook log error", "Web Hook Denied, Hash Mismatch", hash, monnifySignature, requestBody)
		return http.StatusUnauthorized, fmt.Errorf("web Hook Denied, Hash Mismatch")
	}

	extReq.Logger.Info("monnify webhhook log info", string(requestBody))
	webhookLog := models.WebhookLog{
		Log:      string(requestBody),
		Provider: "monnify",
	}
	err := webhookLog.CreateWebhookLog(db.Payment)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if req.EventType == "" {
		extReq.Logger.Error("monnify webhhook log error", "no event type specified")
		return http.StatusBadRequest, fmt.Errorf("no event type specified")
	}

	if req.EventData != nil {
		data = *req.EventData
	} else {
		extReq.Logger.Error("monnify webhhook log error", "event data not found")
		return http.StatusBadRequest, fmt.Errorf("event data not found")
	}

	if data.PaymentReference != nil {
		paymentReference = *data.PaymentReference
	}

	if data.TransactionReference != nil {
		transactionReference = *data.TransactionReference
	}

	if data.Product != nil {
		if data.Product.Reference != nil {
			generatedReference = *data.Product.Reference
		}
	}

	if data.AmountPaid != nil {
		amountPaid = *data.AmountPaid
	}
	if data.PaidOn != nil {
		dateTime, err := time.Parse("2006-01-02 15:04:05.000", *data.PaidOn)
		if err != nil {
			extReq.Logger.Error("monnify webhhook log error", "could not parse paid on", data.PaidOn)
		}
		paidOn = dateTime
	}
	if data.Customer != nil {
		if data.Customer.Email != nil {
			customerEmail = *data.Customer.Email
		}
	}
	if data.Currency != nil {
		currency = strings.ToUpper(*data.Currency)
	}

	if data.DestinationAccountInformation != nil {
		if data.DestinationAccountInformation.AccountNumber != nil {
			accountNumber = *data.DestinationAccountInformation.AccountNumber
		}
		if data.DestinationAccountInformation.BankCode != nil {
			bankCode = *data.DestinationAccountInformation.BankCode
		}
		if data.DestinationAccountInformation.BankName != nil {
			bankName = *data.DestinationAccountInformation.BankName
		}

	}

	if data.PaymentSourceInformation != nil {
		if len(*data.PaymentSourceInformation) > 0 {
			arr := *data.PaymentSourceInformation
			paymentSourceInformation = *arr[0].AccountName
		}
	}

	err = SlackNotify(extReq, paymentChannelD, `
					Web Hook Received Monnify
					Environment: `+config.GetConfig().App.Name+`
					Event: `+req.EventType+`
					Reference: `+fmt.Sprintf("transaction reference:%v, generated reference: %v, payment reference:%v", transactionReference, generatedReference, paymentReference)+`
					Status: SUCCESSFUL
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	go func(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, transactionReference, generatedReference, paymentReference, paymentSourceInformation, customerEmail, accountNumber, bankCode, bankName, currency string, amountPaid float64, paidOn time.Time) {
		code, err := handleMonnifyWebhookRequest(c, extReq, db, transactionReference, generatedReference, paymentReference, paymentSourceInformation, customerEmail, accountNumber, bankCode, bankName, currency, amountPaid, paidOn)
		if err != nil {
			extReq.Logger.Error("monnify webhhook log error:", err.Error(), "code:", code)
		}
	}(c, extReq, db, transactionReference, generatedReference, paymentReference, paymentSourceInformation, customerEmail, accountNumber, bankCode, bankName, currency, amountPaid, paidOn)

	return http.StatusOK, nil
}

func handleMonnifyWebhookRequest(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, transactionReference, generatedReference, paymentReference, paymentSourceInformation, customerEmail, accountNumber, bankCode, bankName, currency string, amountPaid float64, paidOn time.Time) (int, error) {
	var (
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
		monnify         = Monnify{ExtReq: extReq}
	)
	paymentAccount := models.PaymentAccount{PaymentAccountID: generatedReference}
	code, err := paymentAccount.GetByPaymentAccountIDAndTransactionIDNotNull(db.Payment)
	if code == http.StatusInternalServerError {
		extReq.Logger.Error("monnify webhhook log error", err.Error())
		return code, err
	}

	if err != nil {
		user, err := GetUserWithEmail(extReq, customerEmail)
		if err != nil {
			extReq.Logger.Error("monnify webhhook log error", "user not found", err.Error())
			return http.StatusOK, fmt.Errorf("user not found")
		}

		payment := models.Payment{
			PaymentID:    utility.RandomString(10),
			TotalAmount:  amountPaid,
			EscrowCharge: 0,
			IsPaid:       false,
			AccountID:    int64(user.AccountID),
			BusinessID:   int64(user.AccountID),
			Currency:     currency,
		}

		err = payment.CreatePayment(db.Payment)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		paymentAccount = models.PaymentAccount{
			PaymentAccountID: generatedReference,
			PaymentID:        payment.PaymentID,
			AccountNumber:    accountNumber,
			BankCode:         bankCode,
			BankName:         bankName,
			Status:           "ACTIVE",
			IsUsed:           true,
			ExpiresAfter:     strconv.Itoa(int(time.Now().Add(720 * time.Hour).Unix())), // 30 days
			BusinessID:       strconv.Itoa(int(user.AccountID)),
		}
		err = paymentAccount.CreatePaymentAccount(db.Payment)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	businessID, _ := strconv.Atoi(paymentAccount.BusinessID)
	payment := models.Payment{PaymentID: paymentAccount.PaymentID}
	_, err = payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		payment = models.Payment{
			PaymentID:    utility.RandomString(10),
			TotalAmount:  amountPaid,
			EscrowCharge: 0,
			IsPaid:       false,
			AccountID:    int64(businessID),
			BusinessID:   int64(businessID),
			Currency:     currency,
		}

		err = payment.CreatePayment(db.Payment)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	payment.PaymentMadeAt = paidOn
	payment.PaidBy = paymentSourceInformation
	payment.UpdateAllFields(db.Payment)

	if payment.IsPaid {
		return http.StatusOK, nil
	}

	verified, _, err := monnify.VerifyTrans(generatedReference, amountPaid)
	if err != nil {
		extReq.Logger.Error("monnify webhhook log error", "error verifying transaction", err.Error())
		return http.StatusInternalServerError, fmt.Errorf("error verifying transaction")
	}

	if verified {
		if payment.TransactionID != "" {
			transaction, tErr := ListTransactionsByID(extReq, payment.TransactionID)
			escrowChargeBearerParty := transaction.Parties["escrow_charge_bearer"]
			if transaction.EscrowWallet == "yes" {
				payment.WalletFunded = "ESCROW_" + currency
			} else {
				payment.WalletFunded = currency
			}
			payment.UpdateAllFields(db.Payment)
			sendTransactionConfirmed(extReq, db, &payment, paymentReference, amountPaid)

			err = SlackNotify(extReq, paymentChannelD, `
			Bank Transfer Payment | WEB HOOK MONNIFY
			Environment: `+config.GetConfig().App.Name+`
			Transaction ID: `+transaction.TransactionID+`
			Payment ID: `+payment.PaymentID+`
			Account Number: `+paymentAccount.AccountNumber+`
			Account Name: `+paymentAccount.AccountName+`
			Bank: `+paymentAccount.BankName+`
			Amount: `+fmt.Sprintf("%v %v", currency, amountPaid)+`
			Escrow Charge: `+fmt.Sprintf("%v %v", currency, payment.EscrowCharge)+`
			Escrow Charge Bearer: `+fmt.Sprintf("%v", escrowChargeBearerParty.AccountID)+`
			Status: SUCCESSFUL
			`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}
			businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, int(transaction.BusinessID))
			if businessProfileData.Webhook_uri != "" {
				InitWebhook(extReq, db, businessProfileData.Webhook_uri, "bank-transfer.success", map[string]interface{}{
					"payment_id": payment.PaymentID,
					"amount":     amountPaid,
					"status":     "success",
				}, businessProfileData.AccountID)
			}
			if tErr != nil {
				return http.StatusBadRequest, fmt.Errorf("transaction with ID %v not found", payment.TransactionID)
			}
			fundAccount(extReq, db, amountPaid, currency, generatedReference, customerEmail, paymentReference, &payment, transaction)
		} else if paymentAccount.TransactionID != "" {
			transaction, err := ListTransactionsByID(extReq, paymentAccount.TransactionID)
			if err != nil {
				return http.StatusBadRequest, fmt.Errorf("transaction with ID %v not found", paymentAccount.TransactionID)
			}
			fundAccount(extReq, db, amountPaid, currency, generatedReference, customerEmail, paymentReference, &payment, transaction)
		} else {
			fundAccount(extReq, db, amountPaid, currency, generatedReference, customerEmail, paymentReference, &payment, external_models.TransactionByID{})
		}

	} else {
		if payment.TransactionID != "" {
			transaction, _ := ListTransactionsByID(extReq, payment.TransactionID)
			if transaction.EscrowWallet == "yes" {
				payment.WalletFunded = "ESCROW_" + currency
			} else {
				payment.WalletFunded = currency
			}
			payment.UpdateAllFields(db.Payment)

			businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, int(transaction.BusinessID))
			if businessProfileData.Webhook_uri != "" {
				InitWebhook(extReq, db, businessProfileData.Webhook_uri, "bank-transfer.failed", map[string]interface{}{
					"payment_id": payment.PaymentID,
					"amount":     amountPaid,
					"status":     "failed",
				}, businessProfileData.AccountID)
			}
		}

	}
	return http.StatusOK, nil
}

func MonnifyDisbursementCallbackService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.MonnifyWebhookRequest, requestBody []byte) (int, error) {
	var (
		secret               = config.GetConfig().Monnify.MonnifySecret
		data                 models.MonnifyWebhookRequestEventData
		monnifySignature     = utility.GetHeader(c, "monnify-signature")
		generatedReference   string
		amountPaid           float64
		disbursementStatus   string
		disbursementChannelD = config.GetConfig().Slack.DisbursementChannelID
	)

	hash := utility.Sha512Hmac(secret, requestBody)
	if hash != monnifySignature {
		extReq.Logger.Error("monnify callback log error", "Web Hook Denied, Hash Mismatch", hash, monnifySignature, requestBody)
		return http.StatusUnauthorized, fmt.Errorf("web Hook Denied, Hash Mismatch")
	}

	extReq.Logger.Info("monnify callback log info", string(requestBody))
	webhookLog := models.WebhookLog{
		Log:      string(requestBody),
		Provider: "monnify",
	}
	err := webhookLog.CreateWebhookLog(db.Payment)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if req.EventType == "" {
		extReq.Logger.Error("monnify callback log error", "no event type specified")
		return http.StatusBadRequest, fmt.Errorf("no event type specified")
	}

	if req.EventType != "FAILED_DISBURSEMENT" && req.EventType != "SUCCESSFUL_DISBURSEMENT" {
		extReq.Logger.Error("monnify callback log error", "not a disbursement event type: ", req.EventType)
		return http.StatusBadRequest, fmt.Errorf("not a disbursement event type")
	}

	if req.EventData != nil {
		data = *req.EventData
	} else {
		extReq.Logger.Error("monnify callback log error", "event data not found")
		return http.StatusBadRequest, fmt.Errorf("event data not found")
	}

	if data.Reference != nil {
		generatedReference = *data.Reference
	}

	if data.Amount != nil {
		amountPaid = *data.Amount
	}

	if data.Status != nil {
		disbursementStatus = *data.Status
	}

	err = SlackNotify(extReq, disbursementChannelD, `
					CallBack Received Monnify
					Environment: `+config.GetConfig().App.Name+`
					Event: `+req.EventType+`
					Reference: `+fmt.Sprintf("generated reference: %v", generatedReference)+`
					Status: SUCCESSFUL
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	if strings.EqualFold(disbursementStatus, "FAILED") && strings.EqualFold(req.EventType, "FAILED_DISBURSEMENT") {
		disbursement := models.Disbursement{Reference: generatedReference}
		code, err := disbursement.GetDisbursementByReference(db.Payment)
		if err != nil {
			extReq.Logger.Error("monnify callback log error", fmt.Sprintf("disbursement with reference %v, not found. Error: %v", generatedReference, err.Error()))
			return code, err
		}

		_, err = CreditWallet(extReq, db, amountPaid, disbursement.Currency, disbursement.BusinessID, true, "no", "")
		if err != nil {
			extReq.Logger.Error("monnify callback log error", fmt.Sprintf("crediting wallet for business id: %v, amount: %v %v, Error: %v", disbursement.BusinessID, disbursement.Currency, amountPaid, err.Error()))
		}
	}

	return http.StatusOK, nil
}

func fundAccount(extReq request.ExternalRequest, db postgresql.Databases, amountPaid float64, currency, generatedReference, customerEmail, paymentReference string, payment *models.Payment, transaction external_models.TransactionByID) error {
	var (
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)

	paymentAccount := models.PaymentAccount{PaymentAccountID: generatedReference}
	_, err := paymentAccount.GetLatestPaymentAccountByPaymentAccountID(db.Payment)
	if err != nil {
		extReq.Logger.Error("monnify webhhook log error", err.Error())
		return err
	}
	finalAmount := amountPaid

	user, _ := GetUserWithEmail(extReq, customerEmail)
	var fundingCharge float64 = 0
	if amountPaid < 500000 {
		fundingCharge = 500
	}

	if amountPaid > 500000 && amountPaid < 1000000 {
		fundingCharge = 1000
	}
	if amountPaid > 1000000 {
		fundingCharge = 2000
	}

	if transaction.TransactionID == "" {
		finalAmount = amountPaid - fundingCharge
	} else {
		fundingCharge = transaction.EscrowCharge
		finalAmount = amountPaid - transaction.EscrowCharge
	}

	if paymentReference != paymentAccount.PaymentReference {
		pendingTransferFunding := models.PendingTransferFunding{Reference: generatedReference}
		pendingTransferFunding.GetPendingTransferFundingByReference(db.Payment)
		pendingTransferFunding.Delete(db.Payment)

		paymentAccount.PaymentReference = paymentReference
		paymentAccount.UpdateAllFields(db.Payment)

		_, err = CreditWallet(extReq, db, finalAmount, currency, int(user.AccountID), false, thisOrThatStr(transaction.EscrowWallet, "no"), transaction.TransactionID)
		if err != nil {
			return err
		}
		payment.EscrowCharge = fundingCharge
		payment.IsPaid = true
		payment.BusinessID = int64(user.BusinessId)
		err = payment.UpdateAllFields(db.Payment)
		if err != nil {
			return err
		}

		err = SlackNotify(extReq, paymentChannelD, `
			[Web Hook MONNIFY] Wallet Funding For Customer #`+strconv.Itoa(int(user.AccountID))+`
			Environment: `+config.GetConfig().App.Name+`
			Account ID: `+strconv.Itoa(int(user.AccountID))+`
			Amount: `+fmt.Sprintf("%v %v", currency, amountPaid)+`
			Settled Amount: `+fmt.Sprintf("%v %v", currency, finalAmount)+`
			Fee: `+fmt.Sprintf("%v %v", currency, fundingCharge)+`
			Status: PAID
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}

		businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, int(user.AccountID))
		if businessProfileData.Webhook_uri != "" {
			InitWebhook(extReq, db, businessProfileData.Webhook_uri, "bank-transfer.success", map[string]interface{}{
				"reference": generatedReference,
				"amount":    amountPaid,
				"status":    "success",
			}, businessProfileData.AccountID)
		}
		paymentAccountByte, _ := json.Marshal(paymentAccount)
		extReq.Logger.Info("Monnify Bank Transfer Confirmed", "data:", string(paymentAccountByte))

	} else {
		businessProfileData, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, int(user.AccountID))
		if businessProfileData.Webhook_uri != "" {
			InitWebhook(extReq, db, businessProfileData.Webhook_uri, "bank-transfer.failed", map[string]interface{}{
				"reference": generatedReference,
				"amount":    amountPaid,
				"status":    "failed",
			}, businessProfileData.AccountID)
		}
		paymentAccountByte, _ := json.Marshal(paymentAccount)
		extReq.Logger.Error("Monnify Bank Transfer Not Confirmed", "data:", string(paymentAccountByte))

	}
	return nil
}
