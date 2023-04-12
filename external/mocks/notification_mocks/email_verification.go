package notification_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/mocks/auth_mocks"
	"github.com/vesicash/payment-ms/utility"
)

func SendVerificationEmail(logger *utility.Logger, idata interface{}) (interface{}, error) {

	_, ok := idata.(external_models.EmailNotificationRequest)
	if !ok {
		logger.Error("get user", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Error("verification email", nil, err)
		return nil, err
	}

	return nil, nil
}

func SendWelcomeEmail(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.AccountIDRequestModel)
	if !ok {
		logger.Error("welcome email", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Error("welcome email", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("welcome email", data)

	return nil, nil
}

func SendEmailVerifiedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.AccountIDRequestModel)
	if !ok {
		return nil, fmt.Errorf("request data format error")
	}

	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Error("email verified notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("email verified notification", data)

	return nil, nil
}
