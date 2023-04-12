package transactions

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) GetEscrowCharge() (external_models.GetEscrowChargeResponseData, error) {
	var (
		outBoundResponse external_models.GetEscrowChargeResponse
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.GetEscrowChargeRequest)
	if !ok {
		logger.Error("get escrow charge", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("get escrow charge", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("get escrow charge", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get escrow charge", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get escrow charge", outBoundResponse)

	return outBoundResponse.Data, nil
}
