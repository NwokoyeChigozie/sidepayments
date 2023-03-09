package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func InitiatePaymentService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.InitiatePaymentRequest) (models.InitiatePaymentResponse, int, error) {
	var (
		payment                     = models.Payment{TransactionID: req.TransactionID}
		paymentGateway              = req.PaymentGateway
		reference                   = ""
		rave                        = Rave{ExtReq: extReq}
		response                    = models.InitiatePaymentResponse{}
		maxUSDAmountNigeria float64 = 100
		onlinePayment               = config.GetConfig().ONLINE_PAYMENT
	)

	transaction, err := ListTransactionsByID(extReq, req.TransactionID)
	if err != nil {
		return response, http.StatusBadRequest, err
	}

	buyerParty, ok := transaction.Parties["buyer"]
	if !ok {
		return response, http.StatusBadRequest, fmt.Errorf("buyer not found")
	}

	buyer, err := GetUserWithAccountID(extReq, buyerParty.AccountID)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("buyer account not found: %v", err.Error())
	}

	if buyer.EmailAddress == "" {
		return response, http.StatusBadRequest, fmt.Errorf("buyer does not have an email address")
	}

	code, err := payment.GetPaymentByTransactionID(db.Payment)
	if err != nil {
		return response, code, err
	}

	if payment.TotalAmount > onlinePayment.Max {
		return response, http.StatusBadRequest, fmt.Errorf("payable amount exceeds online payment max")
	}

	isNigerian, err := isRequestIPNigerian(extReq, c)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if isNigerian {
		if strings.ToUpper(transaction.Currency) == "USD" {
			if payment.TotalAmount > maxUSDAmountNigeria {
				return response, http.StatusBadRequest, fmt.Errorf("payable amount exceeds CBN's Limit. You can only pay $%v and below", maxUSDAmountNigeria)
			}
		}
	}

	_, err = GetBusinessProfileByAccountID(extReq, extReq.Logger, transaction.BusinessID)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("business profile not found: %v", err.Error())
	}

	businessCharge, err := getBusinessChargeWithBusinessIDAndCountry(extReq, transaction.BusinessID, transaction.Country.CountryCode)
	if err != nil {
		businessCharge, err = initBusinessCharge(extReq, transaction.BusinessID, transaction.Country.CurrencyCode)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}
	}

	if paymentGateway == "" {
		paymentGateway = businessCharge.PaymentGateway
	}

	if paymentGateway == "campay" {
		id, _ := uuid.NewV4()
		reference = id.String()
	} else {
		reference = fmt.Sprint("VC%v", strconv.Itoa(utility.GetRandomNumbersInRange(1000000000, 9999999999)))
	}

	if payment.IsPaid {
		return response, http.StatusBadRequest, fmt.Errorf("Transaction has been paid for.")
	}

	finalCharge := payment.EscrowCharge
	shippingFee := payment.ShippingFee
	brokerFee := payment.BrokerCharge
	var amount float64 = 0
	var totalAmount float64 = 0

	chargeBearerParty, chargeBearerPartyCheck := transaction.Parties["charge_bearer"]
	if chargeBearerPartyCheck {
		if chargeBearerParty.AccountID == buyerParty.AccountID {
			amount = payment.TotalAmount + finalCharge
			totalAmount = payment.TotalAmount + finalCharge
		} else {
			amount = payment.TotalAmount
			totalAmount = payment.TotalAmount
		}
	} else {
		amount = payment.TotalAmount
		totalAmount = payment.TotalAmount
	}

	if shippingFee != 0.00 {
		if chargeBearerPartyCheck {
			if chargeBearerParty.AccountID == buyerParty.AccountID {
				amount = totalAmount + shippingFee
			}
		}
	}

	if strings.ToLower(transaction.Type) == "broker" {
		if chargeBearerPartyCheck {
			if chargeBearerParty.AccountID == buyerParty.AccountID {
				amount = totalAmount + shippingFee + brokerFee
			}
		}
	}

	buyerUserProfile, err := GetUserProfileByAccountID(extReq, extReq.Logger, int(buyer.AccountID))
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("buyer's user profile not found, %v", err.Error())
	}

	if buyerUserProfile.Country == "" {
		buyerUserProfile.Country = "NG"
	}

	successPage, err := utility.URLDecode(req.SuccessPage)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("error decoding success_page: %v", err.Error())
	}
	utility.AddQueryParam(&successPage, "reference", reference)

	paymentRef := ""

	if paymentGateway == "" {
		paymentGateway = "rave"
	}

	if paymentGateway != "wallet" {
		if strings.ToUpper(transaction.Currency) == "USD" {
			paymentGateway = "rave"
		}
	}

	paymentUrl, paymentRequest, err := rave.InitPayment(reference, buyer.EmailAddress, transaction.Currency, successPage, amount)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("initiating payment failed: %v", err.Error())
	}

	paymentInfo := models.PaymentInfo{
		PaymentID:   payment.PaymentID,
		Reference:   reference,
		Status:      "pending",
		Gateway:     paymentGateway,
		RedirectUrl: successPage,
	}

	err = paymentInfo.CreatePaymentInfo(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	paymentRequestByte, err := json.Marshal(paymentRequest)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	paymentLog := models.PaymentLog{
		PaymentID: payment.PaymentID,
		Log:       string(paymentRequestByte),
	}

	err = paymentLog.CreatePaymentLog(db.Transaction)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if paymentGateway == "campay" {
		response = models.InitiatePaymentResponse{
			Ref:           reference,
			ExternalRef:   paymentRef,
			PaymentStatus: "initiated",
		}
	} else {
		response = models.InitiatePaymentResponse{
			Link:          utility.Stripslashes(paymentUrl),
			Ref:           reference,
			ExternalRef:   paymentRef,
			TransactionID: transaction.TransactionID,
		}
	}

	return response, http.StatusOK, nil
}

func InitiatePaymentHeadlessService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.InitiatePaymentHeadlessRequest) (models.InitiatePaymentResponse, int, error) {
	var (
		response                    = models.InitiatePaymentResponse{}
		rave                        = Rave{ExtReq: extReq}
		monnify                     = Monnify{ExtReq: extReq}
		amount              float64 = 0
		paymentGateway              = req.PaymentGateway
		country                     = ""
		successPage                 = ""
		failPage                    = ""
		paymentUrl                  = ""
		reference                   = ""
		charge              float64 = 0
		paymentRequest      interface{}
		paymentRef                  = ""
		maxUSDAmountNigeria float64 = 100
		onlinePayment               = config.GetConfig().ONLINE_PAYMENT
	)

	if req.Amount > onlinePayment.Max {
		return response, http.StatusBadRequest, fmt.Errorf("payable amount exceeds online payment max")
	}

	isNigerian, err := isRequestIPNigerian(extReq, c)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if isNigerian {
		if strings.ToUpper(req.Currency) == "USD" {
			if req.Amount > maxUSDAmountNigeria {
				return response, http.StatusBadRequest, fmt.Errorf("payable amount exceeds CBN's Limit. You can only pay $%v and below", maxUSDAmountNigeria)
			}
		}
	}
	accessToken, err := GetAccessTokenByKeyFromRequest(extReq, c)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if req.Initialize {
		amount = 1
	} else {
		chargeObj, err := GetEscrowCharge(extReq, accessToken.AccountID, req.Amount)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}
		charge = chargeObj.Charge
		amount = req.Amount + charge
	}

	if paymentGateway == "campay" {
		id, _ := uuid.NewV4()
		reference = id.String()
	} else {
		reference = fmt.Sprint("VC%v", strconv.Itoa(utility.GetRandomNumbersInRange(1000000000, 9999999999)))
	}

	businessProfile, err := GetBusinessProfileByAccountID(extReq, extReq.Logger, accessToken.AccountID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	businessCharge, err := getBusinessChargeWithBusinessIDAndCountry(extReq, businessProfile.AccountID, strings.ToUpper(req.Currency))
	if err != nil {
		businessCharge, err = initBusinessCharge(extReq, businessProfile.AccountID, strings.ToUpper(req.Currency))
		if err != nil {
			return response, http.StatusInternalServerError, err
		}
	}

	if paymentGateway == "" {
		paymentGateway = businessCharge.PaymentGateway
	}

	if req.Country != "" {
		_, err := GetCountryByNameOrCode(extReq, extReq.Logger, strings.ToUpper(req.Country))
		if err != nil {
			return response, http.StatusBadRequest, fmt.Errorf("country does not exist")
		}
		country = strings.ToUpper(req.Country)
	} else {
		country = businessCharge.Country
	}

	if req.SuccessUrl == "" {
		successPage = utility.GenerateGroupByURL(c, "pay/successful", map[string]string{"website": businessProfile.Website})
	} else {
		successPage = req.SuccessUrl
	}

	if req.FailUrl == "" {
		failPage = utility.GenerateGroupByURL(c, "pay/failed", map[string]string{"website": businessProfile.Website})
	} else {
		failPage = req.FailUrl
	}

	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if strings.ToUpper(req.Currency) == "USD" {
		req.PaymentGateway = "rave"
	}

	callback := utility.GenerateGroupByURL(c, "pay/status", map[string]string{"reference": reference, "success_page": successPage, "failure_page": failPage, "fund_wallet": req.FundWallet})

	switch strings.ToLower(paymentGateway) {
	case "rave":
		url, reqData, err := rave.InitPayment(reference, user.EmailAddress, strings.ToUpper(req.Currency), callback, amount)
		if err != nil {
			return response, http.StatusInternalServerError, fmt.Errorf("initiating payment failed: %v", err.Error())
		}
		paymentUrl, paymentRequest = url, reqData
	case "monnify":
		url, reqData, err := monnify.InitPayment(amount, fmt.Sprintf("%v %v", user.Lastname, user.Firstname), user.EmailAddress, reference, "Payment For Vesicash", req.Currency, callback)
		if err != nil {
			return response, http.StatusInternalServerError, fmt.Errorf("initiating payment failed: %v", err.Error())
		}
		paymentUrl, paymentRequest = url, reqData
	default:
		url, reqData, err := rave.InitPayment(reference, user.EmailAddress, strings.ToUpper(req.Currency), callback, amount)
		if err != nil {
			return response, http.StatusInternalServerError, fmt.Errorf("initiating payment failed: %v", err.Error())
		}
		paymentUrl, paymentRequest = url, reqData
	}

	payment := models.Payment{
		PaymentID:    utility.RandomString(10),
		TotalAmount:  req.Amount,
		EscrowCharge: charge,
		IsPaid:       false,
		AccountID:    int64(req.AccountID),
		BusinessID:   int64(businessProfile.AccountID),
		Currency:     strings.ToUpper(req.Currency),
	}

	err = payment.CreatePayment(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	paymentInfo := models.PaymentInfo{
		PaymentID:   payment.PaymentID,
		Reference:   reference,
		Status:      "pending",
		Gateway:     paymentGateway,
		RedirectUrl: successPage,
		FailUrl:     failPage,
	}

	err = paymentInfo.CreatePaymentInfo(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	paymentRequestByte, err := json.Marshal(paymentRequest)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	paymentLog := models.PaymentLog{
		PaymentID: payment.PaymentID,
		Log:       string(paymentRequestByte),
	}

	err = paymentLog.CreatePaymentLog(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if paymentGateway == "campay" {
		response = models.InitiatePaymentResponse{
			Ref:           reference,
			ExternalRef:   paymentRef,
			PaymentStatus: "initiated",
		}
	} else {
		response = models.InitiatePaymentResponse{
			Link:        utility.Stripslashes(paymentUrl),
			Ref:         reference,
			ExternalRef: paymentRef,
		}
	}
	extReq.Logger.Info(country)

	return response, http.StatusOK, nil
}

func FundWalletService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.FundWalletRequest) (map[string]interface{}, int, error) {
	var (
		response        = map[string]interface{}{}
		rave            = Rave{ExtReq: extReq}
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)
	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	businessCharge, err := getBusinessChargeWithBusinessIDAndCountry(extReq, req.AccountID, req.Currency)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("Missing business charges data, visit https://vesicash.com to set this up.")
	}

	chargeString := businessCharge.ProcessingFee
	charge, _ := strconv.ParseFloat(chargeString, 64)
	finalAmount := req.Amount + float64(charge)

	if user.Firstname == "" || user.Lastname == "" {
		return response, http.StatusBadRequest, fmt.Errorf("Missing firstname or lastname in user data.")
	}

	referenceID, _ := uuid.NewV4()
	reference := "VESICASH_FA_" + referenceID.String()
	accountName := fmt.Sprintf("%v %v", user.Firstname, user.Lastname)

	payment := models.Payment{
		PaymentID:   utility.RandomString(10),
		TotalAmount: req.Amount,
		IsPaid:      false,
		AccountID:   int64(req.AccountID),
		BusinessID:  int64(req.AccountID),
		Currency:    strings.ToUpper(req.Currency),
	}

	err = payment.CreatePayment(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	paymentInfo := models.PaymentInfo{
		PaymentID: payment.PaymentID,
		Reference: reference,
		Status:    "pending",
		Gateway:   "flutterwave",
	}

	err = paymentInfo.CreatePaymentInfo(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	accountDetails, err := rave.ReserveAccount(reference, accountName, user.EmailAddress, user.Firstname, user.Lastname, finalAmount)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	bank, err := rave.GetBank("NGN", accountDetails.BankName)
	bankCode := "0"
	if err == nil {
		bankCode = bank.Code
	}

	fID, _ := uuid.NewV4()
	fundingAccount := models.FundingAccount{
		FundingAccountID:  fID.String(),
		Reference:         reference,
		AccountID:         req.AccountID,
		AccountName:       accountName,
		AccountNumber:     accountDetails.AccountNumber,
		BankCode:          bankCode,
		BankName:          accountDetails.BankName,
		LastFundingAmount: fmt.Sprintf("%v", finalAmount),
		EscrowWallet:      req.Escrow_wallet,
	}

	err = fundingAccount.CreateFundingAccount(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	err = SlackNotify(paymentChannelD, `
		Funding Account Generated
		Environment: `+config.GetConfig().App.Name+`
		Account ID: `+strconv.Itoa(req.AccountID)+`
		Account Number: `+accountDetails.AccountNumber+`
		Account Name: `+accountName+`
		Bank: `+accountDetails.BankName+`
		Status: SUCCESSFUL
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	pendingTransferFunding := models.PendingTransferFunding{
		Reference: fundingAccount.Reference,
		Status:    "pending",
		Type:      "funding",
	}

	err = pendingTransferFunding.CreatePendingTransferFunding(db.Payment)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	return gin.H{"funding_account": accountDetails, "currency": req.Currency, "amount": finalAmount, "payment": payment, "reference": reference}, http.StatusOK, nil

}

func FundWalletVerifyService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, reference, currency string) (map[string]interface{}, int, error) {
	var (
		fundingAccount  = models.FundingAccount{Reference: reference}
		response        = map[string]interface{}{}
		rave            = Rave{ExtReq: extReq}
		paymentChannelD = config.GetConfig().Slack.PaymentChannelID
	)

	code, err := fundingAccount.GetFundingAccountByReference(db.Payment)
	if err != nil {
		if code == http.StatusInternalServerError {
			return response, code, err
		}
		return response, code, fmt.Errorf("No funding account attached to this reference ID. Kindly re-generate a reference using the fund wallet endpoint")
	}

	user, err := GetUserWithAccountID(extReq, fundingAccount.AccountID)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("Funding account user could not be found")
	}

	businessCharge, _ := getBusinessChargeWithBusinessIDAndCurrency(extReq, fundingAccount.AccountID, strings.ToUpper(currency))

	processingFee, err := strconv.ParseFloat(businessCharge.ProcessingFee, 64)
	lastFundingAmount, err := strconv.ParseFloat(fundingAccount.LastFundingAmount, 64)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("last funding amount is not of float value")
	}

	amount := lastFundingAmount - processingFee
	if amount < 0 {
		return response, http.StatusInternalServerError, fmt.Errorf("amount is less than 0")
	}

	_, status, _ := rave.VerifyTransactionByTxRef(reference)
	if status == "pending" {
		return response, http.StatusBadRequest, fmt.Errorf("account has not received any payments or error verifying payment")
	}

	businessProfile, _ := GetBusinessProfileByAccountID(extReq, extReq.Logger, int(user.AccountID))

	if status != "success" {
		if businessProfile.Webhook_uri != "" {
			InitWebhook(extReq, db, businessProfile.Webhook_uri, "bank-transfer.failed", map[string]interface{}{"reference": reference, "amount": amount, "status": "failed"}, businessProfile.AccountID)
		}

	}

	err = CreditWallet(extReq, db, amount, currency, int(user.AccountID), false, fundingAccount.EscrowWallet, "")
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("error crediting wallet: %v", err.Error())
	}

	err = SlackNotify(paymentChannelD, `
		Wallet Funding For Customer #`+strconv.Itoa(int(user.AccountID))+`
		Environment: `+config.GetConfig().App.Name+`
		Account ID: `+strconv.Itoa(int(user.AccountID))+`
		Beneficiary Name: `+user.Firstname+` `+user.Lastname+`
		Amount: `+fmt.Sprintf("%v", amount)+`
		Status: PAID
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	err = SlackNotify(paymentChannelD, `
		Bank Transfer Payment
		Environment: `+config.GetConfig().App.Name+`
		Reference: `+reference+`
		Account Number: `+fundingAccount.AccountNumber+`
		Account Name: `+fundingAccount.AccountName+`
		Bank: `+fundingAccount.BankName+`
		Amount: `+fmt.Sprintf("%v", amount)+`
		Status: SUCCESSFUL
			`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	if businessProfile.Webhook_uri != "" {
		InitWebhook(extReq, db, businessProfile.Webhook_uri, "bank-transfer.success", map[string]interface{}{"reference": reference, "amount": amount, "status": "success"}, businessProfile.AccountID)
	}

	return gin.H{"reference": reference, "amount": amount, "currency": strings.ToUpper(currency)}, http.StatusOK, nil
}

func FundWalletLogsService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, accountID int, paginator postgresql.Pagination) ([]models.FundingAccount, postgresql.PaginationResponse, int, error) {
	var (
		fundingAccount = models.FundingAccount{AccountID: accountID}
	)

	fundingAccounts, pagination, err := fundingAccount.GetFundingAccountsByAccountID(db.Payment, "id", "desc", paginator)
	if err != nil {
		return fundingAccounts, pagination, http.StatusInternalServerError, err
	}

	return fundingAccounts, pagination, http.StatusOK, nil
}
