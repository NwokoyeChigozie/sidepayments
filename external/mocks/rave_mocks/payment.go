package rave_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func RaveInitPayment(logger *utility.Logger, idata interface{}) (external_models.RaveInitPaymentResponse, error) {

	var (
		outBoundResponse external_models.RaveInitPaymentResponse
	)

	data, ok := idata.(external_models.RaveInitPaymentRequest)
	if !ok {
		logger.Error("init payment rave", idata, "request data format error")
		return external_models.RaveInitPaymentResponse{}, fmt.Errorf("request data format error")
	}

	logger.Info("init payment rave", outBoundResponse, data)

	return external_models.RaveInitPaymentResponse{
		Status:  "success",
		Message: "success",
		Data: struct {
			Link string "json:\"link\""
		}{
			Link: "https://rave.payment.link",
		},
	}, nil
}

func RaveReserveAccount(logger *utility.Logger, idata interface{}) (external_models.RaveReserveAccountResponseData, error) {

	var (
		outBoundResponse external_models.RaveReserveAccountResponse
	)

	data, ok := idata.(external_models.RaveReserveAccountRequest)
	if !ok {
		logger.Error("rave reserve account", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("rave reserve account", outBoundResponse, data)
	return external_models.RaveReserveAccountResponseData{
		ResponseCode:    "00",
		ResponseMessage: "Transaction in progress",
		FlwRef:          data.TxRef,
		OrderRef:        data.TxRef,
		AccountNumber:   "1357958081",
		Frequency:       nil,
		BankName:        "vesicash bank",
		Note:            "Please make a bank transfer to Earth Gang",
		Amount:          nil,
	}, nil
}

func RaveVerifyTransactionByTxRef(logger *utility.Logger, idata interface{}) (external_models.RaveVerifyTransactionResponseData, error) {

	var (
		outBoundResponse external_models.RaveVerifyTransactionResponse
	)

	data, ok := idata.(string)
	if !ok {
		logger.Error("rave verify transaction by tx_ref", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("rave verify transaction by tx_ref", outBoundResponse)

	return external_models.RaveVerifyTransactionResponseData{
		ID:            12,
		TxRef:         data,
		FlwRef:        data + "yvyv",
		Amount:        200,
		ChargedAmount: 200,
		Currency:      "NGN",
		AmountSettled: 200,
		Card: &external_models.RaveVerifyTransactionResponseDataCard{
			First6digits: "553188",
			Last4digits:  "2950",
			Type:         "MASTERCARD",
			Token:        "flw-t1nf-f9b3bf384cd30d6fca42b6df9d27bd2f-m03k",
			Expiry:       "09/22",
			Country:      "NIGERIA NG",
		},
		Status: "successful",
	}, nil
}

func RaveChargeCard(logger *utility.Logger, idata interface{}) (external_models.RaveVerifyTransactionResponseData, error) {

	var (
		outBoundResponse external_models.RaveVerifyTransactionResponse
	)

	data, ok := idata.(external_models.RaveChargeCardRequest)
	if !ok {
		logger.Error("rave charge card", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("rave charge card", outBoundResponse, data)

	return external_models.RaveVerifyTransactionResponseData{
		ID:            22,
		TxRef:         data.TxRef,
		Amount:        data.Amount,
		AmountSettled: data.Amount,
		ChargedAmount: data.Amount,
		Currency:      data.Currency,
		Card: &external_models.RaveVerifyTransactionResponseDataCard{
			First6digits: "553188",
			Last4digits:  "2950",
			Type:         "MASTERCARD",
			Token:        "flw-t1nf-f9b3bf384cd30d6fca42b6df9d27bd2f-m03k",
			Expiry:       "09/22",
			Country:      "NIGERIA NG",
		},
	}, nil
}
