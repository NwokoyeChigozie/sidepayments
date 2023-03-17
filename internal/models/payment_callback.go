package models

import (
	"fmt"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PaymentCallback struct {
	ID             uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Log            string    `gorm:"column:log; type:text" json:"log"`
	CreatedAt      time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	PaymentLogType string    `gorm:"column:payment_log_type; type:varchar(255); default: webhook; comment: wehbook|initiation" json:"payment_log_type"`
	Reference      string    `gorm:"column:reference; type:varchar(255)" json:"reference"`
}

func (p *PaymentCallback) CreatePaymentCallback(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payment callback creation failed: %v", err.Error())
	}
	return nil
}
