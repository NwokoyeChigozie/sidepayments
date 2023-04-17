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

func SuccessfulRefundNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.OnlyTransactionIDAndAccountIDRequest)
	if !ok {
		logger.Error("successful refund notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	logger.Info("successful refund notification", outBoundResponse, data)

	return nil, nil
}
func EscrowDisbursedSellerNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("escrow disbursed seller notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	logger.Info("escrow disbursed seller notification", outBoundResponse, data)

	return nil, nil
}
func EscrowDisbursedBuyerNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("escrow disbursed buyer notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	logger.Info("escrow disbursed buyer notification", outBoundResponse, data)

	return nil, nil
}
func TransactionClosedBuyerNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("transaction closed buyer notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	logger.Info("transaction closed buyer notification", outBoundResponse, data)

	return nil, nil
}
func TransactionClosedSellerNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("transaction closed seller notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	logger.Info("transaction closed seller notification", outBoundResponse, data)

	return nil, nil
}
