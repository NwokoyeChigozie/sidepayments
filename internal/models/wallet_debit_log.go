package models

import (
	"fmt"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type WalletDebitLog struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	Amount    float64   `gorm:"column:amount; type:decimal(20,2); not null" json:"available"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Currency  string    `gorm:"column:currency; type:varchar(255)" json:"currency"`
}

func (w *WalletDebitLog) CreateWalletDebitLog(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &w)
	if err != nil {
		return fmt.Errorf("wallet debit log creation failed: %v", err.Error())
	}
	return nil
}
