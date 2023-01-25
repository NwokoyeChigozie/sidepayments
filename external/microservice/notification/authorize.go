package notification

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) SendAuthorizedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.AuthorizeNotificationRequest)
	if !ok {
		logger.Info("authorized notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("authorized notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("authorized notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("authorized notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("authorized notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) SendAuthorizationNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.AuthorizeNotificationRequest)
	if !ok {
		logger.Info("authorization notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("authorization notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("verification failed notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("authorization notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("authorization notification", outBoundResponse)

	return nil, nil
}
