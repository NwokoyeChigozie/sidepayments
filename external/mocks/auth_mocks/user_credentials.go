package auth_mocks

import (
	"fmt"
	"net/http"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/utility"
)

var (
	UsersCredential *external_models.UsersCredential
)

func GetUserCredential(logger *utility.Logger, idata interface{}) (external_models.GetUserCredentialResponse, error) {

	_, ok := idata.(external_models.GetUserCredentialModel)
	if !ok {
		logger.Error("get user credential", idata, "request data format error")
		return external_models.GetUserCredentialResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "request data format error",
			Data:    external_models.UsersCredential{},
		}, fmt.Errorf("request data format error")
	}

	if UsersCredential == nil {
		logger.Error("get user credential", UsersCredential, "user credential not provided")
		return external_models.GetUserCredentialResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "user not provided",
			Data:    external_models.UsersCredential{},
		}, fmt.Errorf("user not provided")
	}

	logger.Info("get user credential", UsersCredential, "user credential found")
	return external_models.GetUserCredentialResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "success",
		Data:    *UsersCredential,
	}, nil
}

func CreateUserCredential(logger *utility.Logger, idata interface{}) (external_models.GetUserCredentialResponse, error) {

	_, ok := idata.(external_models.CreateUserCredentialModel)
	if !ok {
		logger.Error("create user credential", idata, "request data format error")
		return external_models.GetUserCredentialResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "request data format error",
			Data:    external_models.UsersCredential{},
		}, fmt.Errorf("request data format error")
	}

	if UsersCredential == nil {
		logger.Error("create user credential", UsersCredential, "user credential not provided")
		return external_models.GetUserCredentialResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "user not provided",
			Data:    external_models.UsersCredential{},
		}, fmt.Errorf("user not provided")
	}

	logger.Info("create user credential", UsersCredential, "user credential found")
	return external_models.GetUserCredentialResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "success",
		Data:    *UsersCredential,
	}, nil

}

func UpdateUserCredential(logger *utility.Logger, idata interface{}) (external_models.GetUserCredentialResponse, error) {

	_, ok := idata.(external_models.UpdateUserCredentialModel)
	if !ok {
		logger.Error("update user credential", idata, "request data format error")
		return external_models.GetUserCredentialResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "request data format error",
			Data:    external_models.UsersCredential{},
		}, fmt.Errorf("request data format error")
	}

	if UsersCredential == nil {
		logger.Error("update user credential", UsersCredential, "user credential not provided")
		return external_models.GetUserCredentialResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "user not provided",
			Data:    external_models.UsersCredential{},
		}, fmt.Errorf("user not provided")
	}

	logger.Info("update user credential", UsersCredential, "user credential found")
	return external_models.GetUserCredentialResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "success",
		Data:    *UsersCredential,
	}, nil
}
