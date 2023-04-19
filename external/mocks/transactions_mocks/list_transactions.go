package transactions_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	ListTransactionsByIDObj *external_models.TransactionByID
	ListTransactionsObj     []external_models.TransactionByID
)

func ListTransactionsByID(logger *utility.Logger, idata interface{}) (external_models.TransactionByID, error) {
	var (
		outBoundResponse external_models.ListTransactionsByIDResponse
	)
	_, ok := idata.(string)
	if !ok {
		logger.Error("list transactions by id", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	if ListTransactionsByIDObj == nil {
		logger.Error("list transactions by id", ListTransactionsByIDObj, "ListTransactionsByIDObj not provided")
		return external_models.TransactionByID{}, fmt.Errorf("ListTransactionsByIDObj not provided")
	}

	logger.Info("list transactions by id", outBoundResponse)

	return *ListTransactionsByIDObj, nil
}

func ListTransactions(logger *utility.Logger, idata interface{}) ([]external_models.TransactionByID, error) {
	var (
		outBoundResponse external_models.ListTransactionsResponse
		queryParam       = ""
	)

	dataMid, ok := idata.(external_models.ListTransactionsRequestMid)
	if !ok {
		logger.Error("list transactions", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("list transactions", dataMid)
	data := external_models.ListTransactionsRequest{
		Status:     dataMid.Status,
		StatusCode: dataMid.StatusCode,
		Filter:     dataMid.Filter,
	}

	if dataMid.Limit != 0 {
		queryParam += fmt.Sprintf("?limit=%v", dataMid.Limit)
		if dataMid.Page != 0 {
			queryParam += fmt.Sprintf("&page=%v", dataMid.Page)
		}
	}

	logger.Info("list transactions", data)

	return ListTransactionsObj, nil
}
