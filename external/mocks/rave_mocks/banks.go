package rave_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func ListBanksWithRave(logger *utility.Logger, idata interface{}) ([]external_models.BanksResponse, error) {

	var (
		outBoundResponse external_models.ListBanksResponse
	)

	_, ok := idata.(string)
	if !ok {
		logger.Info("list banks with rave", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	return []external_models.BanksResponse{
		{
			ID:   1,
			Code: "221",
			Name: "vesicash bank",
		},
	}, nil
}

func ConvertCurrencyWithRave(logger *utility.Logger, idata interface{}) (external_models.ConvertCurrencyData, error) {

	var (
		outBoundResponse external_models.ConvertCurrencyResponse
	)

	data, ok := idata.(external_models.ConvertCurrencyRequest)
	if !ok {
		logger.Info("convert currency with rave", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("convert currency with rave", outBoundResponse)

	return external_models.ConvertCurrencyData{
		Rate: 0.8,
		Source: external_models.ConvertCurrencyDataSourceOrDestination{
			Amount:   data.Amount * 0.8,
			Currency: data.To,
		},
		Destination: external_models.ConvertCurrencyDataSourceOrDestination{
			Amount:   data.Amount,
			Currency: data.From,
		},
	}, nil
}
