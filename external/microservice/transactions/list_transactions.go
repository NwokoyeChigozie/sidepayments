package transactions

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) ListTransactionsByID() (external_models.TransactionByID, error) {
	var (
		outBoundResponse external_models.ListTransactionsByIDResponse
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(string)
	if !ok {
		logger.Error("list transactions by id", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("list transactions by id", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("list transactions by id", data)
	err = r.getNewSendRequestObject(data, headers, "/"+data).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("list transactions by id", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("list transactions by id", outBoundResponse)

	return outBoundResponse.Data, nil
}
