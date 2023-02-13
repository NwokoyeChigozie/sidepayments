package payment

import (
	"fmt"
	"net/http"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/payment-ms/utility"
)

func CreatePaymentService(extReq request.ExternalRequest, db postgresql.Databases, req models.CreatePaymentRequest, user external_models.User) (models.Payment, int, error) {
	var (
		disburseCurrency = ""
	)
	transaction, err := ListTransactionsByID(extReq, req.TransactionID)
	if err != nil {
		return models.Payment{}, http.StatusBadRequest, err
	}

	sellerParty, ok := transaction.Parties["seller"]
	if !ok {
		return models.Payment{}, http.StatusBadRequest, fmt.Errorf("seller not found")
	}
	buyerParty := transaction.Parties["buyer"]

	receiver, err := GetUserWithAccountID(extReq, sellerParty.AccountID)
	if err == nil {
		receiverProfile, err := GetUserProfileByAccountID(extReq, extReq.Logger, int(receiver.AccountID))
		if err == nil {
			countryCode := receiverProfile.Country
			if countryCode == "" {
				countryCode = "NG"
			}

			country, err := GetCountryByNameOrCode(extReq, extReq.Logger, countryCode)
			if err != nil {
				extReq.Logger.Info("error getting country with country code", countryCode, err)
				return models.Payment{}, http.StatusBadRequest, err
			}
			disburseCurrency = country.CurrencyCode
		}
	}

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	payment := models.Payment{
		PaymentID:        utility.RandomString(10),
		AccountID:        int64(buyerParty.AccountID),
		TransactionID:    req.TransactionID,
		TotalAmount:      req.TotalAmount,
		EscrowCharge:     req.EscrowCharge,
		ShippingFee:      req.ShippingFee,
		BrokerCharge:     req.BrokerCharge,
		IsPaid:           false,
		Currency:         req.Currency,
		DisburseCurrency: disburseCurrency,
	}
	err = payment.CreatePayment(db.Payment)
	if err != nil {
		return payment, http.StatusInternalServerError, err
	}

	return payment, http.StatusCreated, nil

}
