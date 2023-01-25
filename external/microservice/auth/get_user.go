package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetUser() (external_models.User, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetUserRequestModel)
	if !ok {
		logger.Info("get user", idata, "request data format error")
		return outBoundResponse.Data.User, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("get user", outBoundResponse, err.Error())
		return outBoundResponse.Data.User, err
	}
	logger.Info("get user", outBoundResponse)

	return outBoundResponse.Data.User, nil
}

func (r *RequestObj) SetUserAuthorizationRequiredStatus() (bool, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.SetUserAuthorizationRequiredStatusResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.SetUserAuthorizationRequiredStatusModel)
	if !ok {
		logger.Info("set user authorization required status", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("set user authorization required status", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("set user authorization required status", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("set user authorization required status", outBoundResponse)

	return outBoundResponse.Data, nil
}
