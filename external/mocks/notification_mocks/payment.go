package notification_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func WalletfundedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.WalletFundedNotificationRequest)
	if !ok {
		logger.Error("wallet funded notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("wallet funded notification", outBoundResponse, data)

	return nil, nil
}

func WalletDebitNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.WalletDebitNotificationRequest)
	if !ok {
		logger.Error("wallet debit notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("wallet debit notification", outBoundResponse, data)

	return nil, nil
}

func PaymentInvoiceNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.PaymentInvoiceNotificationRequest)
	if !ok {
		logger.Error("payment invoice notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("payment invoice notification", outBoundResponse, data)

	return nil, nil
}

func TransactionPaidNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.OnlyTransactionIDAndAccountIDRequest)
	if !ok {
		logger.Error("transaction paid notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("transaction paid notification", outBoundResponse, data)

	return nil, nil
}
