package auth_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	UserProfile *external_models.UserProfile
)

func GetUserProfile(logger *utility.Logger, idata interface{}) (external_models.UserProfile, error) {

	var (
		outBoundResponse external_models.GetUserProfileResponse
	)

	_, ok := idata.(external_models.GetUserProfileModel)
	if !ok {
		logger.Error("get user profile", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	if UserProfile == nil {
		logger.Error("get UserProfile", UserProfile, "UserProfile not provided")
		return external_models.UserProfile{}, fmt.Errorf("UserProfile not provided")
	}

	logger.Info("get UserProfile", UserProfile, "UserProfile found")
	return *UserProfile, nil
}
