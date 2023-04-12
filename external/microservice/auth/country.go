package auth

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) GetCountry() (external_models.Country, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetCountryResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetCountryModel)
	if !ok {
		logger.Error("get country", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get country", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get country", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get country", outBoundResponse)

	return outBoundResponse.Data, nil
}
