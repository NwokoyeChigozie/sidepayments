package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type FailedDisbursement struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	PaymentID string    `gorm:"column:payment_id; type:varchar(255); not null" json:"payment_id"`
	Reasons   string    `gorm:"column:reasons; type:text" json:"reasons"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (f *FailedDisbursement) GetFailedDisbursementByPaymentID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &f, "payment_id = ? ", f.PaymentID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (f *FailedDisbursement) CreateFailedDisbursement(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &f)
	if err != nil {
		return fmt.Errorf("failed disbursement creation failed: %v", err.Error())
	}
	return nil
}

func (f *FailedDisbursement) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &f)
	return err
}
