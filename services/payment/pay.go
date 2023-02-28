package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func InitiatePaymentService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, req models.InitiatePaymentRequest, user external_models.User) (models.InitiatePaymentResponse, int, error) {
	var (
		payment                     = models.Payment{TransactionID: req.TransactionID}
		onlinePayment               = config.GetConfig().ONLINE_PAYMENT
		maxUSDAmountNigeria float64 = 100
		paymentGateway              = req.PaymentGateway
		reference                   = ""
		rave                        = Rave{ExtReq: extReq}
		response                    = models.InitiatePaymentResponse{}
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
			PaymentStatus: "initiated",
			TransactionID: transaction.TransactionID,
		}
	}

	return response, http.StatusOK, nil
}
