package auth_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	BusinessCharge *external_models.BusinessCharge
)

func GetBusinessCharge(logger *utility.Logger, idata interface{}) (external_models.BusinessCharge, error) {

	var (
		outBoundResponse external_models.GetBusinessChargeResponse
	)

	_, ok := idata.(external_models.GetBusinessChargeModel)
	if !ok {
		logger.Error("get business charge", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	if BusinessCharge == nil {
		logger.Error("get business charge", BusinessCharge, "businessCharge not provided")
		return external_models.BusinessCharge{}, fmt.Errorf("businessCharge not provided")
	}

	logger.Info("get business charge", BusinessCharge, "businessCharge found")

	return *BusinessCharge, nil
}

func InitBusinessCharge(logger *utility.Logger, idata interface{}) (external_models.BusinessCharge, error) {

	var (
		outBoundResponse external_models.GetBusinessChargeResponse
	)

	_, ok := idata.(external_models.InitBusinessChargeModel)
	if !ok {
		logger.Error("init business charge", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	if BusinessCharge == nil {
		logger.Error("init business charge", BusinessCharge, "businessCharge not provided")
		return external_models.BusinessCharge{}, fmt.Errorf("businessCharge not provided")
	}

	logger.Info("init business charge", BusinessCharge, "businessCharge found")

	return *BusinessCharge, nil
}
