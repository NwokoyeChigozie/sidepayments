package models

import (
	"fmt"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type WebhookLog struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Log       string    `gorm:"column:log; type:text; not null" json:"log"`
	Provider  string    `gorm:"column:provider; type:varchar(255)" json:"provider"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (w *WebhookLog) CreateWebhookLog(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("webhook log creation failed: %v", err.Error())
	}
	return nil
}
