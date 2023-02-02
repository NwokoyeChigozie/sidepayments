package monnify_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/utility"
)

func MonnifyLogin(logger *utility.Logger, idata interface{}) (string, error) {

	var (
		base64Key        = config.GetConfig().Monnify.MonnifyBase64Key
		outBoundResponse external_models.MonnifyLoginResponse
	)

	if base64Key == "" {
		logger.Info("monnify login", "monnify base64 key not in env", base64Key)
		return "", fmt.Errorf("monnify base64 key not in env: %v", base64Key)
	}
	logger.Info("monnify login", outBoundResponse)

	return outBoundResponse.ResponseBody.AccessToken, nil
}

func MonnifyMatchBvnDetails(logger *utility.Logger, idata interface{}) (bool, error) {

	var (
		outBoundResponse external_models.MonnifyMatchBvnDetailsResponse
	)

	_, err := MonnifyLogin(logger, nil)
	if err != nil {
		logger.Info("monnify match bvn details", outBoundResponse, err.Error())
		return false, err
	}

	data, ok := idata.(external_models.MonnifyMatchBvnDetailsReq)
	if !ok {
		logger.Info("monnify match bvn details", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	logger.Info("monnify match bvn details", data)
	return true, nil
}
