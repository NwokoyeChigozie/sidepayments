package payment

import (
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

func RaveWebhookService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.RaveWebhookRequest, requestBody []byte) (int, error) {

	if config.GetConfig().Rave.WebhookSecret != utility.GetHeader(c, "verif-hash") {
		return http.StatusUnauthorized, fmt.Errorf("invalid webhook secret")
	}

	extReq.Logger.Info("rave webhhook log info", string(requestBody))

	webhookLog := models.WebhookLog{
		Log:      string(requestBody),
		Provider: "flutterwave",
	}
	err := webhookLog.CreateWebhookLog(db.Payment)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var (
		eventType       = req.Event
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)

	err = SlackNotify(extReq, paymentChannelD, `
					Web Hook Received RAVE
					Environment: `+config.GetConfig().App.Name+`
					Event: `+eventType+`
					Reference: `+eventType+`
					Status: SUCCESSFUL
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	go func(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.RaveWebhookRequest, eventType string) {
		var (
			err     error
			errCode = http.StatusOK
		)
		switch eventType {
		case "charge.completed":
			errCode, err = handleChargeCompleted(c, extReq, db, req)
		case "Transfer":
			errCode, err = handleTransfer(c, extReq, db, req)
		case "transfer.completed":
			errCode, err = handleTransferCompleted(c, extReq, db, req)
		default:
			err, errCode = fmt.Errorf("event type %v, not implemented", eventType), http.StatusNotImplemented
		}

		if err != nil {
			extReq.Logger.Error("rave webhhook log error: ", err.Error(), "error code:", errCode)
		}
	}(c, extReq, db, req, eventType)

	return http.StatusOK, nil
}

func handleChargeCompleted(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.RaveWebhookRequest) (int, error) {
	var (
		data            models.RaveWebhookRequestData
		ref             string
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
		amount          float64
		currency        string
		rave            = Rave{ExtReq: extReq}
	)

	if req.Data != nil {
		data = *req.Data
	} else {
		return http.StatusBadRequest, fmt.Errorf("data not found")
	}

	if data.TxRef != nil {
		ref = *data.TxRef
	} else {
		return http.StatusBadRequest, fmt.Errorf("data txref not found")
	}

	if data.Amount != nil {
		amount = *data.Amount
	}

	if data.Currency != nil {
		currency = *data.Currency
	}

	sts, err := rave.VerifyTrans(ref, amount, currency)
	extReq.Logger.Info(fmt.Sprintf("checking issue: %v, %v, %v", sts, err, amount))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("payment verification failed")
	}

	paymentAccount := models.PaymentAccount{PaymentAccountID: ref}
	code, err := paymentAccount.GetPaymentAccountByPaymentAccountID(db.Payment)
	if err != nil {
		return code, fmt.Errorf("payment account not found: %v", err.Error())
	}

	payment := models.Payment{PaymentID: paymentAccount.PaymentID}
	payment.GetPaymentByPaymentID(db.Payment)

	transaction, _ := ListTransactionsByID(extReq, payment.TransactionID)

	if payment.IsPaid {
		return http.StatusOK, nil
	}
	escrowChargeBearerParty := transaction.Parties["charge_bearer"]
	if sts == "success" {
		if payment.TransactionID != "" {
			transactionPaid(extReq, db, &payment, &transaction, ref, currency, "card_payment")
		} else {
			payment.IsPaid = true
			payment.WalletFunded = strings.ToUpper(currency)
			payment.PaymentMethod = "card_payment"
			payment.PaymentMadeAt = time.Now()
			payment.UpdateAllFields(db.Payment)
		}
		err = SlackNotify(extReq, paymentChannelD, `
			[WEBHOOK RAVE] Card Payment
			Environment: `+config.GetConfig().App.Name+`
			Reference: `+ref+`
			Transaction ID: `+payment.TransactionID+`
			Payment ID: `+payment.PaymentID+`
			Account Number: `+paymentAccount.AccountNumber+`
			Account Name: `+paymentAccount.AccountName+`
			Bank: `+paymentAccount.BankName+`
			Amount: `+fmt.Sprintf("%v %v", currency, amount)+`
			Escrow Charge: `+fmt.Sprintf("%v", payment.EscrowCharge)+`
			Escrow Charge Bearer: `+fmt.Sprintf("%v", escrowChargeBearerParty.AccountID)+`
			Status: SUCCESSFUL
			`)
		if err != nil && !extReq.Test {
			extReq.Logger.Error("error sending notification to slack: ", err.Error())
		}
	}

	return http.StatusOK, nil
}
func handleTransfer(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.RaveWebhookRequest) (int, error) {
	var (
		transfer        models.RaveWebhookRequestTransfer
		meta            models.RaveWebhookRequestTransferMeta
		ref             string
		amount          float64
		currency        string
		merchantID      string
		transaferStatus string
		paymentMadeAt   time.Time
	)

	if req.Transfer != nil {
		transfer = *req.Transfer
	} else {
		return http.StatusInternalServerError, fmt.Errorf("transfer data not found")
	}

	if transfer.Reference != nil {
		ref = *transfer.Reference
	} else {
		return http.StatusBadRequest, fmt.Errorf("transfer reference not found")
	}

	if transfer.Meta != nil {
		meta = *transfer.Meta
	} else {
		return http.StatusInternalServerError, fmt.Errorf("transfer meta not found")
	}

	if transfer.Amount != nil {
		amount = *transfer.Amount
	}

	if transfer.Currency != nil {
		currency = strings.ToUpper(*transfer.Currency)
	}

	if meta.MerchantId != nil {
		merchantID = *meta.MerchantId
	}

	if transfer.Status != nil {
		transaferStatus = *transfer.Status
	}

	if transfer.DateCreated != nil {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", *transfer.DateCreated)
		if err != nil {
			extReq.Logger.Error("rave webhhook log error", "error parsing transfer.DateCreated", *transfer.DateCreated, err.Error())
		}
		paymentMadeAt = t
	}

	business, err := GetBusinessProfileByFlutterwaveMerchantID(extReq, extReq.Logger, merchantID)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("user with this merchant id not found")
	}
	businessID := business.AccountID

	paymentAccount := models.PaymentAccount{BusinessID: strconv.Itoa(businessID), BankCode: "flutterwave", BankName: "flutterwave"}
	_, err = paymentAccount.GetBybusinessIDBankNameAndCodeAndTransactionIDNotNull(db.Payment)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("can't find associated transaction:%v", err.Error())
	}
	transaction, err := ListTransactionsByID(extReq, paymentAccount.TransactionID)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("transaction record not found")
	}

	paymentAccount.PaymentReference = ref
	paymentAccount.UpdateAllFields(db.Payment)

	var fundingCharge float64 = 0
	if amount < 500000 {
		fundingCharge = 500
	}

	if amount > 500000 && amount < 1000000 {
		fundingCharge = 1000
	}
	if amount > 1000000 {
		fundingCharge = 2000
	}

	if strings.EqualFold(transaferStatus, "SUCCESSFUL") {
		_, err = CreditWallet(extReq, db, amount, currency, businessID, false, thisOrThatStr(transaction.EscrowWallet, "yes"), transaction.TransactionID)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		payment := models.Payment{
			PaymentID:     utility.RandomString(10),
			TotalAmount:   amount,
			EscrowCharge:  fundingCharge,
			IsPaid:        true,
			AccountID:     int64(businessID),
			BusinessID:    int64(businessID),
			Currency:      currency,
			PaymentMadeAt: paymentMadeAt,
			PaymentMethod: "wallet-transfer",
			WalletFunded:  currency,
		}
		payment.CreatePayment(db.Payment)

	}
	return http.StatusOK, nil
}
func handleTransferCompleted(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.RaveWebhookRequest) (int, error) {
	var (
		data           models.RaveWebhookRequestData
		accountNumber  string
		ref            string
		amount         float64
		currency       string
		transferStatus string
	)

	if req.Data != nil {
		data = *req.Data
	} else {
		return http.StatusBadRequest, fmt.Errorf("data not found")
	}

	if data.Reference != nil {
		ref = *data.Reference
	} else {
		return http.StatusBadRequest, fmt.Errorf("data reference not found")
	}

	if data.AccountNumber != nil {
		accountNumber = *data.AccountNumber
	}
	if data.Amount != nil {
		amount = *data.Amount
	}

	if data.Currency != nil {
		currency = *data.Currency
	}

	if data.Status != nil {
		transferStatus = *data.Status
	}

	fundingAccounts := models.FundingAccount{AccountNumber: accountNumber}
	code, err := fundingAccounts.GetFundingAccountByAccountNumber(db.Payment)
	if err != nil {
		return code, fmt.Errorf("funding account does not exist: %v", err.Error())
	}
	accountID := fundingAccounts.AccountID
	fundingAccounts.LastFundingReference = ref
	fundingAccounts.UpdateAllFields(db.Payment)

	if strings.EqualFold(transferStatus, "SUCCESSFUL") {
		_, err = CreditWallet(extReq, db, amount, currency, accountID, false, thisOrThatStr(fundingAccounts.EscrowWallet, "no"), "")
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func transactionPaid(extReq request.ExternalRequest, db postgresql.Databases, payment *models.Payment, transaction *external_models.TransactionByID, txref, currency string, paymentMethod string) error {
	var (
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)
	payment.IsPaid = true
	payment.WalletFunded = strings.ToUpper(currency)
	if transaction.EscrowWallet == "yes" {
		payment.WalletFunded = fmt.Sprintf("ESCROW_%v", currency)
	}
	payment.PaymentMethod = thisOrThatStr(paymentMethod, "card_payment")
	payment.PaymentMadeAt = time.Now()
	err := payment.UpdateAllFields(db.Payment)
	if err != nil {
		return err
	}

	extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
		AccountID:     transaction.BusinessID,
		TransactionID: transaction.TransactionID,
		MilestoneID:   transaction.MilestoneID,
		Status:        "ip",
	})

	var amounts float64
	milestones := transaction.Milestones
	for _, v := range milestones {
		amounts += v.Amount
	}

	if amounts == transaction.TotalAmount {
		for _, v := range milestones {
			extReq.SendExternalRequest(request.TransactionUpdateStatus, external_models.UpdateTransactionStatusRequest{
				AccountID:     transaction.BusinessID,
				TransactionID: transaction.TransactionID,
				MilestoneID:   v.MilestoneID,
				Status:        "ip",
			})
		}
	}

	if transaction.Source == "transfer" {
		extReq.SendExternalRequest(request.BuyerSatisfied, external_models.OnlyTransactionIDRequiredRequest{
			TransactionID: payment.TransactionID,
		})
	}

	buyerParty := transaction.Parties["buyer"]
	sellerParty := transaction.Parties["seller"]
	brokerParty := transaction.Parties["broker"]
	chargeBearerParty := transaction.Parties["charge_bearer"]
	shippingChargeBearerParty := transaction.Parties["shipping_charge_bearer"]
	extReq.SendExternalRequest(request.TransactionPaidNotification, external_models.OnlyTransactionIDAndAccountIDRequest{
		TransactionID: payment.TransactionID,
		AccountID:     buyerParty.AccountID,
	})

	if transaction.Type == "broker" {
		extReq.SendExternalRequest(request.TransactionPaidNotification, external_models.OnlyTransactionIDAndAccountIDRequest{
			TransactionID: payment.TransactionID,
			AccountID:     brokerParty.AccountID,
		})
	}

	inspectionPeriod, _ := strconv.Atoi(transaction.InspectionPeriod)
	inspectionPeriodAsDate := ""
	if inspectionPeriod > 0 {
		t := time.Unix(int64(inspectionPeriod), 0)
		inspectionPeriodAsDate = t.Format("2006-01-02")
	}
	extReq.SendExternalRequest(request.PaymentInvoiceNotification, external_models.PaymentInvoiceNotificationRequest{
		Reference:                 txref,
		PaymentID:                 payment.PaymentID,
		TransactionType:           transaction.Type,
		TransactionID:             payment.TransactionID,
		Buyer:                     buyerParty.AccountID,
		Seller:                    sellerParty.AccountID,
		InspectionPeriodFormatted: inspectionPeriodAsDate,
		ExpectedDelivery:          transaction.DueDate,
		Title:                     transaction.Title,
		Currency:                  transaction.Currency,
		Amount:                    transaction.TotalAmount,
		EscrowCharge:              payment.EscrowCharge,
		BrokerCharge:              payment.BrokerCharge,
	})

	if chargeBearerParty.AccountID == sellerParty.AccountID {
		payment.TotalAmount = payment.TotalAmount - payment.EscrowCharge
		payment.UpdateAllFields(db.Payment)
	}

	if transaction.Type == "broker" && (chargeBearerParty.AccountID == brokerParty.AccountID) {
		payment.BrokerCharge = payment.BrokerCharge - payment.EscrowCharge
		if payment.BrokerCharge < 0 {
			payment.BrokerCharge = 0
		}
		payment.UpdateAllFields(db.Payment)
	}

	if payment.ShippingFee != 0 && (shippingChargeBearerParty.AccountID == sellerParty.AccountID) {
		payment.TotalAmount = payment.TotalAmount - payment.ShippingFee
		payment.UpdateAllFields(db.Payment)
	}

	businessEscrowCharge, err := GetBusinessChargeWithBusinessIDAndCountry(extReq, transaction.BusinessID, transaction.Country.CountryCode)
	if err == nil {
		vesicashCharge, _ := strconv.ParseFloat(businessEscrowCharge.VesicashCharge, 64)
		businessPerc, _ := strconv.ParseFloat(businessEscrowCharge.BusinessCharge, 64)
		_, err = CreditWallet(extReq, db, utility.PercentageOf(payment.TotalAmount, vesicashCharge), transaction.Currency, 1, false, "no", transaction.TransactionID)
		if err != nil {
			return err
		}

		_, err = CreditWallet(extReq, db, utility.PercentageOf(payment.TotalAmount, businessPerc), transaction.Currency, transaction.BusinessID, false, transaction.EscrowWallet, transaction.TransactionID)
		if err != nil {
			return err
		}
	}

	if payment.ShippingFee != 0 && shippingChargeBearerParty.AccountID != 0 {
		_, err = CreditWallet(extReq, db, payment.ShippingFee, transaction.Currency, shippingChargeBearerParty.AccountID, false, transaction.EscrowWallet, transaction.TransactionID)
		if err != nil {
			return err
		}
	}

	err = SlackNotify(extReq, paymentChannelD, `
			[WEBHOOK RAVE] Payment Status For Transaction #`+payment.PaymentID+`
			Environment: `+config.GetConfig().App.Name+`
			Payment ID: `+payment.PaymentID+`
			Amount: `+fmt.Sprintf("%v %v", payment.Currency, payment.TotalAmount)+`
			Status: SUCCESSFUL
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	return nil
}
