package auth

import (
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
