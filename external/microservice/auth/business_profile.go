package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetBusinessProfile() (external_models.BusinessProfile, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetBusinessProfileResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetBusinessProfileModel)
	if !ok {
		logger.Error("get business profile", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get business profile", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get business profile", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get business profile", outBoundResponse)

	return outBoundResponse.Data, nil
}
