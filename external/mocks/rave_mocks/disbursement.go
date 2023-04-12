package rave_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func RaveInitTransfer(logger *utility.Logger, idata interface{}) (external_models.RaveInitTransferResponse, error) {

	var (
		outBoundResponse external_models.RaveInitTransferResponse
	)

	data, ok := idata.(external_models.RaveInitTransferRequest)
	if !ok {
		logger.Error("rave init transfer", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	logger.Info("rave init transfer", outBoundResponse)

	return external_models.RaveInitTransferResponse{
		Status:  "success",
		Message: "Transfer Queued Successfully",
		Data: external_models.RaveInitTransferResponseData{
			ID:            262251,
			AccountNumber: "0690000040",
			BankCode:      "221",
			FullName:      data.BeneficiaryName,
			Currency:      data.Currency,
			DebitCurrency: data.DebitCurrency,
			Amount:        data.Amount,
			Fee:           0,
			Status:        "NEW",
			Reference:     data.Reference,
			IsApproved:    1,
			BankName:      "vesicash bank",
		},
	}, nil
}
