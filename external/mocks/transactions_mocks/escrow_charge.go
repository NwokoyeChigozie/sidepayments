package transactions_mocks

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	EscrowCharge external_models.GetEscrowChargeResponseData
)

func GetEscrowCharge(logger *utility.Logger, idata interface{}) (external_models.GetEscrowChargeResponseData, error) {
	var (
		outBoundResponse external_models.GetEscrowChargeResponse
	)
	data, ok := idata.(external_models.GetEscrowChargeRequest)
	if !ok {
		logger.Info("get escrow charge", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("get escrow charge", outBoundResponse)

	return external_models.GetEscrowChargeResponseData{
		Amount: data.Amount - 7,
		Charge: 7,
	}, nil
}
