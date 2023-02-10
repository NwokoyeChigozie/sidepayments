package payment

import (
	"fmt"
	"net/http"

	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/internal/models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
)

func ListPaymentByTransactionIDService(extReq request.ExternalRequest, db postgresql.Databases, transactionID string) (models.ListPaymentsResponse, int, error) {
	var (
		resp = models.ListPaymentsResponse{}
	)

	transaction, err := ListTransactionsByID(extReq, transactionID)
	if err != nil {
		return models.ListPaymentsResponse{}, http.StatusBadRequest, err
	}
	resp.Transaction = transaction

	_, ok := transaction.Parties["buyer"]
	if !ok {
		return resp, http.StatusBadRequest, fmt.Errorf("transaction lacks a Buyer Party")
	}

	payment := models.Payment{TransactionID: transactionID}
	code, err := payment.GetPaymentByTransactionID(db.Payment)
	if err != nil {
		return resp, code, err
	}

	escrowCharge := 0
	brokerCharge := 0
	shippingFee := 0

	resp.Payment = models.ListPayment{
		ID:               payment.ID,
		PaymentID:        payment.PaymentID,
		TransactionID:    payment.TransactionID,
		TotalAmount:      payment.TotalAmount,
		EscrowCharge:     payment.EscrowCharge,
		IsPaid:           payment.IsPaid,
		PaymentMadeAt:    payment.PaymentMadeAt,
		DeletedAt:        payment.DeletedAt,
		CreatedAt:        payment.CreatedAt,
		UpdatedAt:        payment.UpdatedAt,
		AccountID:        payment.AccountID,
		BusinessID:       payment.BusinessID,
		Currency:         payment.Currency,
		ShippingFee:      payment.ShippingFee,
		DisburseCurrency: payment.DisburseCurrency,
		PaymentType:      payment.PaymentType,
		BrokerCharge:     payment.BrokerCharge,
		SummedAmount:     payment.TotalAmount + float64(shippingFee) + float64(brokerCharge) + float64(escrowCharge),
	}
	return resp, http.StatusOK, nil
}
