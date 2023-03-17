package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type FundingAccount struct {
	ID                   uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID            int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	FundingAccountID     string    `gorm:"column:funding_account_id; type:varchar(255); not null" json:"funding_account_id"`
	Reference            string    `gorm:"column:reference; type:varchar(255); not null" json:"reference"`
	AccountName          string    `gorm:"column:account_name; type:varchar(255)" json:"account_name"`
	BankName             string    `gorm:"column:bank_name; type:varchar(255)" json:"bank_name"`
	BankCode             string    `gorm:"column:bank_code; type:varchar(255)" json:"bank_code"`
	AccountNumber        string    `gorm:"column:account_number; type:varchar(255); not null" json:"account_number"`
	CreatedAt            time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	LastFundingAmount    string    `gorm:"column:last_funding_amount; type:varchar(255)" json:"last_funding_amount"`
	LastFundingReference string    `gorm:"column:last_funding_reference; type:varchar(255)" json:"last_funding_reference"`
	EscrowWallet         string    `gorm:"column:escrow_wallet; type:varchar(255)" json:"escrow_wallet"`
}

func (f *FundingAccount) CreateFundingAccount(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &f)
	if err != nil {
		return fmt.Errorf("funding accounts creation failed: %v", err.Error())
	}
	return nil
}

func (f *FundingAccount) GetFundingAccountByReference(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &f, "reference = ?", f.Reference)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (f *FundingAccount) GetFundingAccountsByAccountID(db *gorm.DB, orderBy, order string, paginator postgresql.Pagination) ([]FundingAccount, postgresql.PaginationResponse, error) {
	var (
		details = []FundingAccount{}
	)

	totalPages, err := postgresql.SelectAllFromDbOrderByPaginated(db, orderBy, order, paginator, &details, "account_id = ?", f.AccountID)
	if err != nil {
		return details, totalPages, err
	}
	return details, totalPages, nil
}
