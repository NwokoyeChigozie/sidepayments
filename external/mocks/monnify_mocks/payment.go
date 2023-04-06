package monnify_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func MonnifyInitPayment(logger *utility.Logger, idata interface{}) (external_models.MonnifyInitPaymentResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyInitPaymentResponse
	)

	data, ok := idata.(external_models.MonnifyInitPaymentRequest)
	if !ok {
		logger.Info("monnify init payment", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	return external_models.MonnifyInitPaymentResponseBody{
		TransactionReference: data.PaymentReference,
		PaymentReference:     data.PaymentReference,
		CheckoutUrl:          "https://monnify.payment.url",
	}, nil
}

func MonnifyVerifyTransactionByReference(logger *utility.Logger, idata interface{}) (external_models.MonnifyVerifyByReferenceResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyVerifyByReferenceResponse
	)

	data, ok := idata.(string)
	if !ok {
		logger.Info("monnify verify transaction by reference", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	return external_models.MonnifyVerifyByReferenceResponseBody{
		CreatedOn:            "2023-01-09T18:52:45.000+0000",
		Amount:               200,
		CurrencyCode:         "NGN",
		CustomerName:         "test",
		CustomerEmail:        "test@gmail.com",
		PaymentStatus:        "PAID",
		TransactionReference: data,
		PaymentReference:     data,
	}, nil
}

func MonnifyReserveAccount(logger *utility.Logger, idata interface{}) (external_models.MonnifyReserveAccountResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyReserveAccountResponse
	)

	data, ok := idata.(external_models.MonnifyReserveAccountRequest)
	if !ok {
		logger.Info("monnify reserve account", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	return external_models.MonnifyReserveAccountResponseBody{
		AccountNumber:        "7727632865",
		AccountName:          "test",
		BankName:             "vesicash bank",
		BankCode:             "221",
		ReservationReference: data.AccountReference,
		CustomerEmail:        data.CustomerEmail,
		CurrencyCode:         data.CurrencyCode,
	}, nil
}

func GetMonnifyReserveAccountTransactions(logger *utility.Logger, idata interface{}) (external_models.GetMonnifyReserveAccountTransactionsResponseBody, error) {

	var (
		outBoundResponse external_models.GetMonnifyReserveAccountTransactionsResponse
	)

	data, ok := idata.(string)
	if !ok {
		logger.Info("get monnify reserve account transactions", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	return external_models.GetMonnifyReserveAccountTransactionsResponseBody{
		Content: []external_models.GetMonnifyReserveAccountTransactionsResponseBodyContent{
			{
				PaymentReference:     data,
				TransactionReference: data,
				Amount:               200,
				AmountPaid:           200,
			},
		},
	}, nil
}
