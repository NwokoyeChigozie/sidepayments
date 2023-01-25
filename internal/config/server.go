package config

type ServerConfiguration struct {
	Port                      string
	Secret                    string
	AccessTokenExpireDuration int
}
type App struct {
	Name string
	Key  string
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
