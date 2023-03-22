package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PaymentAccount struct {
	ID                   uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	PaymentAccountID     string    `gorm:"column:payment_account_id; type:varchar(255); not null" json:"payment_account_id"`
	PaymentID            string    `gorm:"column:payment_id; type:varchar(255)" json:"payment_id"`
	TransactionID        string    `gorm:"column:transaction_id; type:varchar(255)" json:"transaction_id"`
	AccountNumber        string    `gorm:"column:account_number; type:varchar(255); not null" json:"account_number"`
	BankCode             string    `gorm:"column:bank_code; type:varchar(255); not null" json:"bank_code"`
	ExpiresAfter         string    `gorm:"column:expires_after; type:varchar(255); not null" json:"expires_after"`
	IsUsed               bool      `gorm:"column:is_used;not null" json:"is_used"`
	CreatedAt            time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	BankName             string    `gorm:"column:bank_name; type:varchar(255)" json:"bank_name"`
	ReservationReference string    `gorm:"column:reservation_reference; type:varchar(255)" json:"reservation_reference"`
	Status               string    `gorm:"column:status; type:varchar(255)" json:"status"`
	AccountName          string    `gorm:"column:account_name; type:varchar(255)" json:"account_name"`
	BusinessID           string    `gorm:"column:business_id; type:varchar(255)" json:"business_id"`
	PaymentReference     string    `gorm:"column:paymentReference; type:varchar(255)" json:"paymentReference"`
}

func (p *PaymentAccount) CreatePaymentAccount(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("payment accounts creation failed: %v", err.Error())
	}
	return nil
}

func (p *PaymentAccount) GetPaymentAccountByBusinessIDAndTransactionID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "business_id = ? and transaction_id=?", p.BusinessID, p.TransactionID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (p *PaymentAccount) GetPaymentAccountByPaymentAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "payment_account_id = ?", p.PaymentAccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (p *PaymentAccount) GetPaymentAccountByPaymentID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "payment_id = ?", p.PaymentID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *PaymentAccount) GetPaymentAccountByBusinessID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "business_id = ?", p.BusinessID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *PaymentAccount) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &p)
	return err
}

func (p *PaymentAccount) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &p)
	if err != nil {
		return fmt.Errorf("payment account delete failed: %v", err.Error())
	}
	return nil
}
