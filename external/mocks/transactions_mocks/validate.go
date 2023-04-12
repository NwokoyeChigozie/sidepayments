package transactions_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func ValidateOnTransactions(logger *utility.Logger, idata interface{}) (bool, error) {

	_, ok := idata.(external_models.ValidateOnDBReq)
	if !ok {
		logger.Error("validate on transaction", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	logger.Info("validate on transaction", true)

	return true, nil
}
