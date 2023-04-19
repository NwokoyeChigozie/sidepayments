package cronjobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/services/payment"
	"github.com/vesicash/payment-ms/utility"
)

func Disbursement(extReq request.ExternalRequest, db postgresql.Databases) {

	transactions, err := payment.ListTransactionsByStatusCode(extReq, "cdp", 1, 20)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting transactions, err: %v", err.Error()))
		return
	}

	for _, transaction := range transactions {
		err := beginDisbursement(extReq, db, transaction)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error running disbursement for transaction %v; err: %v", transaction.TransactionID, err.Error()))
		} else {
			extReq.Logger.Info(fmt.Sprintf("disbursement for transaction %v complete", transaction.TransactionID))
		}
	}
}

func beginDisbursement(extReq request.ExternalRequest, db postgresql.Databases, transaction external_models.TransactionByID) error {
	var (
		baseTime             time.Time
		disbursementChannelD = config.GetConfig().Slack.DisbursementChannelID
	)
	extReq.Logger.Info(fmt.Sprintf("%v : Disbursing...", transaction.TransactionID))
	paymnt := models.Payment{TransactionID: transaction.TransactionID}

	code, err := paymnt.GetPaymentByTransactionID(db.Payment)

	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("err getting payment for transaction %v, error: %v", transaction.TransactionID, err))
		if code == http.StatusInternalServerError {
			return err
		}
		return fmt.Errorf("payment record not found")
	}

	var (
		restricted = []string{}
		issues     = []string{}
	)
	sellerParty := transaction.Parties["seller"]
	user, _ := payment.GetUserWithAccountID(extReq, sellerParty.AccountID)

	bankDetails, err := payment.GetBankDetail(extReq, 0, int(user.AccountID), "", paymnt.Currency)
	if err != nil {
		restricted = append(restricted, "Bank Account")
	}

	if !paymnt.IsPaid && paymnt.PaymentMadeAt == baseTime {
		issues = append(issues, "Transaction has not been paid for")
	}

	if user.AccountType == "" {
		restricted = append(restricted, "Account Type")
	}

	if strings.EqualFold(user.AccountType, "individual") {
		profile, _ := payment.GetUserProfileByAccountID(extReq, extReq.Logger, int(user.AccountID))
		if strings.EqualFold(profile.Country, "Nigeria") || strings.EqualFold(profile.Country, "NG") {
			if !payment.HasVerification(extReq, user.AccountID, "bvn") {
				restricted = append(restricted, "Identity Card")
			}
		}
	}

	if strings.EqualFold(user.AccountType, "business") {
		businessprofile, _ := payment.GetBusinessProfileByAccountID(extReq, extReq.Logger, int(user.AccountID))
		if strings.EqualFold(businessprofile.Country, "Nigeria") || strings.EqualFold(businessprofile.Country, "NG") {
			if strings.EqualFold(businessprofile.BusinessType, "social_commerce") {
				if !payment.HasVerification(extReq, user.AccountID, "id") {
					restricted = append(restricted, "Identity Card")
				}
			} else {
				if !payment.HasVerification(extReq, user.AccountID, "cac") {
					restricted = append(restricted, "CAC")
				}
				if !payment.HasVerification(extReq, user.AccountID, "utilitybill") {
					restricted = append(restricted, "Utility Bill")
				}
			}

		}
	}

	businessId := transaction.BusinessID
	businessProfile, _ := payment.GetBusinessProfileByAccountID(extReq, extReq.Logger, businessId)
	businessCharge, _ := payment.GetBusinessChargeWithBusinessIDAndCountry(extReq, transaction.BusinessID, transaction.Country.CountryCode)

	if businessCharge.DisbursementGateway == "rave_momo" {
		momoBankDetails, err := payment.GetBankDetail(extReq, 0, int(user.AccountID), "", "", true)
		if err != nil {
			restricted = append(restricted, "Mobile Money Operator")
		}

		if !utility.InStringSlice(momoBankDetails.MobileMoneyOperator, []string{"mps", "mtn", "tigo", "vodafone", "airtel"}) {
			restricted = append(restricted, "Mobile Money Operator Not Supported")
		}
	}

	businesscheck, err := payment.GetBusinessProfileByAccountID(extReq, extReq.Logger, user.BusinessId)
	if businesscheck.IsVerificationWaved {
		restricted = []string{}
	}

	if len(restricted) > 0 {
		extReq.Logger.Info("transaction %v: Disbursement skipped: user has certain restrictions. %v", transaction.TransactionID, strings.Join(restricted, ", "))
		failedDisbursements := models.FailedDisbursement{PaymentID: paymnt.PaymentID}
		_, err := failedDisbursements.GetFailedDisbursementByPaymentID(db.Payment)
		if err == nil {
			return fmt.Errorf("transaction %v: Disbursement skipped: user has certain restrictions. %v", transaction.TransactionID, strings.Join(restricted, ", "))
		}
		failedDisbursements = models.FailedDisbursement{
			PaymentID: paymnt.PaymentID,
			Reasons:   strings.Join(restricted, ", "),
		}

		err = failedDisbursements.CreateFailedDisbursement(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("err creating failed disbursement %v", err.Error()))
			return err
		}

		err = payment.SlackNotify(extReq, disbursementChannelD, `
			Disbursement For Transaction #`+paymnt.TransactionID+`.
			Environment: `+config.GetConfig().App.Name+`
			Account ID: `+strconv.Itoa(int(user.AccountID))+`
			Beneficiary Name: `+fmt.Sprintf("%v %v", user.Firstname, user.Lastname)+`
			Amount: `+fmt.Sprintf("%v %v", paymnt.Currency, paymnt.TotalAmount)+`
			Status: NOT PROCESSED
			Reason: `+strings.Join(restricted, ", ")+`
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}

		return fmt.Errorf("transaction %v: Disbursement skipped: user has certain restrictions. %v", transaction.TransactionID, strings.Join(restricted, ", "))
	} else if len(issues) > 0 {
		extReq.Logger.Info("transaction %v: Disbursement skipped: user has certain issues. %v", transaction.TransactionID, strings.Join(issues, ", "))
		return fmt.Errorf("transaction %v: Disbursement skipped: user has certain issues. %v", transaction.TransactionID, strings.Join(issues, ", "))
	}

	disbursement := models.Disbursement{PaymentID: paymnt.PaymentID}
	code, err = disbursement.GetDisbursementByPaymentID(db.Payment)
	if err == nil {
		return nil
	}

	if code == http.StatusInternalServerError {
		extReq.Logger.Error(fmt.Sprintf("error getting disbursement with payment_id: %v; error: %v", paymnt.PaymentID, err.Error()))
		return fmt.Errorf("error getting disbursement with payment_id: %v; error: %v", paymnt.PaymentID, err.Error())
	}

	disbursementID := utility.GetRandomNumbersInRange(1000000000, 9999999999)
	reference := utility.GetRandomNumbersInRange(1000000000, 9999999999)
	if bankDetails.ID == 0 {
		extReq.Logger.Error(fmt.Sprintf("User %v does not have bank details for currency %v", user.AccountID, paymnt.Currency))
		return fmt.Errorf("User %v does not have bank details on file for currency %v", user.AccountID, paymnt.Currency)
	}

	bank, err := payment.GetBank(extReq, bankDetails.BankID, "", "", "")
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("bank with id %v not found, error: %v", bankDetails.BankID, err.Error()))
		return fmt.Errorf("bank with id %v not found, error: %v", bankDetails.BankID, err.Error())
	}

	var (
		bankCode = bank.Code
		amount   float64
		currency string
		rave     = payment.Rave{ExtReq: extReq}
		monnify  = payment.Monnify{ExtReq: extReq}
	)

	if paymnt.DisburseCurrency != "" {
		if strings.EqualFold(paymnt.Currency, paymnt.DisburseCurrency) {
			converted, err := rave.ConvertCurrency(paymnt.TotalAmount, paymnt.Currency, paymnt.DisburseCurrency)
			if err != nil {
				extReq.Logger.Error(fmt.Sprintf("error converting currency with rave, from : %v, to: %v, amount: %v", paymnt.Currency, paymnt.DisburseCurrency, paymnt.TotalAmount))
				return fmt.Errorf("error converting currency with rave, from : %v, to: %v, amount: %v", paymnt.Currency, paymnt.DisburseCurrency, paymnt.TotalAmount)
			}
			amount = converted.Converted
			currency = paymnt.DisburseCurrency
		}
	} else {
		amount = paymnt.TotalAmount
		currency = paymnt.Currency
	}

	callback := utility.GenerateGroupByURL(config.GetConfig().App.Url, "/disbursement/callback", map[string]string{})
	disbursement = models.Disbursement{
		RecipientID:           int(user.AccountID),
		PaymentID:             paymnt.PaymentID,
		DisbursementID:        disbursementID,
		Reference:             strconv.Itoa(reference),
		Currency:              currency,
		BusinessID:            businessId,
		Amount:                fmt.Sprintf("%v", amount),
		CallbackUrl:           callback,
		BeneficiaryName:       bankDetails.AccountName,
		BankAccountNumber:     bankDetails.AccountNo,
		BankName:              bankDetails.BankName,
		DestinationBranchCode: "",
		DebitCurrency:         currency,
		Status:                "pending",
		Type:                  "credit",
	}

	if strings.EqualFold(businessProfile.DisbursementSettings, "wallet") {
		_, err := payment.CreditWallet(extReq, db, paymnt.TotalAmount, transaction.Currency, int(user.AccountID), false, "no", transaction.TransactionID)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error crediting wallet for user %v, amount:%v, currency:%v, %v", user.AccountID, paymnt.TotalAmount, transaction.Currency, err.Error()))
			return fmt.Errorf("error crediting wallet for user %v, amount:%v, currency:%v, %v", user.AccountID, paymnt.TotalAmount, transaction.Currency, err.Error())
		}

		disbursement.Gateway = "wallet"
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "completed"
		err = disbursement.CreateDisbursement(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error creating disbursement %v", err.Error()))
			return fmt.Errorf("error creating disbursement %v", err.Error())
		}

		extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
			AccountID:     int(paymnt.AccountID),
			TransactionID: transaction.TransactionID,
			MilestoneID:   transaction.MilestoneID,
			Status:        "closed",
		})

		extReq.SendExternalRequest(request.EscrowDisbursedSellerNotification, external_models.OnlyTransactionIDRequiredRequest{
			TransactionID: paymnt.TransactionID,
		})
		extReq.SendExternalRequest(request.EscrowDisbursedBuyerNotification, external_models.OnlyTransactionIDRequiredRequest{
			TransactionID: paymnt.TransactionID,
		})
		extReq.SendExternalRequest(request.TransactionClosedBuyerNotification, external_models.OnlyTransactionIDRequiredRequest{
			TransactionID: paymnt.TransactionID,
		})
		extReq.SendExternalRequest(request.TransactionClosedSellerNotification, external_models.OnlyTransactionIDRequiredRequest{
			TransactionID: paymnt.TransactionID,
		})

		err = payment.SlackNotify(extReq, disbursementChannelD, `
		 	Payment Disbursement For Transaction #`+paymnt.TransactionID+` has been completed successfully.
            Environment: `+config.GetConfig().App.Name+`
            Account ID: `+fmt.Sprintf("%v", disbursement.RecipientID)+`
            Beneficiary Name: `+disbursement.BeneficiaryName+`
            Amount: `+fmt.Sprintf("%v %v", disbursement.DebitCurrency, disbursement.Amount)+`
            Status: COMPLETED
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}

		return nil
	}

	err = disbursement.CreateDisbursement(db.Payment)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error creating disbursement %v", err.Error()))
		return fmt.Errorf("error creating disbursement %v", err.Error())
	}

	gateway := "rave"
	paymentInfo := models.PaymentInfo{PaymentID: paymnt.PaymentID, Status: "paid"}
	paymentInfo.GetPaymentInfoByPaymentIDAndStatus(db.Payment)

	if paymentInfo.Gateway != "" {
		gateway = paymentInfo.Gateway
	}

	disbursement.Gateway = gateway
	disbursement.UpdateAllFields(db.Payment)
	if err != nil {
		return err
	}

	if !strings.EqualFold(transaction.Currency, "NGN") {
		disbursement.Status = "manual"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			return err
		}
		extReq.Logger.Info(fmt.Sprintf("payment disbursement requires manual check for transaction %v", transaction.TransactionID))
		err = payment.SlackNotify(extReq, disbursementChannelD, `
		 	Disbursement For Transaction #`+paymnt.TransactionID+` requires manual check.
            Environment: `+config.GetConfig().App.Name+`
            Account ID: `+fmt.Sprintf("%v", disbursement.RecipientID)+`
            Beneficiary Name: `+disbursement.BeneficiaryName+`
            Amount: `+fmt.Sprintf("%v %v", disbursement.DebitCurrency, disbursement.Amount)+`
            Status: REVIEW
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
		return nil
	}

	var (
		narration  = "Vesicash"
		requestLog interface{}
	)

	switch strings.ToLower(gateway) {
	case "rave":
		resData, err := rave.InitTransfer(bankCode, bankDetails.AccountNo, amount, narration, currency, strconv.Itoa(reference), "")
		if err != nil {
			return err
		}
		requestLog = resData
	case "monnify":
		resData, err := monnify.InitTransfer(amount, strconv.Itoa(reference), narration, bankCode, bankDetails.AccountNo, currency, fmt.Sprintf("%v %v", user.Firstname, user.Lastname))
		if err != nil {
			return err
		}
		requestLog = resData
	default:
		resData, err := rave.InitTransfer(bankCode, bankDetails.AccountNo, amount, narration, currency, strconv.Itoa(reference), "")
		if err != nil {
			return err
		}
		requestLog = resData
	}

	payment.LogDisbursement(db, disbursementID, requestLog)

	disbursement.Status = "pending"
	err = disbursement.UpdateAllFields(db.Payment)
	if err != nil {
		return err
	}

	requestlogbyte, _ := json.Marshal(requestLog)

	disbursementRequestLog := models.DisbursementRequestLog{
		DisbursementID: strconv.Itoa(disbursementID),
		Log:            string(requestlogbyte),
	}
	err = disbursementRequestLog.CreateDisbursementRequestLog(db.Payment)
	if err != nil {
		return err
	}

	extReq.Logger.Info(fmt.Sprintf("Payment disbursement request sent for transaction %v", transaction.TransactionID))
	return nil
}
