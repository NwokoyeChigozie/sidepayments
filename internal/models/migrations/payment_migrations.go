package migrations

import "github.com/vesicash/payment-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Disbursement{},
		models.FundingAccount{},
		models.PaymentCallback{},
		models.PaymentCardInfo{},
		models.PaymentInfo{},
		models.PaymentLog{},
		models.Payment{},
		models.PendingTransferFunding{},
		models.WalletDebitLog{},
		models.WalletEarningLog{},
		models.Webhook{},
	}
}
