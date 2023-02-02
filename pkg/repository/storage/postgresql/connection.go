package postgresql

import (
	"fmt"
	"os"

	"log"

	"github.com/vesicash/payment-ms/internal/config"
	"github.com/vesicash/payment-ms/utility"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

type Databases struct {
	Admin         *gorm.DB
	Auth          *gorm.DB
	Notifications *gorm.DB
	Payment       *gorm.DB
	Reminder      *gorm.DB
	Subscription  *gorm.DB
	Transaction   *gorm.DB
	Verification  *gorm.DB
	Cron          *gorm.DB
}

var DB Databases

// Connection gets connection of mysqlDB database
func Connection() Databases {
	return DB
}

func ConnectToDatabases(logger *utility.Logger, configDatabases config.Databases) Databases {
	dbsCV := configDatabases
	databases := Databases{}
	utility.LogAndPrint(logger, "connecting to databases")
	// databases.Admin = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.ADMIN_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Auth = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.AUTH_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Notifications = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.NOTIFICATIONS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Payment = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.PAYMENT_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Reminder = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.REMINDERS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Subscription = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.SUBSCRIPTIONS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Transaction = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.TRANSACTIONS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Verification = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.VERIFICATION_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	// databases.Cron = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.CRON_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)

	utility.LogAndPrint(logger, "connected to databases")

	utility.LogAndPrint(logger, "connected to db")
	// migrations

	DB = databases
	return DB
}

func connectToDb(host, user, password, dbname, port, sslmode, timezone string, logger *utility.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", host, user, password, dbname, port, sslmode, timezone)

	newLogger := lg.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		lg.Config{
			LogLevel:                  lg.Error, // Log level
			IgnoreRecordNotFoundError: true,     // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("connection to %v db failed with: %v", dbname, err))
		panic(err)

	}

	utility.LogAndPrint(logger, fmt.Sprintf("connected to %v db", dbname))
	return db
}

func ReturnDatabase(name string) *gorm.DB {
	databases := DB
	switch name {
	case "admin":
		return DB.Admin
	case "auth":
		return DB.Auth
	case "notifications":
		return DB.Notifications
	case "payment":
		return DB.Payment
	case "reminder":
		return DB.Reminder
	case "subscription":
		return DB.Subscription
	case "transaction":
		return DB.Transaction
	case "verification":
		return DB.Verification
	case "cron":
		return DB.Cron
	default:
		return databases.Auth
	}
}
