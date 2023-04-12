package monnify_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func MonnifyInitTransfer(logger *utility.Logger, idata interface{}) (external_models.MonnifyInitTransferResponse, error) {

	var (
		outBoundResponse external_models.MonnifyInitTransferResponse
	)

	data, ok := idata.(external_models.MonnifyInitTransferRequest)
	if !ok {
		logger.Error("monnify init transfer", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	return external_models.MonnifyInitTransferResponse{
		RequestSuccessful: true,
		ResponseMessage:   "success",
		ResponseCode:      "0",
		ResponseBody: external_models.MonnifyInitTransferResponseBody{
			DestinationAccountNumber: data.DestinationAccountNumber,
			DestinationAccountName:   data.DestinationAccountName,
			Amount:                   data.Amount,
			TotalFee:                 0,
			DestinationBankCode:      data.DestinationBankCode,
			DestinationBankName:      "vesicash bank",
			DateCreated:              "2022-11-03T14:11:12.659+0000",
			Reference:                data.Reference,
			SessionId:                data.Reference,
		},
	}, nil
}
