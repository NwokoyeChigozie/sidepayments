package models

import (
	"fmt"
	"time"

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

func (p *Payment) CreatePayment(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payment creation failed: %v", err.Error())
	}
	return nil
}
