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
		return msg, http.StatusBadRequest, fmt.Errorf("initial amount must be set if rate exist")
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

		convertedAmount := (rate.Amount / req.InitialAmount) * amount
		recipientAmount = convertedAmount
	} else {
		recipientAmount = amount
	}

	senderWallet, debitStatus, err := DebitWallet(extReq, db, amount, req.SenderCurrency, req.SenderAccountID, escrowWallet, req.TransactionID)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	} else if !debitStatus {
		return msg, http.StatusInternalServerError, fmt.Errorf("debit failed for sender")
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

	receiverWallet, err := CreditWallet(extReq, db, recipientAmount, recipientCurrency, req.RecipientAccountID, req.Refund, escrowWallet, req.TransactionID)
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
		data              models.ManualDebitResponse
		currency          = strings.ToUpper(req.Currency)
		walletCurrency    = currency
		businessName      = strconv.Itoa(req.AccountID)
		bankCode          = ""
		bankName          = ""
		bankAccountNumber = ""
		bankAccountName   = ""
		gateway           = thisOrThatStr(req.Gateway, "monnify")
		beneficiaryName   = ""
		email             = ""
		rave              = Rave{ExtReq: extReq}
		monnify           = Monnify{ExtReq: extReq}
	)

	if !HasBvn(extReq, uint(req.AccountID)) {
		return "", data, http.StatusBadRequest, fmt.Errorf("this account does not have bvn associated with it.")
	}

	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return "", data, http.StatusInternalServerError, fmt.Errorf("could not retrieve user info")
	}

	if !user.CanMakeWithdrawal {
		return "", data, http.StatusBadRequest, fmt.Errorf("user withdrawals not enabled, please contact customer care")
	}

	businessProfile, err := GetBusinessProfileByAccountID(extReq, extReq.Logger, req.AccountID)
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
			return "", data, http.StatusBadRequest, fmt.Errorf("Recipient does not have a first & last name")
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
		return "", data, http.StatusBadRequest, fmt.Errorf("Requested amount is greater than wallet balance")
	}

	walletBalance, debitStatus, err := DebitWallet(extReq, db, amount, currency, req.AccountID, req.EscrowWallet, "")
	if err != nil {
		return "", data, http.StatusInternalServerError, err
	}

	if !debitStatus {
		return "", data, http.StatusInternalServerError, fmt.Errorf("debit failed")
	}

	callback := utility.GenerateGroupByURL(c, "disbursement/callback", map[string]string{})
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
	log := ""

	switch strings.ToLower(gateway) {
	case "rave":
		resData, err := rave.InitTransfer(bankCode, bankAccountNumber, finalAmount, narration, currency, reference, callback)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
		data.Msg = resData.Message
		jsonByte, _ := json.Marshal(resData)
		log, gatewayData = string(jsonByte), resData
	case "monnify":
		resData, err := monnify.InitTransfer(finalAmount, reference, narration, bankCode, bankAccountNumber, currency, bankAccountName)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
		data.Msg = resData.ResponseMessage
		jsonByte, _ := json.Marshal(resData)
		log, gatewayData = string(jsonByte), resData
	default:
		resData, err := rave.InitTransfer(bankCode, bankAccountNumber, finalAmount, narration, currency, reference, callback)
		if err != nil {
			return "", data, http.StatusInternalServerError, err
		}
		data.Msg = resData.Message
		jsonByte, _ := json.Marshal(resData)
		log, gatewayData = string(jsonByte), resData
	}
	fmt.Println(beneficiaryName, email)

	LogDisbursement(db, disbursementID, log)
	data.Response = gatewayData

	return "Wallet Disbursement Queued", data, http.StatusOK, nil
}

func LogDisbursement(db postgresql.Databases, disbursementID int, log string) {
	disbursementLog := models.DisbursementLog{
		DisbursementID: disbursementID,
		Log:            log,
	}
	disbursementLog.CreateDisbursementLog(db.Payment)
}
