package payment

import (
	"fmt"
	"strconv"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/utility"
)

func ListTransactionsByID(extReq request.ExternalRequest, transactionID string) (external_models.TransactionByID, error) {

	transactionInterface, err := extReq.SendExternalRequest(request.ListTransactionsByID, transactionID)

	if err != nil {
		extReq.Logger.Info(err.Error())
		return external_models.TransactionByID{}, fmt.Errorf("transaction could not be retrieved")
	}
	transaction, ok := transactionInterface.(external_models.TransactionByID)
	if !ok {
		return external_models.TransactionByID{}, fmt.Errorf("response data format error")
	}
	if transaction.ID == 0 {
		return external_models.TransactionByID{}, fmt.Errorf("transaction not found")
	}

	return transaction, nil
}

func GetUserWithAccountID(extReq request.ExternalRequest, accountID int) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{AccountID: uint(accountID)})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}

func GetUsersByBusinessID(extReq request.ExternalRequest, BusinessId int) ([]external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUsersByBusinessID, strconv.Itoa(BusinessId))
	if err != nil {
		return []external_models.User{}, err
	}

	us, ok := usItf.([]external_models.User)
	if !ok {
		return []external_models.User{}, fmt.Errorf("response data format error")
	}

	return us, nil
}

func GetUserProfileByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.UserProfile, error) {
	userProfileInterface, err := extReq.SendExternalRequest(request.GetUserProfile, external_models.GetUserProfileModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Info(err.Error())
		return external_models.UserProfile{}, err
	}

	userProfile, ok := userProfileInterface.(external_models.UserProfile)
	if !ok {
		return external_models.UserProfile{}, fmt.Errorf("response data format error")
	}

	if userProfile.ID == 0 {
		return external_models.UserProfile{}, fmt.Errorf("user profile not found")
	}

	return userProfile, nil

}
func GetCountryByNameOrCode(extReq request.ExternalRequest, logger *utility.Logger, NameOrCode string) (external_models.Country, error) {

	countryInterface, err := extReq.SendExternalRequest(request.GetCountry, external_models.GetCountryModel{
		Name: NameOrCode,
	})

	if err != nil {
		logger.Info(err.Error())
		return external_models.Country{}, fmt.Errorf("Your country could not be resolved, please update your profile.")
	}
	country, ok := countryInterface.(external_models.Country)
	if !ok {
		return external_models.Country{}, fmt.Errorf("response data format error")
	}
	if country.ID == 0 {
		return external_models.Country{}, fmt.Errorf("Your country could not be resolved, please update your profile")
	}

	return country, nil
}
