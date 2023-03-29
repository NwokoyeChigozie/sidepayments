package config

import "encoding/json"

type Configuration struct {
	Server         ServerConfiguration
	Databases      Databases
	TestDatabases  Databases
	Microservices  Microservices
	App            App
	Monnify        Monnify
	Appruve        Appruve
	Rave           Rave
	IPStack        IPStack
	ONLINE_PAYMENT OnlinePayment
	Slack          Slack
}

type BaseConfig struct {
	SERVER_PORT                      string  `mapstructure:"SERVER_PORT"`
	SERVER_SECRET                    string  `mapstructure:"SERVER_SECRET"`
	SERVER_ACCESSTOKENEXPIREDURATION int     `mapstructure:"SERVER_ACCESSTOKENEXPIREDURATION"`
	REQUEST_PER_SECOND               float64 `mapstructure:"REQUEST_PER_SECOND"`
	TRUSTED_PROXIES                  string  `mapstructure:"TRUSTED_PROXIES"`
	EXEMPT_FROM_THROTTLE             string  `mapstructure:"EXEMPT_FROM_THROTTLE"`
	METRICS_SERVER_PORT              string  `mapstructure:"METRICS_SERVER_PORT"`

	APP_NAME string `mapstructure:"APP_NAME"`
	APP_KEY  string `mapstructure:"APP_KEY"`
	APP_MODE string `mapstructure:"APP_MODE"`
	SITE_URL string `mapstructure:"SITE_URL"`
	APP_URL  string `mapstructure:"APP_URL"`

	DB_HOST          string `mapstructure:"DB_HOST"`
	DB_PORT          string `mapstructure:"DB_PORT"`
	DB_CONNECTION    string `mapstructure:"DB_CONNECTION"`
	TIMEZONE         string `mapstructure:"TIMEZONE"`
	SSLMODE          string `mapstructure:"SSLMODE"`
	USERNAME         string `mapstructure:"USERNAME"`
	PASSWORD         string `mapstructure:"PASSWORD"`
	ADMIN_DB         string `mapstructure:"ADMIN_DB"`
	AUTH_DB          string `mapstructure:"AUTH_DB"`
	NOTIFICATIONS_DB string `mapstructure:"NOTIFICATIONS_DB"`
	PAYMENT_DB       string `mapstructure:"PAYMENT_DB"`
	REMINDERS_DB     string `mapstructure:"REMINDERS_DB"`
	SUBSCRIPTIONS_DB string `mapstructure:"SUBSCRIPTIONS_DB"`
	TRANSACTIONS_DB  string `mapstructure:"TRANSACTIONS_DB"`
	VERIFICATION_DB  string `mapstructure:"VERIFICATION_DB"`
	CRON_DB          string `mapstructure:"CRON_DB"`
	MIGRATE          bool   `mapstructure:"MIGRATE"`

	TEST_DB_HOST          string `mapstructure:"TEST_DB_HOST"`
	TEST_DB_PORT          string `mapstructure:"TEST_DB_PORT"`
	TEST_DB_CONNECTION    string `mapstructure:"TEST_DB_CONNECTION"`
	TEST_TIMEZONE         string `mapstructure:"TEST_TIMEZONE"`
	TEST_SSLMODE          string `mapstructure:"TEST_SSLMODE"`
	TEST_USERNAME         string `mapstructure:"TEST_USERNAME"`
	TEST_PASSWORD         string `mapstructure:"TEST_PASSWORD"`
	TEST_ADMIN_DB         string `mapstructure:"TEST_ADMIN_DB"`
	TEST_AUTH_DB          string `mapstructure:"TEST_AUTH_DB"`
	TEST_NOTIFICATIONS_DB string `mapstructure:"TEST_NOTIFICATIONS_DB"`
	TEST_PAYMENT_DB       string `mapstructure:"TEST_PAYMENT_DB"`
	TEST_REMINDERS_DB     string `mapstructure:"TEST_REMINDERS_DB"`
	TEST_SUBSCRIPTIONS_DB string `mapstructure:"TEST_SUBSCRIPTIONS_DB"`
	TEST_TRANSACTIONS_DB  string `mapstructure:"TEST_TRANSACTIONS_DB"`
	TEST_VERIFICATION_DB  string `mapstructure:"TEST_VERIFICATION_DB"`
	TEST_CRON_DB          string `mapstructure:"TEST_CRON_DB"`
	TEST_MIGRATE          bool   `mapstructure:"TEST_MIGRATE"`

	MS_ADMIN        string `mapstructure:"MS_ADMIN"`
	MS_AUTH         string `mapstructure:"MS_AUTH"`
	MS_BOILERPLATE  string `mapstructure:"MS_BOILERPLATE"`
	MS_CRON         string `mapstructure:"MS_CRON"`
	MS_FEEDBACK     string `mapstructure:"MS_FEEDBACK"`
	MS_INTERNALDOCS string `mapstructure:"MS_INTERNALDOCS"`
	MS_NOTIFICATION string `mapstructure:"MS_NOTIFICATION"`
	MS_PAYMENT      string `mapstructure:"MS_PAYMENT"`
	MS_PRODUCTLINK  string `mapstructure:"MS_PRODUCTLINK"`
	MS_REFERRAL     string `mapstructure:"MS_REFERRAL"`
	MS_REMINDERS    string `mapstructure:"MS_REMINDERS"`
	MS_ROLES        string `mapstructure:"MS_ROLES"`
	MS_SUBSCRIPTION string `mapstructure:"MS_SUBSCRIPTION"`
	MS_TRANSACTIONS string `mapstructure:"MS_TRANSACTIONS"`
	MS_UPLOAD       string `mapstructure:"MS_UPLOAD"`
	MS_VERIFICATION string `mapstructure:"MS_VERIFICATION"`
	MS_WIDGET       string `mapstructure:"MS_WIDGET"`

	MONNIFY_API_KEY                   string `mapstructure:"MONNIFY_API_KEY"`
	MONNIFY_SECRET                    string `mapstructure:"MONNIFY_SECRET"`
	MONNIFY_ENDPOINT                  string `mapstructure:"MONNIFY_ENDPOINT"`
	MONNIFY_CONTRACT_CODE             string `mapstructure:"MONNIFY_CONTRACT_CODE"`
	MONNIFY_BASE64_KEY                string `mapstructure:"MONNIFY_BASE64_KEY"`
	MONNIFY_API                       string `mapstructure:"MONNIFY_API"`
	MONNIFY_DISBURSEMENT_ENDPOINT     string `mapstructure:"MONNIFY_DISBURSEMENT_ENDPOINT"`
	MONNIFY_DISBURSEMENT_USERNAME     string `mapstructure:"MONNIFY_DISBURSEMENT_USERNAME"`
	MONNIFY_DISBURSEMENT_PASSWORD     string `mapstructure:"MONNIFY_DISBURSEMENT_PASSWORD"`
	MONNIFY_DISBURSEMENT_ACCOUNT      string `mapstructure:"MONNIFY_DISBURSEMENT_ACCOUNT"`
	MONNIFY_DISBURSEMENT_ACCOUNT_NAME string `mapstructure:"MONNIFY_DISBURSEMENT_ACCOUNT_NAME"`

	APPRUVE_TEST_ACCESS_TOKEN string `mapstructure:"APPRUVE_TEST_ACCESS_TOKEN"`
	APPRUVE_BASE_URL          string `mapstructure:"APPRUVE_BASE_URL"`

	RAVE_PUBLIC_KEY          string `mapstructure:"RAVE_PUBLIC_KEY"`
	RAVE_SECRET_KEY          string `mapstructure:"RAVE_SECRET_KEY"`
	RAVE_BASE_URL            string `mapstructure:"RAVE_BASE_URL"`
	RAVE_PAYMENT_URL         string `mapstructure:"RAVE_PAYMENT_URL"`
	RAVE_KEY                 string `mapstructure:"RAVE_KEY"`
	RAVE_WEBHOOK_SECRET      string `mapstructure:"RAVE_WEBHOOK_SECRET"`
	FLUTTERWAVE_MERCHANT_ID  string `mapstructure:"FLUTTERWAVE_MERCHANT_ID"`
	FLUTTERWAVE_ACCOUNT_NAME string `mapstructure:"FLUTTERWAVE_ACCOUNT_NAME"`

	IPSTACK_KEY      string `mapstructure:"IPSTACK_KEY"`
	IPSTACK_BASE_URL string `mapstructure:"IPSTACK_BASE_URL"`

	ONLINE_PAYMENT_MAX  float64 `mapstructure:"ONLINE_PAYMENT_MAX"`
	DISBURSEMENT_CHARGE float64 `mapstructure:"DISBURSEMENT_CHARGE"`
	NAIRA_THRESHOLD     float64 `mapstructure:"NAIRA_THRESHOLD"`

	SLACK_OAUTH_TOKEN             string `mapstructure:"SLACK_OAUTH_TOKEN"`
	SLACK_PAYMENT_CHANNELID       string `mapstructure:"SLACK_PAYMENT_CHANNELID"`
	SLACK_DISBURSEMENTS_CHANNELID string `mapstructure:"SLACK_DISBURSEMENTS_CHANNELID"`
}

func (config *BaseConfig) SetupConfigurationn() *Configuration {
	trustedProxies := []string{}
	exemptFromThrottle := []string{}
	json.Unmarshal([]byte(config.TRUSTED_PROXIES), &trustedProxies)
	json.Unmarshal([]byte(config.EXEMPT_FROM_THROTTLE), &exemptFromThrottle)
	return &Configuration{
		Server: ServerConfiguration{
			Port:                      config.SERVER_PORT,
			Secret:                    config.SERVER_SECRET,
			AccessTokenExpireDuration: config.SERVER_ACCESSTOKENEXPIREDURATION,
			RequestPerSecond:          config.REQUEST_PER_SECOND,
			TrustedProxies:            trustedProxies,
			ExemptFromThrottle:        exemptFromThrottle,
			MetricsPort:               config.METRICS_SERVER_PORT,
		},
		App: App{
			Name:    config.APP_NAME,
			Key:     config.APP_KEY,
			Mode:    config.APP_MODE,
			SiteUrl: config.SITE_URL,
			Url:     config.APP_URL,
		},
		Databases: Databases{
			DB_HOST:          config.DB_HOST,
			DB_PORT:          config.DB_PORT,
			DB_CONNECTION:    config.DB_CONNECTION,
			USERNAME:         config.USERNAME,
			PASSWORD:         config.PASSWORD,
			TIMEZONE:         config.TIMEZONE,
			SSLMODE:          config.SSLMODE,
			ADMIN_DB:         config.ADMIN_DB,
			AUTH_DB:          config.AUTH_DB,
			NOTIFICATIONS_DB: config.NOTIFICATIONS_DB,
			PAYMENT_DB:       config.PAYMENT_DB,
			REMINDERS_DB:     config.REMINDERS_DB,
			SUBSCRIPTIONS_DB: config.SUBSCRIPTIONS_DB,
			TRANSACTIONS_DB:  config.TRANSACTIONS_DB,
			VERIFICATION_DB:  config.VERIFICATION_DB,
			CRON_DB:          config.CRON_DB,
			Migrate:          config.MIGRATE,
		},
		TestDatabases: Databases{
			DB_HOST:          config.TEST_DB_HOST,
			DB_PORT:          config.TEST_DB_PORT,
			DB_CONNECTION:    config.TEST_DB_CONNECTION,
			USERNAME:         config.TEST_USERNAME,
			PASSWORD:         config.TEST_PASSWORD,
			TIMEZONE:         config.TEST_TIMEZONE,
			SSLMODE:          config.TEST_SSLMODE,
			ADMIN_DB:         config.TEST_ADMIN_DB,
			AUTH_DB:          config.TEST_AUTH_DB,
			NOTIFICATIONS_DB: config.TEST_NOTIFICATIONS_DB,
			PAYMENT_DB:       config.TEST_PAYMENT_DB,
			REMINDERS_DB:     config.TEST_REMINDERS_DB,
			SUBSCRIPTIONS_DB: config.TEST_SUBSCRIPTIONS_DB,
			TRANSACTIONS_DB:  config.TEST_TRANSACTIONS_DB,
			VERIFICATION_DB:  config.TEST_VERIFICATION_DB,
			CRON_DB:          config.TEST_CRON_DB,
			Migrate:          config.TEST_MIGRATE,
		},
		Microservices: Microservices{
			Admin:        config.MS_ADMIN,
			Auth:         config.MS_AUTH,
			Boilerplate:  config.MS_BOILERPLATE,
			Cron:         config.MS_CRON,
			Feedback:     config.MS_FEEDBACK,
			Internaldocs: config.MS_INTERNALDOCS,
			Notification: config.MS_NOTIFICATION,
			Payment:      config.MS_PAYMENT,
			Productlink:  config.MS_PRODUCTLINK,
			Referral:     config.MS_REFERRAL,
			Reminders:    config.MS_REMINDERS,
			Roles:        config.MS_ROLES,
			Subscription: config.MS_SUBSCRIPTION,
			Transactions: config.MS_TRANSACTIONS,
			Upload:       config.MS_UPLOAD,
			Verification: config.MS_VERIFICATION,
			Widget:       config.MS_WIDGET,
		},
		Monnify: Monnify{
			MonnifyApiKey:                  config.MONNIFY_API_KEY,
			MonnifySecret:                  config.MONNIFY_SECRET,
			MonnifyEndpoint:                config.MONNIFY_ENDPOINT,
			MonnifyContractCode:            config.MONNIFY_CONTRACT_CODE,
			MonnifyBase64Key:               config.MONNIFY_BASE64_KEY,
			MonnifyApi:                     config.MONNIFY_API,
			MonnifyDisbursementEndpoint:    config.MONNIFY_DISBURSEMENT_ENDPOINT,
			MonnifyDisbursementUsername:    config.MONNIFY_DISBURSEMENT_USERNAME,
			MonnifyDisbursementPassword:    config.MONNIFY_DISBURSEMENT_PASSWORD,
			MonnifyDisbursementAccount:     config.MONNIFY_DISBURSEMENT_ACCOUNT,
			MonnifyDisbursementAccountName: config.MONNIFY_DISBURSEMENT_ACCOUNT_NAME,
		},
		Appruve: Appruve{
			AccessToken: config.APPRUVE_TEST_ACCESS_TOKEN,
			BaseUrl:     config.APPRUVE_BASE_URL,
		},
		Rave: Rave{
			PublicKey:     config.RAVE_PUBLIC_KEY,
			SecretKey:     config.RAVE_SECRET_KEY,
			BaseUrl:       config.RAVE_BASE_URL,
			PaymentUrl:    config.RAVE_PAYMENT_URL,
			Key:           config.RAVE_KEY,
			MerchantId:    config.FLUTTERWAVE_MERCHANT_ID,
			AccountName:   config.FLUTTERWAVE_ACCOUNT_NAME,
			WebhookSecret: config.RAVE_WEBHOOK_SECRET,
		},

		IPStack: IPStack{
			Key:     config.IPSTACK_KEY,
			BaseUrl: config.IPSTACK_BASE_URL,
		},

		ONLINE_PAYMENT: OnlinePayment{
			Max:                config.ONLINE_PAYMENT_MAX,
			DisbursementCharge: config.DISBURSEMENT_CHARGE,
			NairaThreshold:     config.NAIRA_THRESHOLD,
		},
		Slack: Slack{
			OauthToken:            config.SLACK_OAUTH_TOKEN,
			PaymentChannelID:      config.SLACK_PAYMENT_CHANNELID,
			DisbursementChannelID: config.SLACK_DISBURSEMENTS_CHANNELID,
		},
	}
}
