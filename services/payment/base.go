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

type ExchangeTransactionStatus string

var (
	ExchangeTransactionCompleted ExchangeTransactionStatus = "completed"
	ExchangeTransactionPending   ExchangeTransactionStatus = "pending"
	ExchangeTransactionFailed    ExchangeTransactionStatus = "failed"
)

func ListTransactionsByID(extReq request.ExternalRequest, transactionID string) (external_models.TransactionByID, error) {

	transactionInterface, err := extReq.SendExternalRequest(request.ListTransactionsByID, transactionID)

	if err != nil {
		extReq.Logger.Error(err.Error())
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
func GetUserWithEmail(extReq request.ExternalRequest, email string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{EmailAddress: email})
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
		logger.Error(err.Error())
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
		logger.Error(err.Error())
		return external_models.Country{}, fmt.Errorf("your country could not be resolved, please update your profile")
	}
	country, ok := countryInterface.(external_models.Country)
	if !ok {
		return external_models.Country{}, fmt.Errorf("response data format error")
	}
	if country.ID == 0 {
		return external_models.Country{}, fmt.Errorf("your country could not be resolved, please update your profile")
	}

	return country, nil
}

func isRequestIPNigerian(extReq request.ExternalRequest, c *gin.Context) (bool, error) {
	if extReq.Test {
		return true, nil
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return false, err
	}

	ipInterface, err := extReq.SendExternalRequest(request.ResolveIP, ip)
	if err != nil {
		extReq.Logger.Error(err.Error())
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
		logger.Error(err.Error())
		return external_models.BusinessProfile{}, fmt.Errorf("business lacks a profile")
	}

	businessProfile, ok := businessProfileInterface.(external_models.BusinessProfile)
	if !ok {
		return external_models.BusinessProfile{}, fmt.Errorf("response data format error")
	}

	if businessProfile.ID == 0 {
		return external_models.BusinessProfile{}, fmt.Errorf("business lacks a profile")
	}
	return businessProfile, nil
}
func GetBusinessProfileByFlutterwaveMerchantID(extReq request.ExternalRequest, logger *utility.Logger, merchantID string) (external_models.BusinessProfile, error) {
	businessProfileInterface, err := extReq.SendExternalRequest(request.GetBusinessProfile, external_models.GetBusinessProfileModel{
		FlutterwaveMerchantID: merchantID,
	})
	if err != nil {
		logger.Error(err.Error())
		return external_models.BusinessProfile{}, err
	}

	businessProfile, ok := businessProfileInterface.(external_models.BusinessProfile)
	if !ok {
		return external_models.BusinessProfile{}, fmt.Errorf("response data format error")
	}

	if businessProfile.ID == 0 {
		return external_models.BusinessProfile{}, fmt.Errorf("no business profile found")
	}
	return businessProfile, nil
}

func getBusinessChargeWithBusinessIDAndCurrency(extReq request.ExternalRequest, businessID int, currency string) (external_models.BusinessCharge, error) {
	dataInterface, err := extReq.SendExternalRequest(request.GetBusinessCharge, external_models.GetBusinessChargeModel{
		BusinessID: uint(businessID),
		Currency:   strings.ToUpper(currency),
	})

	if err != nil {
		extReq.Logger.Error(err.Error())
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
		extReq.Logger.Error(err.Error())
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
		extReq.Logger.Error(err.Error())
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

func GetCountryByCurrency(extReq request.ExternalRequest, logger *utility.Logger, currencyCode string) (external_models.Country, error) {

	countryInterface, err := extReq.SendExternalRequest(request.GetCountry, external_models.GetCountryModel{
		CurrencyCode: currencyCode,
	})

	if err != nil {
		logger.Error(err.Error())
		return external_models.Country{}, fmt.Errorf("your country could not be resolved, please update your profile")
	}
	country, ok := countryInterface.(external_models.Country)
	if !ok {
		return external_models.Country{}, fmt.Errorf("response data format error")
	}
	if country.ID == 0 {
		return external_models.Country{}, fmt.Errorf("your country could not be resolved, please update your profile")
	}

	return country, nil
}

func GetRateByID(extReq request.ExternalRequest, rateID int) (external_models.Rate, error) {

	rateInterface, err := extReq.SendExternalRequest(request.GetRateByID, rateID)

	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.Rate{}, err
	}
	rate, ok := rateInterface.(external_models.Rate)
	if !ok {
		return rate, fmt.Errorf("response data format error")
	}
	if rate.ID == 0 {
		return rate, fmt.Errorf("rate with id %v not found", rateID)
	}

	return rate, nil
}

func CreateExchangeTransaction(extReq request.ExternalRequest, accountID, rateID int, initialAmount, finalAmount float64, status ExchangeTransactionStatus) error {

	_, err := extReq.SendExternalRequest(request.CreateExchangeTransaction, external_models.CreateExchangeTransactionRequest{
		AccountID:     accountID,
		InitialAmount: initialAmount,
		FinalAmount:   finalAmount,
		RateID:        rateID,
		Status:        string(status),
	})

	if err != nil {
		extReq.Logger.Error(err.Error())
		return err
	}

	return nil
}
func GetUserCredentialByAccountIdAndType(extReq request.ExternalRequest, accountID int, iType string) (external_models.UsersCredential, error) {

	userCredInterface, err := extReq.SendExternalRequest(request.GetUserCredential, external_models.GetUserCredentialModel{
		AccountID:          uint(accountID),
		IdentificationType: iType,
	})

	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.UsersCredential{}, err
	}

	userCred, ok := userCredInterface.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCred.Data, fmt.Errorf("response data format error")
	}
	if userCred.Data.ID == 0 {
		return userCred.Data, fmt.Errorf("user credential not found")
	}

	return userCred.Data, nil
}
func GetBankDetail(extReq request.ExternalRequest, id, accountID int, country, currency string) (external_models.BankDetail, error) {

	data := external_models.GetBankDetailModel{}

	if id != 0 {
		data = external_models.GetBankDetailModel{
			ID: uint(id),
		}
	} else {
		data = external_models.GetBankDetailModel{
			AccountID: uint(accountID),
			Country:   country,
			Currency:  currency,
		}
	}
	userBankInterface, err := extReq.SendExternalRequest(request.GetBankDetails, data)

	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.BankDetail{}, err
	}

	bankDetail, ok := userBankInterface.(external_models.BankDetail)
	if !ok {
		return bankDetail, fmt.Errorf("response data format error")
	}
	if bankDetail.ID == 0 {
		return bankDetail, fmt.Errorf("bank detail not found")
	}

	return bankDetail, nil
}

func GetBank(extReq request.ExternalRequest, id int, name, code, country string) (external_models.Bank, error) {

	data := external_models.GetBankRequest{
		ID:      uint(id),
		Name:    name,
		Code:    code,
		Country: country,
	}
	bankInterface, err := extReq.SendExternalRequest(request.GetBank, data)

	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.Bank{}, err
	}

	bank, ok := bankInterface.(external_models.Bank)
	if !ok {
		return bank, fmt.Errorf("response data format error")
	}
	if bank.ID == 0 {
		return bank, fmt.Errorf("bank not found")
	}

	return bank, nil
}

func HasBvn(extReq request.ExternalRequest, accountID uint) bool {
	userCredential, err := GetUserCredentialByAccountIdAndType(extReq, int(accountID), "bvn")
	if err != nil {
		return false
	}
	return userCredential.Bvn != ""
}

func thisOrThatStr(this, that string) string {
	if this == "" {
		return that
	}
	return this
}
