package models

import (
	"fmt"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PaymentLog struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	PaymentID string    `gorm:"column:payment_id; type:varchar(255); not null" json:"payment_id"`
	Log       string    `gorm:"column:log; type:text; not null" json:"log"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (p *PaymentLog) CreatePaymentLog(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payment log creation failed: %v", err.Error())
	}
	return nil
}
