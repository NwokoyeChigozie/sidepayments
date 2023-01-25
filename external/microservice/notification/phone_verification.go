package notification

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) SendSendSMSToPhone() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.SMSToPhoneNotificationRequest)
	if !ok {
		logger.Info("sms to phone", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("send sms to phone", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("send sms to phone", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("send sms to phone", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("send sms to phone", outBoundResponse)

	return nil, nil
}
