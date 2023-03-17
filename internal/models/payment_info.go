package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PaymentInfo struct {
	ID          uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	PaymentID   string    `gorm:"column:payment_id; type:varchar(255); not null" json:"payment_id"`
	Reference   string    `gorm:"column:reference; type:varchar(255); not null" json:"reference"`
	Status      string    `gorm:"column:status; type:varchar(255); not null" json:"status"`
	Gateway     string    `gorm:"column:gateway; type:varchar(255); not null" json:"gateway"`
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt   time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	RedirectUrl string    `gorm:"column:redirecturl; type:text" json:"redirecturl"`
	FailUrl     string    `gorm:"column:failurl; type:varchar(255)" json:"failurl"`
}

func (p *PaymentInfo) CreatePaymentInfo(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payment info creation failed: %v", err.Error())
	}
	return nil
}

func (p *PaymentInfo) GetPaymentInfoByReference(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "reference = ?", p.Reference)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (p *PaymentInfo) GetPaymentInfoByPaymentID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "payment_id = ?", p.PaymentID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *PaymentInfo) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &p)
	return err
}
