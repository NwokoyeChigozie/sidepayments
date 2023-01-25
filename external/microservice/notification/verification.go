package notification

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) VerificationFailedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.VerificationFailedModel)
	if !ok {
		logger.Info("verification failed notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("verification failed notification", outBoundResponse, err.Error())
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
		logger.Info("verification failed notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("verification failed notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) VerificationSuccessfulNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.VerificationSuccessfulModel)
	if !ok {
		logger.Info("verification successful notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("verification successful notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("verification successful notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("verification successful notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("verification successful notification", outBoundResponse)

	return nil, nil
}
