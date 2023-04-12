package monnify

import (
	"fmt"

	"github.com/vesicash/payment-ms/external/external_models"
)

func (r *RequestObj) MonnifyLogin() (string, error) {

	var (
		base64Key        = getBase64Token(r.IsLiveMust)
		outBoundResponse external_models.MonnifyLoginResponse
		logger           = r.Logger
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + base64Key,
	}

	err := r.getNewSendRequestObject(nil, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("monnify login", outBoundResponse, err.Error())
		return "", err
	}
	logger.Info("monnify login", outBoundResponse)

	return outBoundResponse.ResponseBody.AccessToken, nil
}

func (r *RequestObj) MonnifyMatchBvnDetails() (bool, error) {

	var (
		outBoundResponse external_models.MonnifyMatchBvnDetailsResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	token, err := r.getMonnifyLoginObject(true).MonnifyLogin()
	if err != nil {
		logger.Error("monnify match bvn details", outBoundResponse, err.Error())
		return false, err
	}

	data, ok := idata.(external_models.MonnifyMatchBvnDetailsReq)
	if !ok {
		logger.Error("monnify match bvn details", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	logger.Info("monnify match bvn details", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("monnify match bvn details", outBoundResponse, err.Error())
		return false, err
	}
	logger.Info("monnify match bvn details", outBoundResponse)
	if !outBoundResponse.RequestSuccessful {
		logger.Error("monnify match bvn details", "request not successful: "+outBoundResponse.ResponseMessage)
		return false, fmt.Errorf("request not successful: " + outBoundResponse.ResponseMessage)
	}

	if outBoundResponse.ResponseBody.DateOfBirth != "FULL_MATCH" {
		logger.Error("monnify match bvn details", "bvn does not match date of birth", outBoundResponse.ResponseBody.DateOfBirth)
		return false, fmt.Errorf("bvn does not match date of birth: %v", outBoundResponse.ResponseBody.DateOfBirth)
	}

	return true, nil
}
