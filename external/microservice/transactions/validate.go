package transactions

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) ValidateOnTransactions() (bool, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.ValidateOnDBReqModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.ValidateOnDBReq)
	if !ok {
		logger.Error("validate on transactions", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("validate on transactions", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("validate on transactions", outBoundResponse, err.Error())
		return false, err
	}
	logger.Info("validate on transactions", outBoundResponse)

	return outBoundResponse.Data, nil
}
