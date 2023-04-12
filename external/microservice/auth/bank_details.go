package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetBankDetails() (external_models.BankDetail, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetBankDetailResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetBankDetailModel)
	if !ok {
		logger.Error("get bank detail", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get bank detail", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get bank detail", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get bank detail", outBoundResponse)

	return outBoundResponse.Data, nil
}
