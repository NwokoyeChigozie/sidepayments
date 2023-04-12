package appruve_mocks

import (
	"fmt"
	"net/http"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

func AppruveVerifyID(logger *utility.Logger, idata interface{}) (int, error) {

	var (
		outBoundResponse map[string]interface{}
	)

	_, ok := idata.(external_models.AppruveReqModelFirst)
	if !ok {
		logger.Error("appruve_verify_id request data format error", idata)
		return http.StatusInternalServerError, fmt.Errorf("request data format error")
	}

	logger.Info("appruve_verify_id", outBoundResponse)

	return http.StatusOK, nil
}
