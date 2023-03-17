package config

type ServerConfiguration struct {
	Port                      string
	Secret                    string
	AccessTokenExpireDuration int
	RequestPerSecond          float64
	TrustedProxies            []string
	ExemptFromThrottle        []string
}
type App struct {
	Name    string
	Key     string
	SiteUrl string
	Url     string
}

type Microservices struct {
	Admin        string
	Auth         string
	Boilerplate  string
	Cron         string
	Feedback     string
	Internaldocs string
	Notification string
	Payment      string
	Productlink  string
	Referral     string
	Reminders    string
	Roles        string
	Subscription string
	Transactions string
	Upload       string
	Verification string
	Widget       string
}

type OnlinePayment struct {
	Max                float64
	DisbursementCharge float64
	NairaThreshold     float64
}
