package appruve

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vesicash/payment-ms/external"
	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) AppruveVerifyID() (int, error) {

	var (
		token            = config.GetConfig().Appruve.AccessToken
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	fdata, ok := idata.(external_models.AppruveReqModelFirst)
	if !ok {
		logger.Error("appruve_verify_id", idata, "request data format error")
		return http.StatusInternalServerError, fmt.Errorf("request data format error")
	}

	data := external_models.AppruveReqModelMain{
		ID:           fdata.ID,
		FirstName:    fdata.FirstName,
		LastName:     fdata.LastName,
		MiddleName:   fdata.MiddleName,
		Gender:       fdata.Gender,
		Phone_number: fdata.Phone_number,
		DateOfBirth:  fdata.DateOfBirth,
	}
	logger.Info("appruve_verify_id", data)
	endpoint := "/" + strings.ToLower(fdata.CountryCode) + "/" + fdata.Endpoint

	err := r.getNewSendRequestObject(data, headers, endpoint).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("appruve_verify_id", outBoundResponse, err.Error())
		code := http.StatusInternalServerError
		if external.ResponseCode != 0 {
			code = external.ResponseCode
		}
		return code, err
	}
	logger.Info("appruve_verify_id", outBoundResponse)

	return external.ResponseCode, nil
}
