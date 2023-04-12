package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetBusinessCharge() (external_models.BusinessCharge, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetBusinessChargeResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetBusinessChargeModel)
	if !ok {
		logger.Error("get business charge", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get business charge", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get business charge", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get business charge", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) InitBusinessCharge() (external_models.BusinessCharge, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetBusinessChargeResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.InitBusinessChargeModel)
	if !ok {
		logger.Error("init business charge", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("init business charge", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("init business charge", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("init business charge", outBoundResponse)

	return outBoundResponse.Data, nil
}
