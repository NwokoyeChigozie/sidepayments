package migrations

import (
	"github.com/vesicash/payment-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func RunAllMigrations(db postgresql.Databases) {

	// payment migration
	MigrateModels(db.Transaction, AuthMigrationModels())

}

func MigrateModels(db *gorm.DB, models []interface{}) {
	_ = db.AutoMigrate(models...)
}
