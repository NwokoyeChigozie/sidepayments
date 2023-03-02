package rave

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) RaveInitPayment() (external_models.RaveInitPaymentResponse, error) {

	var (
		outBoundResponse external_models.RaveInitPaymentResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(external_models.RaveInitPaymentRequest)
	if !ok {
		logger.Info("init payment rave", idata, "request data format error")
		return external_models.RaveInitPaymentResponse{}, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("init payment rave", outBoundResponse, err.Error())
		return external_models.RaveInitPaymentResponse{}, err
	}
	logger.Info("init payment rave", outBoundResponse)

	return outBoundResponse, nil
}
