package models

import (
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Disbursement struct {
	ID                    uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	DisbursementID        int       `gorm:"column:disbursement_id; type:int; not null" json:"disbursement_id"`
	RecipientID           int       `gorm:"column:recipient_id; type:int; not null" json:"recipient_id"`
	PaymentID             string    `gorm:"column:payment_id; type:varchar(255); not null" json:"payment_id"`
	BusinessID            int       `gorm:"column:business_id; type:int; not null" json:"business_id"`
	Amount                string    `gorm:"column:amount; type:varchar(255); not null" json:"amount"`
	Narration             string    `gorm:"column:narration; type:varchar(255)" json:"narration"`
	Currency              string    `gorm:"column:currency; type:varchar(255); not null" json:"currency"`
	Reference             string    `gorm:"column:reference; type:varchar(255); not null" json:"reference"`
	CallbackUrl           string    `gorm:"column:callback_url; type:varchar(255)" json:"callback_url"`
	BeneficiaryName       string    `gorm:"column:beneficiary_name; type:varchar(255)" json:"beneficiary_name"`
	DestinationBranchCode string    `gorm:"column:destination_branch_code; type:varchar(255)" json:"destination_branch_code"`
	DebitCurrency         string    `gorm:"column:debit_currency; type:varchar(255)" json:"debit_currency"`
	Gateway               string    `gorm:"column:gateway; type:varchar(255)" json:"gateway"`
	Type                  string    `gorm:"column:type; type:varchar(255)" json:"type"`
	Status                string    `gorm:"column:status; type:varchar(255); default: new" json:"status"`
	PaymentReleasedAt     string    `gorm:"column:payment_released_at; type:varchar(255)" json:"payment_released_at"`
	DeletedAt             time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt             time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Fee                   int       `gorm:"column:fee; type:int; default: 0" json:"fee"`
	Tries                 int       `gorm:"column:tries; type:int; default: 0" json:"tries"`
	TryAgainAt            time.Time `gorm:"column:try_again_at" json:"try_again_at"`
	BankAccountNumber     string    `gorm:"column:bank_account_number; type:varchar(255)" json:"bank_account_number"`
	BankName              string    `gorm:"column:bank_name; type:varchar(255)" json:"bank_name"`
}

type WalletTransferRequest struct {
	SenderAccountID    int     `json:"sender_account_id"  validate:"required" pgvalidate:"exists=auth$users$account_id"`
	RecipientAccountID int     `json:"recipient_account_id"  validate:"required" pgvalidate:"exists=auth$users$account_id"`
	InitialAmount      float64 `json:"initial_amount" validate:"gt=0"`
	FinalAmount        float64 `json:"final_amount" validate:"required,gt=0"`
	RateID             int     `json:"rate_id" pgvalidate:"exists=transaction$rates$id"`
	SenderCurrency     string  `json:"sender_currency" validate:"required"`
	RecipientCurrency  string  `json:"recipient_currency" validate:"required"`
	TransactionID      string  `json:"transaction_id" pgvalidate:"exists=transaction$transactions$transaction_id"`
	Refund             bool    `json:"refund"`
}

func (d *Disbursement) GetDisbursementsByRecipientID(db *gorm.DB, paginator postgresql.Pagination) ([]Disbursement, postgresql.PaginationResponse, error) {
	details := []Disbursement{}
	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, "recipient_id = ?", d.RecipientID)
	if err != nil {
		return details, pagination, err
	}
	return details, pagination, nil
}

func (d *Disbursement) GetDisbursementsByAccountID(db *gorm.DB, paginator postgresql.Pagination) ([]Disbursement, postgresql.PaginationResponse, error) {
	details := []Disbursement{}
	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, "account_id = ?", d.BusinessID)
	if err != nil {
		return details, pagination, err
	}
	return details, pagination, nil
}
