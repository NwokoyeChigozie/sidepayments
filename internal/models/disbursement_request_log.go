package models

import (
	"fmt"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type DisbursementRequestLog struct {
	ID             uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	DisbursementID string    `gorm:"column:disbursement_id; type:varchar(255); not null" json:"disbursement_id"`
	Log            string    `gorm:"column:log; type:text" json:"log"`
	CreatedAt      time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (d *DisbursementRequestLog) CreateDisbursementRequestLog(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &d)
	if err != nil {
		return fmt.Errorf("disbursement request log creation failed: %v", err.Error())
	}
	return nil
}
