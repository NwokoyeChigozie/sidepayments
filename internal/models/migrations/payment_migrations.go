package migrations

import "github.com/vesicash/payment-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Payment{},
	}
}
