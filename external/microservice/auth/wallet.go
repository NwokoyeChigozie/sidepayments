package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) CreateWalletBalance() (external_models.WalletBalance, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.WalletBalanceResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CreateWalletRequest)
	if !ok {
		logger.Error("create wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create wallet", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("create wallet", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("create wallet", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) GetWalletBalanceByAccountIDAndCurrency() (external_models.WalletBalance, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.WalletBalanceResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetWalletRequest)
	if !ok {
		logger.Error("get wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get wallet", data)
	err := r.getNewSendRequestObject(data, headers, fmt.Sprintf("/%v/%v", strconv.Itoa(int(data.AccountID)), strings.ToUpper(data.Currency))).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get wallet", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get wallet", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) UpdateWalletBalance() (external_models.WalletBalance, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.WalletBalanceResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.UpdateWalletRequest)
	if !ok {
		logger.Error("update wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update wallet", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("update wallet", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("update wallet", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) CreateWalletHistory() (external_models.WalletHistory, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.WalletHistoryResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CreateWalletHistoryRequest)
	if !ok {
		logger.Error("create wallet history", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create wallet history", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("create wallet history", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("create wallet history", outBoundResponse)

	return outBoundResponse.Data, nil
}
func (r *RequestObj) CreateWalletTransaction() (external_models.WalletTransaction, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.WalletTransactionResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CreateWalletTransactionRequest)
	if !ok {
		logger.Error("create wallet transaction", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create wallet transaction", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("create wallet transaction", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("create wallet transaction", outBoundResponse)

	return outBoundResponse.Data, nil
}
