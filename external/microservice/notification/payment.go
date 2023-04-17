package notification

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) WalletfundedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.WalletFundedNotificationRequest)
	if !ok {
		logger.Error("wallet funded notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("wallet funded notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("wallet funded notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("wallet funded notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("wallet funded notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) WalletDebitNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.WalletDebitNotificationRequest)
	if !ok {
		logger.Error("wallet debit notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("wallet debit notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("wallet debit notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("wallet debit notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("wallet debit notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) PaymentInvoiceNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.PaymentInvoiceNotificationRequest)
	if !ok {
		logger.Error("payment invoice notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("payment invoice notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("payment invoice notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("payment invoice notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("payment invoice notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) TransactionPaidNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.OnlyTransactionIDAndAccountIDRequest)
	if !ok {
		logger.Error("transaction paid notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("transaction paid notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction paid notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("transaction paid notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction paid notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) SuccessfulRefundNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.OnlyTransactionIDAndAccountIDRequest)
	if !ok {
		logger.Error("successful refund notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("successful refund notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("successful refund notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("successful refund notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("successful refund notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) EscrowDisbursedSellerNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("escrow disbursed seller notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("escrow disbursed seller notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("escrow disbursed seller notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("escrow disbursed seller notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("escrow disbursed seller notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) EscrowDisbursedBuyerNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("escrow disbursed buyer notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("escrow disbursed buyer notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("escrow disbursed buyer notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("escrow disbursed buyer notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("escrow disbursed buyer notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) TransactionClosedBuyerNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("transaction closed buyer notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("transaction closed buyer notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction closed buyer notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("transaction closed buyer notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction closed buyer notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) TransactionClosedSellerNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.OnlyTransactionIDRequiredRequest)
	if !ok {
		logger.Error("transaction closed seller notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Error("transaction closed seller notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction closed seller notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("transaction closed seller notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction closed seller notification", outBoundResponse)

	return nil, nil
}
