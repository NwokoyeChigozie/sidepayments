package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) ValidateOnAuth() (bool, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.ValidateOnDBReqModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.ValidateOnDBReq)
	if !ok {
		logger.Info("validate on auth", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("validate on auth", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("validate on auth", outBoundResponse, err.Error())
		return false, err
	}
	logger.Info("validate on auth", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) ValidateAuthorization() (external_models.ValidateAuthorizationDataModel, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.ValidateAuthorizationModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.ValidateAuthorizationReq)
	if !ok {
		logger.Info("validate authorization", idata, "request data format error")
		return external_models.ValidateAuthorizationDataModel{}, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("validate authorization", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("validate authorization", outBoundResponse, err.Error())
		return external_models.ValidateAuthorizationDataModel{}, err
	}
	logger.Info("validate authorization", outBoundResponse.Data)

	return outBoundResponse.Data, nil
}
