package transactions

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) CreateActivityLog() (interface{}, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.CreateActivityLogRequest
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.UpdateTransactionAmountPaidRequest)
	if !ok {
		logger.Info("create activity log", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create activity log", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("update wallet", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("create activity log", outBoundResponse)

	return nil, nil
}
