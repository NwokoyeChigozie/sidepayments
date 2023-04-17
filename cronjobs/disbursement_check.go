package cronjobs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/services/payment"
)

var (
	maxTries = 3
)

func DisbursementCheck(extReq request.ExternalRequest, db postgresql.Databases) {
	var (
		rave         = payment.Rave{ExtReq: extReq}
		monnify      = payment.Monnify{ExtReq: extReq}
		disbursement = models.Disbursement{}
	)
	allPendingDisbursements, err := disbursement.GetAllForStatuses(db.Payment, []string{"new", "pending"})
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting new and pending disbursements %v", err.Error()))
		return
	}

	for _, item := range allPendingDisbursements {
		paymentGateway := item.Gateway
		if paymentGateway == "" {
			paymentGateway = "rave"
		}
		var (
			log          interface{}
			status       bool
			statusString string
			reference    = item.Reference
		)

		switch strings.ToLower(paymentGateway) {
		case "rave":
			log, status, statusString, _, err = rave.Status(reference)
			if err != nil {
				dataByte, _ := json.Marshal(log)
				extReq.Logger.Error(fmt.Sprintf("error getting transaction status rave,reference:%v data:%v, error: %v", reference, string(dataByte), err.Error()))
			}
		case "monnify":
			log, status, statusString, _, err = monnify.Status(reference)
			if err != nil {
				dataByte, _ := json.Marshal(log)
				extReq.Logger.Error(fmt.Sprintf("error getting transaction status monnify,reference:%v data:%v, error: %v", reference, string(dataByte), err.Error()))
			}
		default:
			log, status, statusString, _, err = rave.Status(reference)
			if err != nil {
				dataByte, _ := json.Marshal(log)
				extReq.Logger.Error(fmt.Sprintf("error getting transaction status rave,reference:%v data:%v, error: %v", reference, string(dataByte), err.Error()))
			}
		}

		payment.LogDisbursement(db, disbursement.DisbursementID, log)

		if item.PaymentID == "" || item.PaymentID == "0" {
			//TODO complete walletConfirm function
			walletConfirm(extReq, db, disbursement, status, statusString)
		} else {
			//TODO complete transConfirm function
			transConfirm(extReq, db, disbursement, status, statusString)
		}
	}

}

func walletConfirm(extReq request.ExternalRequest, db postgresql.Databases, disbursement models.Disbursement, status bool, statusString string) {
	var (
		disbursementChannelD = config.GetConfig().Slack.DisbursementChannelID
	)
	if strings.EqualFold(statusString, "completed") || strings.EqualFold(statusString, "done") {
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "completed"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}

		businessProfileData, _ := payment.GetBusinessProfileByAccountID(extReq, extReq.Logger, disbursement.BusinessID)
		if businessProfileData.Webhook_uri != "" {
			payment.InitWebhook(extReq, db, businessProfileData.Webhook_uri, "disbursement.success", map[string]interface{}{
				"disbursement_id": disbursement.DisbursementID,
				"reference":       disbursement.Reference,
				"status":          "success",
			}, businessProfileData.AccountID)
		}

		err = payment.SlackNotify(extReq, disbursementChannelD, `
		 	Wallet Disbursement been completed successfully.
            Environment: `+config.GetConfig().App.Name+`
            Account ID: `+fmt.Sprintf("%v", disbursement.RecipientID)+`
            Beneficiary Name: `+disbursement.BeneficiaryName+`
            Amount: `+fmt.Sprintf("%v %v", disbursement.DebitCurrency, disbursement.Amount)+`
            Status: COMPLETED
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement completed", disbursement.Reference))

	} else if strings.EqualFold(statusString, "failed") {
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "failed"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}

		amount, _ := strconv.ParseFloat(disbursement.Amount, 64)
		amount = amount + float64(disbursement.Fee)
		_, err = payment.CreditWallet(extReq, db, amount, disbursement.DebitCurrency, disbursement.RecipientID, true, "no", "")
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error crediting wallet for user %v, amount:%v; currency: %v; disbursementId:%v; error: %v", disbursement.RecipientID, amount, disbursement.DebitCurrency, disbursement.DisbursementID, err.Error()))
		}

		businessProfileData, _ := payment.GetBusinessProfileByAccountID(extReq, extReq.Logger, disbursement.BusinessID)
		if businessProfileData.Webhook_uri != "" {
			payment.InitWebhook(extReq, db, businessProfileData.Webhook_uri, "disbursement.failed", map[string]interface{}{
				"disbursement_id": disbursement.DisbursementID,
				"reference":       disbursement.Reference,
				"status":          "failed",
			}, businessProfileData.AccountID)
		}

		err = payment.SlackNotify(extReq, disbursementChannelD, `
		 	Wallet Disbursement Has Failed.
			Note: Requires manual disbursement
            Environment: `+config.GetConfig().App.Name+`
            Account ID: `+fmt.Sprintf("%v", disbursement.RecipientID)+`
            Beneficiary Name: `+disbursement.BeneficiaryName+`
            Amount: `+fmt.Sprintf("%v %v", disbursement.DebitCurrency, disbursement.Amount)+`
			Fee: `+fmt.Sprintf("%v %v", disbursement.DebitCurrency, disbursement.Fee)+`
            Status: FAILED
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement failed", disbursement.Reference))

	} else if strings.EqualFold(statusString, "ongoing") {
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement initiated/ongoing", disbursement.Reference))
	} else if strings.EqualFold(statusString, "cancelled") {
		disbursement.Status = "cancelled"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}

	} else if statusString == "" {
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "review"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}

		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement null or not found", disbursement.Reference))
	} else {
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement status not be determined", disbursement.Reference))
	}

}

func transConfirm(extReq request.ExternalRequest, db postgresql.Databases, disbursement models.Disbursement, status bool, statusString string) {
	var (
		paymnt               = models.Payment{PaymentID: disbursement.PaymentID}
		amount, _            = strconv.ParseFloat(disbursement.Amount, 64)
		disbursementChannelD = config.GetConfig().Slack.DisbursementChannelID
	)

	tries, isTrue := decideTries(extReq, db, &disbursement)
	if !isTrue {
		return
	}

	_, err := paymnt.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting payment for payment id: %v, error: %v", disbursement.PaymentID, err.Error()))
		return
	}

	transaction, _ := payment.ListTransactionsByID(extReq, paymnt.TransactionID)

	if strings.EqualFold(statusString, "completed") {
		extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
			AccountID:     int(paymnt.AccountID),
			TransactionID: transaction.TransactionID,
			MilestoneID:   transaction.MilestoneID,
			Status:        "closed",
		})

		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "completed"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}
		_, err = payment.DebitWallet(extReq, db, amount, disbursement.DebitCurrency, disbursement.RecipientID, "yes", transaction.TransactionID)
		if err != nil {
			if err != nil {
				extReq.Logger.Error(fmt.Sprintf("error debiting wallet for user %v, amount:%v; currency: %v; disbursementId:%v; error: %v", disbursement.RecipientID, amount, disbursement.DebitCurrency, disbursement.DisbursementID, err.Error()))
			}
		}

		if strings.EqualFold(disbursement.Type, "refund") {
			extReq.SendExternalRequest(request.SuccessfulRefundNotification, external_models.OnlyTransactionIDAndAccountIDRequest{
				TransactionID: paymnt.TransactionID,
				AccountID:     disbursement.RecipientID,
			})
		} else {
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
		}

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
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement completed", disbursement.Reference))

	} else if strings.EqualFold(statusString, "failed") {
		if tries > maxTries {
			extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
				AccountID:     int(paymnt.AccountID),
				TransactionID: transaction.TransactionID,
				MilestoneID:   transaction.MilestoneID,
				Status:        "cmdp",
			})
		}

		err = payment.SlackNotify(extReq, disbursementChannelD, `
		 	Disbursement For Transaction #`+paymnt.TransactionID+`  failed.
            Environment: `+config.GetConfig().App.Name+`
            Account ID: `+fmt.Sprintf("%v", disbursement.RecipientID)+`
            Beneficiary Name: `+disbursement.BeneficiaryName+`
            Amount: `+fmt.Sprintf("%v %v", disbursement.DebitCurrency, disbursement.Amount)+`
            Status: `+disbursement.Status+`
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement failed", disbursement.Reference))

	} else if strings.EqualFold(statusString, "ongoing") {
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement initiated/ongoing", disbursement.Reference))

	} else if strings.EqualFold(statusString, "null") {
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "failed"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement null or not found", disbursement.Reference))

	} else if strings.EqualFold(statusString, "new") {
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "new"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement is NEW", disbursement.Reference))

	} else if strings.EqualFold(statusString, "pending") {
		disbursement.PaymentReleasedAt = time.Now().Format("2006-01-02 15:04:05")
		disbursement.Status = "pending"
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement is Pending", disbursement.Reference))

	} else {
		extReq.Logger.Info(fmt.Sprintf("%v: Payment disbursement status could not be determined", disbursement.Reference))
	}

}

func decideTries(extReq request.ExternalRequest, db postgresql.Databases, disbursement *models.Disbursement) (int, bool) {
	var (
		tries    = disbursement.Tries
		baseTime time.Time
	)

	disbursement.Tries += 1
	err := disbursement.UpdateAllFields(db.Payment)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
	}

	if disbursement.TryAgainAt == baseTime {
		disbursement.TryAgainAt = time.Now().Add(24 * time.Hour)
		err := disbursement.UpdateAllFields(db.Payment)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error updating disbursement %v; error: %v", disbursement.DisbursementID, err.Error()))
		}
	} else {
		if disbursement.TryAgainAt.After(time.Now()) {
			return tries, false
		}
	}

	if tries > maxTries {
		return tries, false
	}

	return tries, true
}
