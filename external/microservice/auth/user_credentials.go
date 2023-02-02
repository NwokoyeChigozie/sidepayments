package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetUserCredential() (external_models.GetUserCredentialResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserCredentialResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetUserCredentialModel)
	if !ok {
		logger.Info("get user credential", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user credential", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("get user credential", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("get user credential", outBoundResponse)

	return outBoundResponse, nil
}

func (r *RequestObj) CreateUserCredential() (external_models.GetUserCredentialResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserCredentialResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CreateUserCredentialModel)
	if !ok {
		logger.Info("create user credential", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create user credential", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("create user credential", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("create user credential", outBoundResponse)

	return outBoundResponse, nil
}

func (r *RequestObj) UpdateUserCredential() (external_models.GetUserCredentialResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserCredentialResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.UpdateUserCredentialModel)
	if !ok {
		logger.Info("update user credential", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update user credential", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("update user credential", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("update user credential", outBoundResponse)

	return outBoundResponse, nil
}
