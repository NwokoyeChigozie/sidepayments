package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetAuthorize() (external_models.GetAuthorizeResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAuthorizeResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetAuthorizeModel)
	if !ok {
		logger.Info("get authorize", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user credential", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("get authorize", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("get authorize", outBoundResponse)

	return outBoundResponse, nil
}

func (r *RequestObj) CreateAuthorize() (external_models.GetAuthorizeResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAuthorizeResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CreateAuthorizeModel)
	if !ok {
		logger.Info("create authorize", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create authorize", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("create authorize", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("create authorize", outBoundResponse)

	return outBoundResponse, nil
}

func (r *RequestObj) UpdateAuthorize() (external_models.GetAuthorizeResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAuthorizeResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.UpdateAuthorizeModel)
	if !ok {
		logger.Info("update authorize", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update authorize", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("update authorize", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("update authorize", outBoundResponse)

	return outBoundResponse, nil
}
