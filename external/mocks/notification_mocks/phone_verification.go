package notification_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/mocks/auth_mocks"
	"github.com/vesicash/payment-ms/utility"
)

func SendSendSMSToPhone(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.SMSToPhoneNotificationRequest)
	if !ok {
		logger.Info("send sms to phone", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("send sms to phone", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("send sms to phone", data)

	return nil, nil
}
