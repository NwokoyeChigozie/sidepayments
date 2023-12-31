package payment

import (
	"fmt"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/config"
)

type Monnify struct {
	ExtReq request.ExternalRequest
}

func (m *Monnify) InitPayment(amount float64, customerName, customerEmail, reference, description, currency, redirectUrl string) (string, external_models.MonnifyInitPaymentRequest, error) {
	var (
		contractCode = config.GetConfig().Monnify.MonnifyContractCode
	)
	data := external_models.MonnifyInitPaymentRequest{
		Amount:             amount,
		CustomerName:       customerName,
		CustomerEmail:      customerEmail,
		PaymentReference:   reference,
		PaymentDescription: description,
		CurrencyCode:       strings.ToUpper(currency),
		ContractCode:       contractCode,
		RedirectUrl:        redirectUrl,
	}
	paymentItf, err := m.ExtReq.SendExternalRequest(request.MonnifyInitPayment, data)
	if err != nil {
		return "", data, err
	}

	paymentData, ok := paymentItf.(external_models.MonnifyInitPaymentResponseBody)
	if !ok {
		return "", data, fmt.Errorf("response data format error")
	}
	return paymentData.CheckoutUrl, data, nil
}

func (m *Monnify) Status(reference string) (external_models.MonnifyVerifyByReferenceResponseBody, bool, string, float64, error) {
	var (
		status       bool
		amount       float64
		statusString string
	)
	paymentItf, err := m.ExtReq.SendExternalRequest(request.MonnifyVerifyTransactionByReference, reference)
	if err != nil {
		return external_models.MonnifyVerifyByReferenceResponseBody{}, status, statusString, amount, err
	}

	data, ok := paymentItf.(external_models.MonnifyVerifyByReferenceResponseBody)
	if !ok {
		return data, status, statusString, amount, fmt.Errorf("response data format error")
	}

	if strings.ToUpper(data.PaymentStatus) == "PAID" {
		status = true
		amount = data.Amount
	} else {
		status = false
		amount = data.Amount
	}

	return data, status, statusString, amount, nil
}

func (m *Monnify) VerifyTrans(accountReference string, amount float64) (bool, float64, error) {
	var (
		status     bool
		paidAmount float64
	)

	transactions, err := m.FetchAccountTrans(accountReference)
	if err != nil {
		return status, 0, err
	}

	for _, v := range transactions {
		if strings.ToUpper(v.PaymentStatus) == "PAID" {
			paidAmount += v.AmountPaid
		}
	}

	if paidAmount >= amount && paidAmount > 0 {
		status = true
	} else {
		status = false
	}

	return status, paidAmount, nil
}

func (m *Monnify) ReserveAccount(reference, accountName, currencyCode, customerEmail string) (external_models.MonnifyReserveAccountResponseBody, error) {
	paymentItf, err := m.ExtReq.SendExternalRequest(request.MonnifyReserveAccount, external_models.MonnifyReserveAccountRequest{
		AccountReference: reference,
		AccountName:      accountName,
		CurrencyCode:     strings.ToUpper(currencyCode),
		ContractCode:     config.GetConfig().Monnify.MonnifyContractCode,
		CustomerEmail:    customerEmail,
	})
	if err != nil {
		return external_models.MonnifyReserveAccountResponseBody{}, err
	}

	data, ok := paymentItf.(external_models.MonnifyReserveAccountResponseBody)
	if !ok {
		return external_models.MonnifyReserveAccountResponseBody{}, fmt.Errorf("response data format error")
	}

	return data, nil
}

func (m *Monnify) FetchAccountTrans(reference string) ([]external_models.GetMonnifyReserveAccountTransactionsResponseBodyContent, error) {

	paymentItf, err := m.ExtReq.SendExternalRequest(request.GetMonnifyReserveAccountTransactions, reference)
	if err != nil {
		return []external_models.GetMonnifyReserveAccountTransactionsResponseBodyContent{}, err
	}

	data, ok := paymentItf.(external_models.GetMonnifyReserveAccountTransactionsResponseBody)
	if !ok {
		return []external_models.GetMonnifyReserveAccountTransactionsResponseBodyContent{}, fmt.Errorf("response data format error")
	}

	return data.Content, nil
}

func (m *Monnify) InitTransfer(amount float64, reference, narration, destinationBankCode, destinationAccountNo, currency, destinationAccountName string) (external_models.MonnifyInitTransferResponse, error) {
	reqData := external_models.MonnifyInitTransferRequest{
		Amount:                   amount,
		Reference:                reference,
		Narration:                narration,
		DestinationBankCode:      destinationBankCode,
		DestinationAccountNumber: destinationAccountNo,
		Currency:                 strings.ToUpper(currency),
		SourceAccountNumber:      config.GetConfig().Monnify.MonnifyDisbursementAccount,
		DestinationAccountName:   destinationAccountName,
	}
	paymentItf, err := m.ExtReq.SendExternalRequest(request.MonnifyInitTransfer, reqData)
	if err != nil {
		return external_models.MonnifyInitTransferResponse{}, err
	}

	data, ok := paymentItf.(external_models.MonnifyInitTransferResponse)
	if !ok {
		return external_models.MonnifyInitTransferResponse{}, fmt.Errorf("response data format error")
	}

	return data, nil
}
