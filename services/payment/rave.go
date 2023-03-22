package payment

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
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
		if strings.EqualFold(bank.Name, bankName) {
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

func (r *Rave) VerifyTrans(reference string, amount float64, currency string) (string, error) {
	paymentItf, err := r.ExtReq.SendExternalRequest(request.RaveVerifyTransactionByTxRef, reference)
	if err != nil {
		return "pending", err
	}

	data, ok := paymentItf.(external_models.RaveVerifyTransactionResponseData)
	if !ok {
		return "pending", fmt.Errorf("response data format error")
	}

	if data.Status == "error" || data.Status == "" {
		return "error", fmt.Errorf("error occured verifying transaction")
	}

	if data.Status == "failed" {
		return "error", fmt.Errorf("transaction failed")
	}

	if !strings.EqualFold(data.Currency, currency) {
		return "error", fmt.Errorf("different currencies")
	}

	if data.ChargedAmount < amount {
		return "pending", fmt.Errorf("incomplete payment")
	}

	return "success", nil
}

func (r *Rave) StatusV3(db postgresql.Databases, payment models.Payment, paymentInfo models.PaymentInfo, reference string) (bool, float64, error) {
	var (
		status bool
		amount float64
	)

	paymentItf, err := r.ExtReq.SendExternalRequest(request.RaveVerifyTransactionByTxRef, reference)
	if err != nil {
		return status, amount, err
	}

	data, ok := paymentItf.(external_models.RaveVerifyTransactionResponseData)
	if !ok {
		return status, amount, fmt.Errorf("response data format error")
	}

	if data.Card != nil {
		expirySlice := strings.Split(data.Card.Expiry, "/")
		expiryMonth := ""
		expiryYear := ""
		if len(expirySlice) >= 2 {
			expiryMonth = expirySlice[0]
			expiryYear = expirySlice[1]
		}

		cardByte, _ := json.Marshal(data.Card)

		paymentCardInfo := models.PaymentCardInfo{AccountID: int(payment.AccountID), LastFourDigits: data.Card.Last4digits, Brand: data.Card.Type}
		_, err := paymentCardInfo.GetPaymentCardInfoByAccountIDLast4DigitsAndBrand(db.Payment)
		if err != nil {
			paymentCardInfo = models.PaymentCardInfo{
				AccountID:         int(payment.AccountID),
				PaymentID:         paymentInfo.PaymentID,
				CcExpiryMonth:     expiryMonth,
				CcExpiryYear:      expiryYear,
				LastFourDigits:    data.Card.Last4digits,
				Brand:             data.Card.Type,
				IssuingCountry:    data.Card.Country,
				CardToken:         string(cardByte),
				CardLifeTimeToken: data.Card.Token,
				Payload:           string(cardByte),
			}
			err := paymentCardInfo.CreatePaymentCardInfo(db.Payment)
			if err != nil {
				return status, amount, err
			}
		}
	}

	if data.Status == "successful" || data.Status == "completed" {
		status = true
		amount = data.ChargedAmount
	} else if data.Status == "failed" {
		status = false
		amount = data.ChargedAmount
	} else if data.Status == "error" {
		status = false
		amount = data.ChargedAmount
	} else {
		status = true
		amount = data.ChargedAmount
	}

	return status, amount, nil
}

func (r *Rave) ChargeCard(token, currency, email, reference string, amount float64) (string, error) {
	data := external_models.RaveChargeCardRequest{
		Token:    token,
		Currency: strings.ToUpper(currency),
		Email:    email,
		TxRef:    reference,
		Amount:   amount,
	}
	paymentItf, err := r.ExtReq.SendExternalRequest(request.RaveChargeCard, data)
	if err != nil {
		return "failed", err
	}

	paymentData, ok := paymentItf.(external_models.RaveVerifyTransactionResponseData)
	if !ok {
		return "failed", fmt.Errorf("response data format error")
	}

	if strings.ToLower(paymentData.Status) != "successful" {
		return "failed", nil
	}

	return "success", nil
}

func (r *Rave) InitTransfer(bank, accountNo string, amount float64, narration, currency, reference, callback string) (external_models.RaveInitTransferResponse, error) {
	data := external_models.RaveInitTransferRequest{
		AccountBank:   bank,
		AccountNumber: accountNo,
		Amount:        amount,
		Narration:     narration,
		Currency:      strings.ToUpper(currency),
		Reference:     reference,
		DebitCurrency: strings.ToUpper(currency),
		CallbackUrl:   callback,
	}
	paymentItf, err := r.ExtReq.SendExternalRequest(request.RaveInitTransfer, data)
	if err != nil {
		return external_models.RaveInitTransferResponse{}, err
	}

	paymentData, ok := paymentItf.(external_models.RaveInitTransferResponse)
	if !ok {
		return paymentData, fmt.Errorf("response data format error")
	}

	if strings.ToLower(paymentData.Status) != "successful" {
		return paymentData, nil
	}

	return paymentData, nil
}
