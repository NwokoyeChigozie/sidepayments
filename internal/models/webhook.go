package models

import (
	"fmt"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Webhook struct {
	ID              uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Event           string    `gorm:"column:event; type:varchar(255); not null" json:"event"`
	BusinessID      string    `gorm:"column:business_id; type:varchar(255); not null" json:"business_id"`
	WebhookUri      string    `gorm:"column:webhook_uri; type:varchar(255)" json:"webhook_uri"`
	RequestPayload  string    `gorm:"column:request_payload; type:text" json:"request_payload"`
	IsReceived      bool      `gorm:"column:is_received; default: false" json:"is_received"`
	ResponsePayload string    `gorm:"column:response_payload; type:text" json:"response_payload"`
	ResponseCode    string    `gorm:"column:response_code; type:varchar(255)" json:"response_code"`
	Tries           int       `gorm:"column:tries; type:int; not null; default: 0" json:"tries"`
	IsAbandoned     bool      `gorm:"column:is_abandoned; default: false" json:"is_abandoned"`
	CreatedAt       time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	RetryAt         time.Time `gorm:"column:retry_at" json:"retry_at"`
}

func (w *Webhook) CreateWebhook(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("webhook creation failed: %v", err.Error())
	}
	return nil
}

func (w *Webhook) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &w)
	return err
}
