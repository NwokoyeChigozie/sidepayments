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

type RaveReserveAccountRequest struct {
	TxRef       string  `json:"tx_ref"`
	Narration   string  `json:"narration"`
	Amount      float64 `json:"amount"`
	Email       string  `json:"email"`
	Frequency   int     `json:"frequency"`
	Firstname   string  `json:"firstname"`
	Lastname    string  `json:"lastname"`
	IsPermanent bool    `json:"is_permanent"`
}
type RaveReserveAccountResponse struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    RaveReserveAccountResponseData `json:"data"`
}
type RaveReserveAccountResponseData struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	FlwRef          string `json:"flw_ref"`
	OrderRef        string `json:"order_ref"`
	AccountNumber   string `json:"account_number"`
	AccountStatus   string `json:"account_status"`
	Frequency       int    `json:"frequency"`
	BankName        string `json:"bank_name"`
	CreatedAt       int    `json:"created_at"`
	ExpiryDate      int    `json:"expiry_date"`
	Note            string `json:"note"`
	Amount          string `json:"amount"`
}
