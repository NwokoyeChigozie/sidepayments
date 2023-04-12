package rave_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	AccountName string
)

func RaveResolveBankAccount(logger *utility.Logger, idata interface{}) (string, error) {

	_, ok := idata.(external_models.ResolveAccountRequest)
	if !ok {
		logger.Error("rave resolve bank account", idata, "request data format error")
		return "", fmt.Errorf("request data format error")
	}

	if AccountName == "" {
		logger.Error("rave resolve bank account", "account name not provided", AccountName)
		return "", fmt.Errorf("account name not provided")
	}

	logger.Info("rave resolve bank account", AccountName)

	return AccountName, nil
}
