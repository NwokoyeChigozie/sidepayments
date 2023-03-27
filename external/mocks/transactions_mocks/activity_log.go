package transactions_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func CreateActivityLog(logger *utility.Logger, idata interface{}) (interface{}, error) {

	var (
		outBoundResponse external_models.CreateActivityLogRequest
	)

	data, ok := idata.(external_models.UpdateTransactionAmountPaidRequest)
	if !ok {
		logger.Info("create activity log", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("create activity log", outBoundResponse, data)

	return nil, nil
}
