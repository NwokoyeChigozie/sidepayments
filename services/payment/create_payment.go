package payment

import (
	"fmt"
	"net/http"
	"time"

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

func CreatePaymentHeadlessService(extReq request.ExternalRequest, db postgresql.Databases, req models.CreatePaymentHeadlessRequest) (models.Payment, int, error) {

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	payment := models.Payment{
		PaymentID:    utility.RandomString(10),
		AccountID:    int64(req.AccountID),
		TotalAmount:  req.TotalAmount,
		EscrowCharge: req.EscrowCharge,
		IsPaid:       false,
		Currency:     req.Currency,
	}

	if req.PaymentMadeAt != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.PaymentMadeAt)
		if err != nil {
			return models.Payment{}, http.StatusBadRequest, fmt.Errorf("incorrect format for  payment_made_at, try 2023-11-23 15:04:05")
		}
		payment.PaymentMadeAt = t
	}
	err := payment.CreatePayment(db.Payment)
	if err != nil {
		return payment, http.StatusInternalServerError, err
	}

	return payment, http.StatusCreated, nil

}

func EditPaymentService(extReq request.ExternalRequest, db postgresql.Databases, req models.EditPaymentRequest, user external_models.User) (int, error) {

	var (
		payment = models.Payment{PaymentID: req.PaymentID}
	)

	code, err := payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		return code, err
	}

	if payment.AccountID != 0 && payment.AccountID != int64(user.AccountID) {
		return http.StatusUnauthorized, fmt.Errorf("not allowed to edit payment")
	}

	payment.EscrowCharge = req.EscrowCharge
	err = payment.UpdateAllFields(db.Payment)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func VerifyTransactionPaymentService(extReq request.ExternalRequest, db postgresql.Databases, transactionID string) (models.VerifyTransactionPaymentResponse, int, error) {

	var (
		payment  = models.Payment{TransactionID: transactionID}
		response = models.VerifyTransactionPaymentResponse{}
	)

	code, err := payment.GetPaymentByTransactionIDAndNotPaymentMadeAt(db.Payment)
	if err != nil {
		return response, code, err
	}

	if payment.IsPaid {
		response.Status = "success"
		response.IsPaid = true
		response.Amount = payment.TotalAmount
		response.Charge = payment.EscrowCharge
		response.Date = payment.PaymentMadeAt
	} else {
		response.Status = "failed"
		response.IsPaid = false
		response.Amount = payment.TotalAmount
		response.Charge = payment.EscrowCharge
		response.Date = payment.PaymentMadeAt
	}

	return response, http.StatusOK, nil
}

func DeletePaymentService(extReq request.ExternalRequest, db postgresql.Databases, paymentID string, user external_models.User) (int, error) {

	var (
		payment = models.Payment{PaymentID: paymentID}
	)

	code, err := payment.GetPaymentByPaymentID(db.Payment)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
		return http.StatusOK, nil
	}

	if payment.AccountID != 0 && payment.AccountID != int64(user.AccountID) {
		return http.StatusUnauthorized, fmt.Errorf("not allowed to delete payment")
	}

	err = payment.Delete(db.Payment)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
