package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Payment struct {
	ID               int64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	PaymentID        string    `gorm:"column:payment_id" json:"payment_id"`
	TransactionID    string    `gorm:"column:transaction_id" json:"transaction_id"`
	TotalAmount      float64   `gorm:"column:total_amount" json:"total_amount"`
	EscrowCharge     float64   `gorm:"column:escrow_charge" json:"escrow_charge"`
	IsPaid           bool      `gorm:"column:is_paid" json:"is_paid"`
	PaymentMadeAt    time.Time `gorm:"column:payment_made_at" json:"payment_made_at"`
	DeletedAt        time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt        time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	AccountID        int64     `gorm:"column:account_id" json:"account_id"`
	BusinessID       int64     `gorm:"column:business_id" json:"business_id"`
	Currency         string    `gorm:"column:currency" json:"currency"`
	ShippingFee      float64   `gorm:"column:shipping_fee" json:"shipping_fee"`
	DisburseCurrency string    `gorm:"column:disburse_currency" json:"disburse_currency"`
	PaymentType      string    `gorm:"column:payment_type" json:"payment_type"`
	BrokerCharge     float64   `gorm:"column:broker_charge" json:"broker_charge"`
}

type CreatePaymentRequest struct {
	TransactionID string  `json:"transaction_id"  validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	TotalAmount   float64 `json:"total_amount"  validate:"required"`
	ShippingFee   float64 `json:"shipping_fee"`
	BrokerCharge  float64 `json:"broker_charge"`
	EscrowCharge  float64 `json:"escrow_charge"`
	Currency      string  `json:"currency"`
}
type InitiatePaymentRequest struct {
	TransactionID  string `json:"transaction_id"  validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	SuccessPage    string `json:"success_page" validate:"required"`
	PaymentGateway string `json:"payment_gateway"`
}
type InitiatePaymentResponse struct {
	PaymentStatus string `json:"payment_status"`
	Link          string `json:"link"`
	Ref           string `json:"ref"`
	ExternalRef   string `json:"external_ref"`
	TransactionID string `json:"transaction_id"`
}

type CreatePaymentHeadlessRequest struct {
	AccountID     int     `json:"account_id"  validate:"required" pgvalidate:"exists=auth$users$account_id"`
	TotalAmount   float64 `json:"total_amount"  validate:"required"`
	EscrowCharge  float64 `json:"escrow_charge"`
	Currency      string  `json:"currency"`
	PaymentMadeAt string  `json:"payment_made_at"`
}
type EditPaymentRequest struct {
	PaymentID    string  `json:"payment_id"  validate:"required" pgvalidate:"exists=payment$payments$payment_id"`
	EscrowCharge float64 `json:"escrow_charge"  validate:"required"`
}
type VerifyTransactionPaymentRequest struct {
	TransactionID string `json:"transaction_id"  validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
}
type VerifyTransactionPaymentResponse struct {
	Status string    `json:"status"`
	IsPaid bool      `json:"is_paid"`
	Amount float64   `json:"amount"`
	Charge float64   `json:"charge"`
	Date   time.Time `json:"date"`
}
type ConvertCurrencyResponse struct {
	Converted float64 `json:"converted"`
	Rate      float64 `json:"rate"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
}

type ListPaymentsRequest struct {
	TransactionID string `json:"transaction_id"  validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
}

type ListPayment struct {
	ID               int64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	PaymentID        string    `gorm:"column:payment_id" json:"payment_id"`
	TransactionID    string    `gorm:"column:transaction_id" json:"transaction_id"`
	TotalAmount      float64   `gorm:"column:total_amount" json:"total_amount"`
	EscrowCharge     float64   `gorm:"column:escrow_charge" json:"escrow_charge"`
	IsPaid           bool      `gorm:"column:is_paid" json:"is_paid"`
	PaymentMadeAt    time.Time `gorm:"column:payment_made_at" json:"payment_made_at"`
	DeletedAt        time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt        time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	AccountID        int64     `gorm:"column:account_id" json:"account_id"`
	BusinessID       int64     `gorm:"column:business_id" json:"business_id"`
	Currency         string    `gorm:"column:currency" json:"currency"`
	ShippingFee      float64   `gorm:"column:shipping_fee" json:"shipping_fee"`
	DisburseCurrency string    `gorm:"column:disburse_currency" json:"disburse_currency"`
	PaymentType      string    `gorm:"column:payment_type" json:"payment_type"`
	BrokerCharge     float64   `gorm:"column:broker_charge" json:"broker_charge"`
	SummedAmount     float64   `gorm:"column:summed_amount" json:"summed_amount"`
}

type ListPaymentsResponse struct {
	Transaction external_models.TransactionByID `json:"transaction"`
	Payment     ListPayment                     `json:"payment"`
}

func (p *Payment) CreatePayment(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payment creation failed: %v", err.Error())
	}
	return nil
}

func (p *Payment) GetPaymentByTransactionID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectLatestFromDb(db, &p, "transaction_id = ?", p.TransactionID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (p *Payment) GetPaymentByPaymentID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectLatestFromDb(db, &p, "payment_id = ?", p.PaymentID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (p *Payment) GetPaymentByTransactionIDAndNotPaymentMadeAt(db *gorm.DB) (int, error) {
	var baseTime time.Time
	err, nilErr := postgresql.SelectLatestFromDb(db, &p, "transaction_id = ? and (payment_made_at is not null and payment_made_at!=?)", p.TransactionID, baseTime)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *Payment) GetPaymentsByAccountIDAndNullTransactionID(db *gorm.DB, paginator postgresql.Pagination) ([]Payment, postgresql.PaginationResponse, error) {
	details := []Payment{}
	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, "account_id = ? and (transaction_id is null or transaction_id='')", p.ID)
	if err != nil {
		return details, pagination, err
	}
	return details, pagination, nil
}
func (p *Payment) GetPaymentsByTransactionIDAndIsPaid(db *gorm.DB, paginator postgresql.Pagination) ([]Payment, postgresql.PaginationResponse, error) {
	details := []Payment{}
	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, "transaction_id = ? and is_paid = ?", p.TransactionID, p.IsPaid)
	if err != nil {
		return details, pagination, err
	}
	return details, pagination, nil
}
func (p *Payment) GetAllPaymentsByAccountIDAndIsPaidAndPaymentMadeAtNotNull(db *gorm.DB) ([]Payment, error) {
	details := []Payment{}
	var baseTime time.Time
	err := postgresql.SelectAllFromDb(db, "desc", &details, "account_id = ? and is_paid = ? and (payment_made_at is not null and payment_made_at!=?)", p.AccountID, p.IsPaid, baseTime)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (p *Payment) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &p)
	return err
}
func (p *Payment) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &p)
	if err != nil {
		return fmt.Errorf("payment delete failed: %v", err.Error())
	}
	return nil
}
