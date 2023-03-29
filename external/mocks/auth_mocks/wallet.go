package auth_mocks

import (
	"fmt"
	"strconv"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func CreateWalletBalance(logger *utility.Logger, idata interface{}) (external_models.WalletBalance, error) {

	var (
		outBoundResponse external_models.WalletBalanceResponse
	)

	data, ok := idata.(external_models.CreateWalletRequest)
	if !ok {
		logger.Info("create wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("create wallet", outBoundResponse)

	return external_models.WalletBalance{
		AccountID: int(data.AccountID),
		Available: data.Available,
		Currency:  data.Currency,
	}, nil
}

func GetWalletBalanceByAccountIDAndCurrency(logger *utility.Logger, idata interface{}) (external_models.WalletBalance, error) {

	var (
		outBoundResponse external_models.WalletBalanceResponse
	)

	data, ok := idata.(external_models.GetWalletRequest)
	if !ok {
		logger.Info("get wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("get wallet", outBoundResponse)

	return external_models.WalletBalance{
		AccountID: int(data.AccountID),
		Currency:  data.Currency,
		Available: 20000,
	}, nil
}

func UpdateWalletBalance(logger *utility.Logger, idata interface{}) (external_models.WalletBalance, error) {

	var (
		outBoundResponse external_models.WalletBalanceResponse
	)

	data, ok := idata.(external_models.UpdateWalletRequest)
	if !ok {
		logger.Info("update wallet", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("update wallet", outBoundResponse)

	return external_models.WalletBalance{
		ID:        data.ID,
		Available: data.Available,
	}, nil
}

func CreateWalletHistory(logger *utility.Logger, idata interface{}) (external_models.WalletHistory, error) {

	var (
		outBoundResponse external_models.WalletHistoryResponse
	)

	data, ok := idata.(external_models.CreateWalletHistoryRequest)
	if !ok {
		logger.Info("create wallet history", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("create wallet history", outBoundResponse)

	return external_models.WalletHistory{
		ID:               20,
		AccountID:        strconv.Itoa(data.AccountID),
		Reference:        data.Reference,
		Currency:         data.Currency,
		Type:             data.Type,
		AvailableBalance: data.AvailableBalance,
		Amount:           data.Amount,
	}, nil
}
func CreateWalletTransaction(logger *utility.Logger, idata interface{}) (external_models.WalletTransaction, error) {

	var (
		outBoundResponse external_models.WalletTransactionResponse
	)

	data, ok := idata.(external_models.CreateWalletTransactionRequest)
	if !ok {
		logger.Info("create wallet transaction", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("create wallet transaction", outBoundResponse)

	return external_models.WalletTransaction{
		SenderAccountID:   strconv.Itoa(data.SenderAccountID),
		ReceiverAccountID: strconv.Itoa(data.ReceiverAccountID),
		SenderAmount:      data.SenderAmount,
		ReceiverAmount:    data.ReceiverAmount,
		SenderCurrency:    data.SenderCurrency,
		ReceiverCurrency:  data.ReceiverCurrency,
		Approved:          data.Approved,
		FirstApproval:     data.FirstApproval,
	}, nil
}
