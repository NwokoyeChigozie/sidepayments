package config

type Databases struct {
	DB_HOST          string
	DB_PORT          string
	DB_CONNECTION    string
	USERNAME         string
	PASSWORD         string
	TIMEZONE         string
	SSLMODE          string
	ADMIN_DB         string
	AUTH_DB          string
	NOTIFICATIONS_DB string
	PAYMENT_DB       string
	REMINDERS_DB     string
	SUBSCRIPTIONS_DB string
	TRANSACTIONS_DB  string
	VERIFICATION_DB  string
	CRON_DB          string
	Migrate          bool
}
