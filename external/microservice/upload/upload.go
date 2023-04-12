package upload

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/internal/config"
)

func (r *RequestObj) UploadFile() (external_models.UploadFileResponseData, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.UploadFileResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.UploadFileRequest)
	if !ok {
		logger.Error("upload one file", idata, "request data format error")
		return external_models.UploadFileResponseData{}, fmt.Errorf("request data format error")
	}

	// requestBody := new(bytes.Buffer)
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("files", data.PlaceHolderName)
	if err != nil {
		logger.Error("upload one file", outBoundResponse, err.Error())
		return external_models.UploadFileResponseData{}, err
	}
	if _, err := part.Write(data.File); err != nil {
		logger.Error("upload one file", outBoundResponse, err.Error())
		return external_models.UploadFileResponseData{}, err
	}

	err = writer.Close()
	if err != nil {
		logger.Error("upload one file", outBoundResponse, err.Error())
		return external_models.UploadFileResponseData{}, err
	}

	r.RequestBody = &requestBody
	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"v-app":        appKey,
	}

	logger.Info("upload one file", string(data.File), data.PlaceHolderName)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("upload one file", outBoundResponse, err.Error())
		return external_models.UploadFileResponseData{}, err
	}
	logger.Info("upload one file", outBoundResponse)

	if len(outBoundResponse.Data) < 1 {
		err = fmt.Errorf("no link returned")
		logger.Error("upload one file", outBoundResponse, err)
		return external_models.UploadFileResponseData{}, err
	}

	return outBoundResponse.Data[0], nil
}
