package payment

import (
	"fmt"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
)

type Rave struct {
	ExtReq request.ExternalRequest
}

func (r *Rave) ListBanks(countryCode string) ([]external_models.BanksResponse, error) {
	banksItf, err := r.ExtReq.SendExternalRequest(request.ListBanksWithRave, strings.ToUpper(countryCode))
	if err != nil {
		return []external_models.BanksResponse{}, err
	}

	banks, ok := banksItf.([]external_models.BanksResponse)
	if !ok {
		return []external_models.BanksResponse{}, fmt.Errorf("response data format error")
	}

	return banks, nil
}

func (r *Rave) GetBank(countryCode, bankName string) (external_models.BanksResponse, error) {
	banks, err := r.ListBanks(countryCode)
	if err != nil {
		return external_models.BanksResponse{}, err
	}

	for _, bank := range banks {
		if strings.ToLower(bank.Name) == strings.ToLower(bankName) {
			return bank, nil
		}
	}

	return external_models.BanksResponse{}, fmt.Errorf("bank not found")
}

func (r *Rave) ResolveAccount(bankCode, accountNumber string) (string, error) {
	bankItf, err := r.ExtReq.SendExternalRequest(request.RaveResolveBankAccount, external_models.ResolveAccountRequest{AccountBank: bankCode, AccountNumber: accountNumber})
	if err != nil {
		return "", err
	}

	bankName, ok := bankItf.(string)
	if !ok {
		return "", fmt.Errorf("response data format error")
	}

	return bankName, nil
}

func (r *Rave) ConvertCurrency(amount float64, from, to string) (models.ConvertCurrencyResponse, error) {
	conversionItf, err := r.ExtReq.SendExternalRequest(request.ConvertCurrencyWithRave, external_models.ConvertCurrencyRequest{Amount: amount, From: from, To: to})
	if err != nil {
		return models.ConvertCurrencyResponse{}, err
	}

	conversionData, ok := conversionItf.(external_models.ConvertCurrencyData)
	if !ok {
		return models.ConvertCurrencyResponse{}, fmt.Errorf("response data format error")
	}
	var converted float64 = conversionData.Source.Amount
	var rate float64 = conversionData.Rate

	if strings.ToUpper(from) == "USD" && strings.ToUpper(to) == "NGN" {
		rate = conversionData.Rate - 5
		converted = amount * rate
	}

	return models.ConvertCurrencyResponse{
		Converted: converted,
		Rate:      rate,
		From:      strings.ToUpper(from),
		To:        strings.ToUpper(to),
		Amount:    amount,
	}, nil
}

func (r *Rave) InitPayment(reference, email, currency, redirectUrl string, amount float64) (string, external_models.RaveInitPaymentRequest, error) {
	data := external_models.RaveInitPaymentRequest{
		TxRef: reference,
		Customer: struct {
			Email string "json:\"email\""
		}{Email: email},
		Currency:    strings.ToUpper(currency),
		Amount:      amount,
		RedirectUrl: redirectUrl,
	}
	paymentItf, err := r.ExtReq.SendExternalRequest(request.RaveInitPayment, data)
	if err != nil {
		return "", data, err
	}

	paymentData, ok := paymentItf.(external_models.RaveInitPaymentResponse)
	if !ok {
		return "", data, fmt.Errorf("response data format error")
	}
	return paymentData.Data.Link, data, nil
}

func (r *Rave) ReserveAccount(reference, narration, email, firstName, lastName string, amount float64) (external_models.RaveReserveAccountResponseData, error) {
	data := external_models.RaveReserveAccountRequest{
		TxRef:       reference,
		Narration:   narration,
		Amount:      amount,
		Email:       email,
		Frequency:   1,
		Firstname:   firstName,
		Lastname:    lastName,
		IsPermanent: false,
	}
	paymentItf, err := r.ExtReq.SendExternalRequest(request.RaveReserveAccount, data)
	if err != nil {
		return external_models.RaveReserveAccountResponseData{}, err
	}

	paymentData, ok := paymentItf.(external_models.RaveReserveAccountResponseData)
	if !ok {
		return external_models.RaveReserveAccountResponseData{}, fmt.Errorf("response data format error")
	}
	return paymentData, nil
}
