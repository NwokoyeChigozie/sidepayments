package auth_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	ValidateAuthorizationRes *external_models.ValidateAuthorizationDataModel
)

func ValidateOnAuth(logger *utility.Logger, idata interface{}) (bool, error) {

	_, ok := idata.(external_models.ValidateOnDBReq)
	if !ok {
		logger.Error("validate on auth", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	logger.Info("validate on auth", true)

	return true, nil
}

func ValidateAuthorization(logger *utility.Logger, idata interface{}) (external_models.ValidateAuthorizationDataModel, error) {

	_, ok := idata.(external_models.ValidateAuthorizationReq)
	if !ok {
		logger.Error("validate authorization", idata, "request data format error")
		return external_models.ValidateAuthorizationDataModel{}, fmt.Errorf("request data format error")
	}

	if ValidateAuthorizationRes == nil {
		logger.Error("validate authorization", User, "validate authorization response not provided")
		return external_models.ValidateAuthorizationDataModel{}, fmt.Errorf("validate authorization response not provided")
	}

	logger.Info("validate authorization", ValidateAuthorizationRes)
	return *ValidateAuthorizationRes, nil
}
