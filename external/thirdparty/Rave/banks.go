package rave

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) ListBanksWithRave() ([]external_models.BanksResponse, error) {

	var (
		outBoundResponse external_models.ListBanksResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(string)
	if !ok {
		logger.Info("list banks with rave", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, headers, "/"+data).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("list banks with rave", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("list banks with rave", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) ConvertCurrencyWithRave() (external_models.ConvertCurrencyData, error) {

	var (
		outBoundResponse external_models.ConvertCurrencyResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(external_models.ConvertCurrencyRequest)
	if !ok {
		logger.Info("convert currency with rave", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, headers, fmt.Sprintf("?amount=%v&destination_currency=%v&source_currency=%v", data.Amount, data.From, data.To)).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("convert currency with rave", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("convert currency with rave", outBoundResponse)

	return outBoundResponse.Data, nil
}
