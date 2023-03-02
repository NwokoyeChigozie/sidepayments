package external_models

type ResolveAccountRequest struct {
	AccountBank   string `json:"account_bank"`
	AccountNumber string `json:"account_number"`
}

type ResolveAccountSuccessResponse struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    ResolveAccountSuccessResponseData `json:"data"`
}

type ResolveAccountSuccessResponseData struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

type ListBanksResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []BanksResponse `json:"data"`
}
type BanksResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ConvertCurrencyRequest struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

type ConvertCurrencyResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    ConvertCurrencyData `json:"data"`
}
type ConvertCurrencyData struct {
	Rate        float64                                `json:"rate"`
	Source      ConvertCurrencyDataSourceOrDestination `json:"source"`
	Destination ConvertCurrencyDataSourceOrDestination `json:"destination"`
}
type ConvertCurrencyDataSourceOrDestination struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type RaveInitPaymentRequest struct {
	TxRef    string `json:"tx_ref"`
	Customer struct {
		Email string `json:"email"`
	} `json:"customer"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	RedirectUrl string  `json:"redirect_url"`
}
type RaveInitPaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}
