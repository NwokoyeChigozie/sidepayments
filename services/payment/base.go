package payment

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
func GetAccessTokenByKeyFromRequest(extReq request.ExternalRequest, c *gin.Context) (external_models.AccessToken, error) {
	privateKey := utility.GetHeader(c, "v-private-key")
	publicKey := utility.GetHeader(c, "v-public-key")
	key := privateKey
	if key == "" {
		key = publicKey
	}
	acItf, err := extReq.SendExternalRequest(request.GetAccessTokenByKey, key)
	if err != nil {
		return external_models.AccessToken{}, err
	}

	accessToken, ok := acItf.(external_models.AccessToken)
	if !ok {
		return external_models.AccessToken{}, fmt.Errorf("response data format error")
	}

	return accessToken, nil
}

func GetEscrowCharge(extReq request.ExternalRequest, businessId int, amount float64) (external_models.GetEscrowChargeResponseData, error) {
	escItf, err := extReq.SendExternalRequest(request.GetEscrowCharge, external_models.GetEscrowChargeRequest{
		BusinessID: businessId,
		Amount:     amount,
	})
	if err != nil {
		return external_models.GetEscrowChargeResponseData{}, err
	}

	escrowChargeData, ok := escItf.(external_models.GetEscrowChargeResponseData)
	if !ok {
		return external_models.GetEscrowChargeResponseData{}, fmt.Errorf("response data format error")
	}

	return escrowChargeData, nil
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

func isRequestIPNigerian(extReq request.ExternalRequest, c *gin.Context) (bool, error) {
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return false, err
	}

	ipInterface, err := extReq.SendExternalRequest(request.ResolveIP, ip)
	if err != nil {
		extReq.Logger.Info(err.Error())
		return false, err
	}

	ipResponse, ok := ipInterface.(external_models.ResolveIpResponse)
	if !ok {
		return false, fmt.Errorf("response data format error")
	}

	if strings.ToUpper(ipResponse.CountryCode) != "NG" {
		return false, nil
	}

	return true, nil
}

func GetBusinessProfileByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.BusinessProfile, error) {
	businessProfileInterface, err := extReq.SendExternalRequest(request.GetBusinessProfile, external_models.GetBusinessProfileModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Info(err.Error())
		return external_models.BusinessProfile{}, fmt.Errorf("Business lacks a profile.")
	}

	businessProfile, ok := businessProfileInterface.(external_models.BusinessProfile)
	if !ok {
		return external_models.BusinessProfile{}, fmt.Errorf("response data format error")
	}

	if businessProfile.ID == 0 {
		return external_models.BusinessProfile{}, fmt.Errorf("Business lacks a profile.")
	}
	return businessProfile, nil
}

func getBusinessChargeWithBusinessIDAndCurrency(extReq request.ExternalRequest, businessID int, currency string) (external_models.BusinessCharge, error) {
	dataInterface, err := extReq.SendExternalRequest(request.GetBusinessCharge, external_models.GetBusinessChargeModel{
		BusinessID: uint(businessID),
		Currency:   strings.ToUpper(currency),
	})

	if err != nil {
		extReq.Logger.Info(err.Error())
		return external_models.BusinessCharge{}, err
	}

	businessCharge, ok := dataInterface.(external_models.BusinessCharge)
	if !ok {
		return external_models.BusinessCharge{}, fmt.Errorf("response data format error")
	}

	if businessCharge.ID == 0 {
		return external_models.BusinessCharge{}, fmt.Errorf("business charge not found")
	}

	return businessCharge, nil
}

func getBusinessChargeWithBusinessIDAndCountry(extReq request.ExternalRequest, businessID int, country string) (external_models.BusinessCharge, error) {
	dataInterface, err := extReq.SendExternalRequest(request.GetBusinessCharge, external_models.GetBusinessChargeModel{
		BusinessID: uint(businessID),
		Country:    strings.ToUpper(country),
	})

	if err != nil {
		extReq.Logger.Info(err.Error())
		return external_models.BusinessCharge{}, err
	}

	businessCharge, ok := dataInterface.(external_models.BusinessCharge)
	if !ok {
		return external_models.BusinessCharge{}, fmt.Errorf("response data format error")
	}

	if businessCharge.ID == 0 {
		return external_models.BusinessCharge{}, fmt.Errorf("business charge not found")
	}

	return businessCharge, nil
}

func initBusinessCharge(extReq request.ExternalRequest, businessID int, currency string) (external_models.BusinessCharge, error) {
	dataInterface, err := extReq.SendExternalRequest(request.InitBusinessCharge, external_models.InitBusinessChargeModel{
		BusinessID: uint(businessID),
		Currency:   strings.ToUpper(currency),
	})

	if err != nil {
		extReq.Logger.Info(err.Error())
		return external_models.BusinessCharge{}, err
	}

	businessCharge, ok := dataInterface.(external_models.BusinessCharge)
	if !ok {
		return external_models.BusinessCharge{}, fmt.Errorf("response data format error")
	}

	if businessCharge.ID == 0 {
		return external_models.BusinessCharge{}, fmt.Errorf("business charge init failed")
	}

	return businessCharge, nil
}

func getCountryByCurrency(extReq request.ExternalRequest, logger *utility.Logger, currencyCode string) (external_models.Country, error) {

	countryInterface, err := extReq.SendExternalRequest(request.GetCountry, external_models.GetCountryModel{
		CurrencyCode: currencyCode,
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
