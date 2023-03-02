package monnify

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) MonnifyInitPayment() (external_models.MonnifyInitPaymentResponseBody, error) {

	var (
		outBoundResponse external_models.MonnifyInitPaymentResponse
		logger           = r.Logger
		idata            = r.RequestData
		token            = getBase64Token()
	)

	data, ok := idata.(external_models.MonnifyInitPaymentRequest)
	if !ok {
		logger.Info("monnify init payment", idata, "request data format error")
		return outBoundResponse.ResponseBody, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + token,
	}

	logger.Info("monnify match bvn details", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("monnify init payment", outBoundResponse, err.Error())
		return outBoundResponse.ResponseBody, err
	}

	return outBoundResponse.ResponseBody, nil
}
