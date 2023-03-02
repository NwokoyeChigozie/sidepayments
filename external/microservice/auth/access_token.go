package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetAccessToken() (external_models.AccessToken, error) {
	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAccessTokenModel
		logger           = r.Logger
	)

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	err := r.getNewSendRequestObject(nil, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("get access_token", outBoundResponse, err)
		return outBoundResponse.Data, err
	}
	logger.Info("get access_token", outBoundResponse)

	return outBoundResponse.Data, nil
}
func (r *RequestObj) GetAccessTokenByKey() (external_models.AccessToken, error) {
	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAccessTokenModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(string)
	if !ok {
		logger.Info("get access token by key", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	err := r.getNewSendRequestObject(nil, headers, fmt.Sprintf("/%v", data)).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("get access token by key", outBoundResponse, err)
		return outBoundResponse.Data, err
	}
	logger.Info("get access token by key", outBoundResponse)

	return outBoundResponse.Data, nil
}
