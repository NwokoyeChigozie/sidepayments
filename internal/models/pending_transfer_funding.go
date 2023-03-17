package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PendingTransferFunding struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Reference string    `gorm:"column:reference; type:varchar(255); not null" json:"reference"`
	Status    string    `gorm:"column:status; type:varchar(255); not null" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Type      string    `gorm:"column:type; type:varchar(255); not null" json:"type"`
}

func (f *PendingTransferFunding) CreatePendingTransferFunding(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &f)
	if err != nil {
		return fmt.Errorf("pending transfer funding creation failed: %v", err.Error())
	}
	return nil
}

func (p *PendingTransferFunding) GetPendingTransferFundingByReference(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "reference = ?", p.Reference)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *PendingTransferFunding) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &p)
	if err != nil {
		return fmt.Errorf("pending transfer funding delete failed: %v", err.Error())
	}
	return nil
}
