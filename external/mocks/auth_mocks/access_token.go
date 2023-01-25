package auth_mocks

import (
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	AccessToken external_models.AccessToken
)

func GetAccessToken(logger *utility.Logger) (external_models.AccessToken, error) {
	logger.Info("get access tokens", "get access tokens called")
	return AccessToken, nil
}
