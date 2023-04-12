package monnify

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) MonnifyInitPayment() (external_models.MonnifyInitPaymentResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyInitPaymentResponse
		logger           = r.Logger
		idata            = r.RequestData
		token            = getBase64Token(false)
	)

	data, ok := idata.(external_models.MonnifyInitPaymentRequest)
	if !ok {
		logger.Error("monnify init payment", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + token,
	}

	logger.Info("monnify init payment", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("monnify init payment", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	return outBoundResponse.ResponseBody, nil
}

func (r *RequestObj) MonnifyVerifyTransactionByReference() (external_models.MonnifyVerifyByReferenceResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyVerifyByReferenceResponse
		logger           = r.Logger
		idata            = r.RequestData
		token            = getBase64Token(false)
	)

	data, ok := idata.(string)
	if !ok {
		logger.Error("monnify verify transaction by reference", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + token,
	}

	logger.Info("monnify verify transaction by reference", data)
	err := r.getNewSendRequestObject(nil, headers, data).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("monnify verify transaction by reference", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	return outBoundResponse.ResponseBody, nil
}

func (r *RequestObj) MonnifyReserveAccount() (external_models.MonnifyReserveAccountResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyReserveAccountResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.MonnifyReserveAccountRequest)
	if !ok {
		logger.Error("monnify reserve account", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	token, err := r.getMonnifyLoginObject(false).MonnifyLogin()
	if err != nil {
		logger.Error("monnify reserve account", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	logger.Info("monnify reserve account", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("monnify reserve account", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	return outBoundResponse.ResponseBody, nil
}

func (r *RequestObj) GetMonnifyReserveAccountTransactions() (external_models.GetMonnifyReserveAccountTransactionsResponseBody, error) {

	var (
		outBoundResponse external_models.GetMonnifyReserveAccountTransactionsResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(string)
	if !ok {
		logger.Error("get monnify reserve account transactions", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	token, err := r.getMonnifyLoginObject(false).MonnifyLogin()
	if err != nil {
		logger.Error("get monnify reserve account transactions", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	logger.Info("get monnify reserve account transactions", data)
	err = r.getNewSendRequestObject(nil, headers, fmt.Sprintf("?accountReference=%v&page=0&size=100", data)).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get monnify reserve account transactions", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	return outBoundResponse.ResponseBody, nil
}
