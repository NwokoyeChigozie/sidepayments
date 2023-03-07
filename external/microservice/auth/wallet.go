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
		logger.Info("create wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create wallet", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("create wallet", outBoundResponse, err.Error())
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
		logger.Info("get wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get wallet", data)
	err := r.getNewSendRequestObject(data, headers, fmt.Sprintf("/%v/%v", strconv.Itoa(int(data.AccountID)), strings.ToUpper(data.Currency))).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("get wallet", outBoundResponse, err.Error())
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
		logger.Info("update wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update wallet", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("update wallet", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("update wallet", outBoundResponse)

	return outBoundResponse.Data, nil
}
