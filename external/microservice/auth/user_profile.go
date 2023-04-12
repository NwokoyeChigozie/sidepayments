package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetUserProfile() (external_models.UserProfile, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserProfileResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetUserProfileModel)
	if !ok {
		logger.Error("get user profile", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user profile", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get user profile", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get user profile", outBoundResponse)

	return outBoundResponse.Data, nil
}
