package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PaymentCardInfo struct {
	ID                uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	PaymentID         string    `gorm:"column:payment_id; type:varchar(255); not null" json:"payment_id"`
	CcExpiryMonth     string    `gorm:"column:cc_expiry_month; type:varchar(255)" json:"cc_expiry_month"`
	CcExpiryYear      string    `gorm:"column:cc_expiry_year; type:varchar(255)" json:"cc_expiry_year"`
	LastFourDigits    string    `gorm:"column:lastFourDigits; type:varchar(255)" json:"last_four_digits"`
	Brand             string    `gorm:"column:brand; type:varchar(255)" json:"brand"`
	IssuingCountry    string    `gorm:"column:issuing_country; type:varchar(255)" json:"issuing_country"`
	CardToken         string    `gorm:"column:card_token; type:text" json:"card_token"`
	CardLifeTimeToken string    `gorm:"column:card_life_time_token; type:text" json:"card_life_time_token"`
	Payload           string    `gorm:"column:payload; type:text" json:"payload"`
	CreatedAt         time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	AccountID         int       `gorm:"column:account_id" json:"account_id"`
}

type CardResponse struct {
	AccountID   int    `json:"account_id"`
	Card        string `json:"card"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
}

type DeleteStoredCardRequest struct {
	AccountID int `json:"account_id"  validate:"required" pgvalidate:"exists=auth$users$account_id"`
}

func (p *PaymentCardInfo) DeleteByAccountID(db *gorm.DB) error {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "account_id = ?", p.AccountID)
	if nilErr != nil {
		return nil
	}
	if err != nil {
		return err
	}

	err = postgresql.DeleteRecordFromDb(db, &p)
	if err != nil {
		return err
	}

	return nil
}

func (p *PaymentCardInfo) GetPaymentCardInfoByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, "account_id = ?", p.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (p *PaymentCardInfo) GetPaymentCardInfoByAccountIDLast4DigitsAndBrand(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &p, `account_id = ? and "lastFourDigits"=? and brand=?`, p.AccountID, p.LastFourDigits, p.Brand)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (p *PaymentCardInfo) GetAllPaymentCardInfosByAccountIDs(db *gorm.DB, accountIds []int) ([]PaymentCardInfo, error) {
	details := []PaymentCardInfo{}
	err := postgresql.SelectAllFromDb(db, "desc", &details, "account_id IN (?)", accountIds)
	if err != nil {
		return details, err
	}
	return details, nil
}
func (p *PaymentCardInfo) CreatePaymentCardInfo(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &p)
	if err != nil {
		return fmt.Errorf("Payment card info creation failed: %v", err.Error())
	}
	return nil
}
