package transactions

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) UpdateTransactionAmountPaid() (external_models.Transaction, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.UpdateTransactionAmountPaidResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.UpdateTransactionAmountPaidRequest)
	if !ok {
		logger.Info("update transaction amount paid", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update transaction amount paid", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("update wallet", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("update transaction amount paid", outBoundResponse)

	return outBoundResponse.Data, nil
}
