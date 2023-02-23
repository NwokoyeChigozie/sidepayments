package models

import "time"

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
