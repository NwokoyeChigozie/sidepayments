package transactions

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) CreateExchangeTransaction() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
		appKey           = config.GetConfig().App.Key
	)
	data, ok := idata.(external_models.CreateExchangeTransactionRequest)
	if !ok {
		logger.Error("create exchange transaction", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create exchange transaction", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("create exchange transaction", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("create exchange transaction", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) GetRateByID() (external_models.Rate, error) {
	var (
		outBoundResponse external_models.RateResponse
		logger           = r.Logger
		idata            = r.RequestData
		appKey           = config.GetConfig().App.Key
	)
	data, ok := idata.(int)
	if !ok {
		logger.Error("get rate by id", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get rate by id", data)
	err := r.getNewSendRequestObject(data, headers, fmt.Sprintf("/%v", data)).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get rate by id", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get rate by id", outBoundResponse)

	return outBoundResponse.Data, nil
}
