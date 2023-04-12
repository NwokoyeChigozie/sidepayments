package rave

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) RaveInitTransfer() (external_models.RaveInitTransferResponse, error) {

	var (
		outBoundResponse external_models.RaveInitTransferResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(external_models.RaveInitTransferRequest)
	if !ok {
		logger.Error("rave init transfer", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("rave init transfer", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("rave init transfer", outBoundResponse)

	return outBoundResponse, nil
}
