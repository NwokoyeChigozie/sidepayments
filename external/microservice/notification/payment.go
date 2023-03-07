package notification

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) WalletfundedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.WalletFundedNotificationRequest)
	if !ok {
		logger.Info("wallet funded notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("wallet funded notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("wallet funded notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("wallet funded notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("wallet funded notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) WalletDebitNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.WalletDebitNotificationRequest)
	if !ok {
		logger.Info("wallet debit notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("wallet debit notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("wallet debit notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("wallet debit notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("wallet debit notification", outBoundResponse)

	return nil, nil
}
