package migrations

import "github.com/vesicash/payment-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.DisbursementLog{},
		models.DisbursementRequestLog{},
		models.Disbursement{},
		models.FailedDisbursement{},
		models.FundingAccount{},
		models.PaymentAccount{},
		models.PaymentCallback{},
		models.PaymentCardInfo{},
		models.PaymentInfo{},
		models.PaymentLog{},
		models.Payment{},
		models.PendingTransferFunding{},
		// models.WalletDebitLog{},
		models.WalletEarningLog{},
		models.WebhookLog{},
		models.Webhook{},
	}
}
