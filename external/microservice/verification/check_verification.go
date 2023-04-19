package verification

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) CheckVerification() (external_models.CheckVerificationResponseData, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.CheckVerificationResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CheckVerificationRequest)
	if !ok {
		logger.Error("check verification", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("check verification", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("check verification", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("check verification", outBoundResponse)

	return outBoundResponse.Data, nil
}
