package verification_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func CheckVerification(logger *utility.Logger, idata interface{}) (external_models.CheckVerificationResponseData, error) {

	var (
		outBoundResponse external_models.CheckVerificationResponse
	)

	data, ok := idata.(external_models.CheckVerificationRequest)
	if !ok {
		logger.Error("check verification", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("check verification", outBoundResponse, data)

	return external_models.CheckVerificationResponseData{
		Verified: true,
	}, nil
}
