package monnify

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) MonnifyInitTransfer() (external_models.MonnifyInitTransferResponse, error) {

	var (
		outBoundResponse external_models.MonnifyInitTransferResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.MonnifyInitTransferRequest)
	if !ok {
		logger.Info("monnify init transfer", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	token, err := r.getMonnifyLoginObject().MonnifyLogin()
	if err != nil {
		logger.Info("monnify init transfer", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	logger.Info("monnify init transfer", data)
	err = r.getNewSendRequestObject(nil, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("monnify init transfer", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
