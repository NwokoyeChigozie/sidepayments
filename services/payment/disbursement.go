package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func WalletTransferService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.WalletTransferRequest) (string, int, error) {
	var (
		msg string
	)
	if req.RateID != 0 && req.InitialAmount == 0 {
		return msg, http.StatusBadRequest, fmt.Errorf("initial amount must be set or greater than 0 if rate exist")
	}
	if req.RateID == 0 && req.InitialAmount != 0 {
		return msg, http.StatusBadRequest, fmt.Errorf("rate must be set if initial amount exist")
	}

	sender, err := GetUserWithAccountID(extReq, req.SenderAccountID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	_, err = GetUserWithAccountID(extReq, req.RecipientAccountID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	if !sender.CanMakeWithdrawal {
		return msg, http.StatusForbidden, fmt.Errorf("user withdrawal not enabled, please contact customer care")
	}

	if req.RateID != 0 && !sender.CanExchange {
		return msg, http.StatusForbidden, fmt.Errorf("user currency exchange not enabled, please contact customer care")
	}

	var (
		recipientCurrency      = strings.ToUpper(req.RecipientCurrency)
		senderCurrency         = strings.ToUpper(req.SenderCurrency)
		senderAvailableBalance float64
		amount                 float64
		initialCurrency        = strings.ReplaceAll(senderCurrency, "ESCROW_", "")
		finalCurrency          = strings.ReplaceAll(recipientCurrency, "ESCROW_", "")
		recipientAmount        float64
		escrowWallet           = "no"
	)

	senderWallet, err := GetWalletBalanceByAccountIdAndCurrency(extReq, req.SenderAccountID, senderCurrency)
	if err != nil {
		return msg, http.StatusBadRequest, fmt.Errorf("sender wallet does not exist")
	}
	senderAvailableBalance = senderWallet.Available

	if req.RateID != 0 && req.InitialAmount > 0 {
		amount = req.InitialAmount
	} else {
		amount = req.FinalAmount
	}

	if senderAvailableBalance < amount {
		return msg, http.StatusBadRequest, fmt.Errorf("insufficient balance")
	}

	if req.InitialAmount > 0 && req.RateID != 0 {
		if initialCurrency != finalCurrency {
			err := CreateExchangeTransaction(extReq, req.RecipientAccountID, req.RateID, req.InitialAmount, req.FinalAmount, ExchangeTransactionCompleted)
			if err != nil {
				return msg, http.StatusInternalServerError, err
			}
		}
	}

	if req.InitialAmount > 0 && req.RateID != 0 {
		rate, err := GetRateByID(extReq, req.RateID)
		if err != nil {
			return msg, http.StatusBadRequest, err
		}

		var multiplier float64 = 0
		if rate.InitialAmount > 0 {
			multiplier = rate.Amount / rate.InitialAmount
		}

		convertedAmount := multiplier * amount
		recipientAmount = convertedAmount
	} else {
		recipientAmount = amount
	}

	senderWallet, err = DebitWallet(extReq, db, amount, req.SenderCurrency, req.SenderAccountID, GetWalletType(escrowWallet, ""), req.TransactionID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	if utility.InStringSlice(strings.ToUpper(req.SenderCurrency), []string{"NGN", "ESCROW_NGN", "USD", "GBP", "ESCROW_USD", "ESCROW_GBP"}) {
		if utility.InStringSlice(strings.ToUpper(req.SenderCurrency), []string{"NGN", "ESCROW_NGN"}) && amount >= config.GetConfig().ONLINE_PAYMENT.NairaThreshold {
			_, err := CreateWalletTransaction(extReq, req.SenderAccountID, req.RecipientAccountID, amount, recipientAmount, senderCurrency, recipientCurrency, WalletTransactionPending, false)
			if err != nil {
				return msg, http.StatusInternalServerError, err
			}
			return "Wallet transfer queued for approval by admin", http.StatusOK, nil
		}
	}

	receiverWallet, err := CreditWallet(extReq, db, recipientAmount, recipientCurrency, req.RecipientAccountID, req.Refund, GetWalletType(escrowWallet, ""), req.TransactionID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	err = SaveWalletHistory(extReq, req.SenderAccountID, req.RecipientAccountID, amount, recipientAmount, senderCurrency, recipientCurrency, senderWallet.Available, receiverWallet.Available)
	if err != nil {
		extReq.Logger.Error("error saving wallet history: ", err.Error())
	}

	if req.TransactionID != "" {
		transaction, _ := ListTransactionsByID(extReq, req.TransactionID)
		transactionTitle := strings.Split(transaction.Title, ";")[0]
		action := "+"
		description := fmt.Sprintf("A sum of %v %v has been paid for this transaction", recipientCurrency, recipientAmount)
		if req.Refund {
			action = "-"
			description = fmt.Sprintf("A sum of %v %v has been deducted based on an excess payment being made on transaction %v", recipientCurrency, recipientAmount, transactionTitle)
		}

		extReq.SendExternalRequest(request.CreateActivityLog, external_models.CreateActivityLogRequest{
			TransactionID: req.TransactionID,
			Description:   description,
		})

		UpdateTransactionAmountPaid(extReq, req.TransactionID, amount, action)
	}

	return "Wallet transfer successful", http.StatusOK, nil
}

func ManualDebitService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.ManualDebitRequest) (string, models.ManualDebitResponse, int, error) {
	var (
		data                 models.ManualDebitResponse
		currency             = strings.ToUpper(req.Currency)
		walletCurrency       = currency
		businessName         = strconv.Itoa(req.AccountID)
		bankCode             = ""
		bankName             = ""
		bankAccountNumber    = ""
		bankAccountName      = ""
		gateway              = thisOrThatStr(req.Gateway, "monnify")
		beneficiaryName      = ""
		email                = ""
		rave                 = Rave{ExtReq: extReq}
		monnify              = Monnify{ExtReq: extReq}
		disbursementChannelD = config.GetConfig().Slack.DisbursementChannelID
	)

	if !HasBvn(extReq, uint(req.AccountID)) {
		return "", data, http.StatusBadRequest, fmt.Errorf("this account does not have bvn associated with it")
	}

	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return "", data, http.StatusInternalServerError, fmt.Errorf("could not retrieve user info")
	}

	if !user.CanMakeWithdrawal {
		return "", data, http.StatusBadRequest, fmt.Errorf("user withdrawals not enabled, please contact customer care")
	}

	businessProfile, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, req.AccountID)
	if businessProfile.BusinessName != "" {
		businessName = businessProfile.BusinessName
	}

	disbursementID := utility.GetRandomNumbersInRange(1000000000, 9999999999)
	reference := fmt.Sprintf("vc%v", utility.GetRandomNumbersInRange(1000000000, 9999999999))

	if req.Firstname != "" && req.Lastname != "" {
		beneficiaryName = fmt.Sprintf("%v %v", req.Firstname, req.Lastname)
	} else if (req.Firstname == "" && req.Lastname == "") && req.BeneficiaryName != "" {
		beneficiaryName = req.BeneficiaryName
	} else {
		if user.Firstname == "" && user.Lastname == "" {
			return "", data, http.StatusBadRequest, fmt.Errorf("recipient does not have a first & last name")
		}
		beneficiaryName = fmt.Sprintf("%v %v", user.Firstname, user.Lastname)
	}

	if req.Email != "" {
		email = req.Email
	} else {
		email = user.EmailAddress
	}

	if req.BankAccountNumber != "" && req.BankCode != "" {
		bankCode, bankAccountNumber, bankAccountName = req.BankCode, req.BankAccountNumber, req.BankAccountName
	} else {
		bankDetails, err := GetBankDetail(extReq, 0, req.AccountID, "", currency)
		if err != nil {
			return "", data, http.StatusBadRequest, fmt.Errorf("user has no %v bank account details", currency)
		}

		bank, err := GetBank(extReq, bankDetails.BankID, "", "", "")
		if err != nil {
			return "", data, http.StatusBadRequest, fmt.Errorf("bank with id %v not found", bankDetails.BankID)
		}
		bankCode = bank.Code
		bankName = bank.Name
		bankAccountNumber = bankDetails.AccountNo
		bankAccountName = bankDetails.AccountName
	}

	if req.EscrowWallet == "yes" {
		walletCurrency = fmt.Sprintf("ESCROW_%v", currency)
	}

	walletBalance, err := GetWalletBalanceByAccountIdAndCurrency(extReq, req.AccountID, walletCurrency)
	if err != nil {
		walletBalance, err = CreateWalletBalance(extReq, req.AccountID, walletCurrency, 0)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
	}

	amount := req.Amount
	disbursementCharge := config.GetConfig().ONLINE_PAYMENT.DisbursementCharge
	finalAmount := amount - disbursementCharge

	if amount > walletBalance.Available {
		return "", data, http.StatusBadRequest, fmt.Errorf("requested amount is greater than wallet balance")
	}

	walletBalance, err = DebitWallet(extReq, db, amount, currency, req.AccountID, GetWalletType(req.EscrowWallet, ""), "")
	if err != nil {
		return "", data, http.StatusInternalServerError, err
	}

	callback := utility.GenerateGroupByURL(config.GetConfig().App.Url, "/disbursement/callback", map[string]string{})
	disbursement := models.Disbursement{Reference: reference, Status: "new"}
	code, err := disbursement.GetDisbursementByReferenceAndNotStatus(db.Payment)
	if err != nil {
		if code == http.StatusInternalServerError {
			return "", data, http.StatusInternalServerError, err
		}

		payment := models.Payment{
			PaymentID:    utility.RandomString(10),
			TotalAmount:  req.Amount,
			EscrowCharge: disbursementCharge,
			IsPaid:       false,
			AccountID:    int64(req.AccountID),
			BusinessID:   int64(businessProfile.AccountID),
			Currency:     strings.ToUpper(req.Currency),
		}

		err = payment.CreatePayment(db.Payment)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}

		disbursement = models.Disbursement{
			RecipientID:           req.AccountID,
			PaymentID:             payment.PaymentID,
			DisbursementID:        disbursementID,
			Reference:             reference,
			Currency:              currency,
			BusinessID:            user.BusinessId,
			Amount:                fmt.Sprintf("%v", finalAmount),
			CallbackUrl:           callback,
			BeneficiaryName:       bankAccountName,
			BankAccountNumber:     bankAccountNumber,
			BankName:              bankName,
			DestinationBranchCode: req.DestinationBranchCode,
			DebitCurrency:         strings.ToUpper(req.DebitCurrency),
			Gateway:               gateway,
			Fee:                   int(disbursementCharge),
			Status:                "new",
			Type:                  "wallet",
			Approved:              "pending",
		}
		err = disbursement.CreateDisbursement(db.Payment)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
	}

	if req.EscrowWallet != "yes" {
		if (currency == "NGN" && amount >= config.GetConfig().ONLINE_PAYMENT.NairaThreshold) || currency == "USD" || currency == "GBP" {
			err = SlackNotify(extReq, disbursementChannelD, `
				Wallet Debit To Bank Account #`+strconv.Itoa(req.AccountID)+`
                Environment: `+config.GetConfig().App.Name+`
                Disbursement ID: `+strconv.Itoa(disbursementID)+`
                User: `+fmt.Sprintf("%v %v", user.Firstname, user.Lastname)+`,
                Amount: `+fmt.Sprintf("%v %v", currency, finalAmount)+`
                Status: Pending Admin Approval
			`)
			if err != nil && !extReq.Test {
				extReq.Logger.Error("error sending notification to slack: ", err.Error())
			}
			return "Wallet Disbursement Queued for approval by admin", models.ManualDebitResponse{DisbursementID: disbursementID, Status: disbursement.Status}, http.StatusOK, nil
		}
	}

	data.DisbursementID = disbursementID
	disbursement.Approved = "yes"
	err = disbursement.UpdateAllFields(db.Payment)
	if err != nil {
		return "", data, http.StatusInternalServerError, err
	}

	data.Status = "new"
	narration := fmt.Sprintf("%v/VES", businessName)
	var gatewayData interface{}

	switch strings.ToLower(gateway) {
	case "rave":
		resData, err := rave.InitTransfer(bankCode, bankAccountNumber, finalAmount, narration, currency, reference, "")
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
		data.Msg = resData.Message
		gatewayData = resData
	case "monnify":
		resData, err := monnify.InitTransfer(finalAmount, reference, narration, bankCode, bankAccountNumber, currency, bankAccountName)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
		data.Msg = resData.ResponseMessage
		gatewayData = resData
	default:
		resData, err := rave.InitTransfer(bankCode, bankAccountNumber, finalAmount, narration, currency, reference, "")
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
		data.Msg = resData.Message
		gatewayData = resData
	}
	fmt.Println(beneficiaryName, email)

	LogDisbursement(db, disbursementID, gatewayData)
	data.Response = gatewayData
	err = SlackNotify(extReq, disbursementChannelD, `
				Wallet Debit To Bank Account #`+strconv.Itoa(req.AccountID)+`
                Environment: `+config.GetConfig().App.Name+`
                Disbursement ID: `+strconv.Itoa(disbursementID)+`
                User: `+fmt.Sprintf("%v %v", user.Firstname, user.Lastname)+`,
                Amount: `+fmt.Sprintf("%v %v", currency, finalAmount)+`
                Status: INITIATED
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}
	return "Wallet Disbursement Queued", data, http.StatusOK, nil
}

func ManualRefundService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.ManualRefundRequest) (string, int, error) {
	var (
		response             string
		rave                 = Rave{ExtReq: extReq}
		disbursementChannelD = config.GetConfig().Slack.DisbursementChannelID
	)

	transaction, err := ListTransactionsByID(extReq, req.TransactionID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
		AccountID:     transaction.BusinessID,
		TransactionID: transaction.TransactionID,
		MilestoneID:   transaction.MilestoneID,
		Status:        "fr",
	})

	payment, code, err := getPaymentByTransactionID(db, req.TransactionID)
	if err != nil {
		return response, code, err
	}

	buyerParty, ok := transaction.Parties["buyer"]
	if !ok {
		return response, http.StatusBadRequest, fmt.Errorf("transaction has no buyer")
	}
	sellerParty, ok := transaction.Parties["seller"]
	if !ok {
		return response, http.StatusBadRequest, fmt.Errorf("transaction has no seller")
	}
	businessID := transaction.BusinessID

	countryCode := transaction.Country.CountryCode
	currencyCode := transaction.Country.CurrencyCode

	if countryCode == "" {
		country, err := GetCountryByCurrency(extReq, extReq.Logger, payment.Currency)
		if err != nil {
			return response, http.StatusBadRequest, fmt.Errorf("error retreiving country: %v", err.Error())
		}
		countryCode = country.CountryCode
		currencyCode = country.CurrencyCode
	}

	businessCharges, err := GetBusinessChargeWithBusinessIDAndCountry(extReq, transaction.BusinessID, countryCode)
	if err != nil {
		businessCharges, err = InitBusinessCharge(extReq, transaction.BusinessID, currencyCode)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}
	}

	disbursement := models.Disbursement{PaymentID: payment.PaymentID}
	code, err = disbursement.GetDisbursementByPaymentID(db.Payment)
	if err != nil {
		if code == http.StatusInternalServerError {
			return response, code, err
		}
	} else {
		return response, http.StatusBadRequest, fmt.Errorf("disbursement already exists")
	}

	// disbursementGateway := businessCharges.DisbursementGateway
	cancellationFee, _ := strconv.ParseFloat(businessCharges.CancellationFee, 64)

	sellerInfo, err := GetUserWithAccountID(extReq, sellerParty.AccountID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}
	buyerInfo, err := GetUserWithAccountID(extReq, buyerParty.AccountID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	disbursementID := utility.GetRandomNumbersInRange(1000000000, 9999999999)
	reference := fmt.Sprintf("vc%v", utility.GetRandomNumbersInRange(1000000000, 9999999999))
	currency := strings.ToUpper(payment.Currency)

	bankDetails, err := GetBankDetail(extReq, 0, buyerParty.AccountID, "", currency)
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("user has no %v bank account details", currency)
	}

	bank, err := GetBank(extReq, bankDetails.BankID, "", "", "")
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("bank with id %v not found", bankDetails.BankID)
	}
	bankCode := bank.Code

	// if strings.EqualFold(disbursementGateway, "rave_momo") {
	// 	bankCode = bankDetails.MobileMoneyOperator
	// }

	businessProfile, err := GetBusinessProfileByAccountID(extReq, extReq.Logger, businessID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	currency = strings.ToUpper(transaction.Currency)

	var realAmount float64

	if cancellationFee != 0 {
		realAmount = payment.TotalAmount - cancellationFee
	} else {
		if payment.TotalAmount == transaction.TotalAmount {
			realAmount = payment.TotalAmount
		} else {
			realAmount = payment.TotalAmount - payment.EscrowCharge
		}
	}
	callback := utility.GenerateGroupByURL(config.GetConfig().App.Url, "/disbursement/callback", map[string]string{})
	disbursement = models.Disbursement{
		RecipientID:           buyerParty.AccountID,
		PaymentID:             payment.PaymentID,
		DisbursementID:        disbursementID,
		Reference:             reference,
		Currency:              currency,
		BusinessID:            businessID,
		Amount:                fmt.Sprintf("%v", realAmount),
		CallbackUrl:           callback,
		BeneficiaryName:       buyerInfo.Firstname,
		DestinationBranchCode: "0",
		DebitCurrency:         currency,
		Status:                "pending",
		Type:                  "refund",
	}

	if businessProfile.DisbursementSettings == "wallet" {
		_, err := CreditWallet(extReq, db, payment.TotalAmount, currency, businessID, true, "no", transaction.TransactionID)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}
		disbursement.Gateway = "wallet"
		disbursement.Status = "completed"
		err = disbursement.CreateDisbursement(db.Payment)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}

		err = SlackNotify(extReq, disbursementChannelD, `
				Wallet Disbursement Re-Fund For Buyer #`+strconv.Itoa(buyerParty.AccountID)+`
                Environment: `+config.GetConfig().App.Name+`
                Disbursement ID: `+strconv.Itoa(disbursementID)+`
                Seller: `+fmt.Sprintf("%v %v", sellerInfo.Firstname, sellerInfo.Lastname)+`,
                Buyer: `+fmt.Sprintf("%v %v", buyerInfo.Firstname, buyerInfo.Lastname)+`,
                Amount: `+fmt.Sprintf("%v %v", currency, payment.TotalAmount)+`
                Status: successful
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
		return fmt.Sprintf("%v : Payment Disbursed to Wallet", transaction.TransactionID), http.StatusOK, nil
	}

	disbursement.Gateway = "rave"
	err = disbursement.CreateDisbursement(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	resData, err := rave.InitTransfer(bankCode, bankDetails.AccountNo, realAmount, "Vesicash Refund", currency, reference, "")
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	disbursement.Status = strings.ToLower(resData.Status)
	disbursement.UpdateAllFields(db.Payment)
	LogDisbursement(db, disbursementID, resData)

	err = SlackNotify(extReq, disbursementChannelD, `
				Bank Account Disbursement Re-Fund For Buyer #`+strconv.Itoa(buyerParty.AccountID)+`
                Environment: `+config.GetConfig().App.Name+`
                Disbursement ID: `+strconv.Itoa(disbursementID)+`
                Seller: `+fmt.Sprintf("%v %v", sellerInfo.Firstname, sellerInfo.Lastname)+`,
                Buyer: `+fmt.Sprintf("%v %v", buyerInfo.Firstname, buyerInfo.Lastname)+`,
                Amount: `+fmt.Sprintf("%v %v", currency, realAmount)+`
                Status: INITIATED
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	return fmt.Sprintf("%v : Refund disbursement request sent", transaction.TransactionID), http.StatusOK, nil
}

func getPaymentByTransactionID(db postgresql.Databases, transactionID string) (models.ListPayment, int, error) {

	payment := models.Payment{TransactionID: transactionID}
	code, err := payment.GetPaymentByTransactionID(db.Payment)
	if err != nil {
		return models.ListPayment{}, code, err
	}

	escrowCharge := 0
	brokerCharge := 0
	shippingFee := 0

	return models.ListPayment{
		ID:               payment.ID,
		PaymentID:        payment.PaymentID,
		TransactionID:    payment.TransactionID,
		TotalAmount:      payment.TotalAmount,
		EscrowCharge:     payment.EscrowCharge,
		IsPaid:           payment.IsPaid,
		PaymentMadeAt:    payment.PaymentMadeAt,
		DeletedAt:        payment.DeletedAt,
		CreatedAt:        payment.CreatedAt,
		UpdatedAt:        payment.UpdatedAt,
		AccountID:        payment.AccountID,
		BusinessID:       payment.BusinessID,
		Currency:         payment.Currency,
		ShippingFee:      payment.ShippingFee,
		DisburseCurrency: payment.DisburseCurrency,
		PaymentType:      payment.PaymentType,
		BrokerCharge:     payment.BrokerCharge,
		SummedAmount:     payment.TotalAmount + float64(shippingFee) + float64(brokerCharge) + float64(escrowCharge),
	}, http.StatusOK, nil
}
func LogDisbursement(db postgresql.Databases, disbursementID int, logData interface{}) {
	jsonByte, _ := json.Marshal(logData)
	disbursementLog := models.DisbursementLog{
		DisbursementID: disbursementID,
		Log:            string(jsonByte),
	}
	disbursementLog.CreateDisbursementLog(db.Payment)
}
