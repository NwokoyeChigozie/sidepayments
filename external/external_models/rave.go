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

// {
//     "status": "success",
//     "message": "Transfer amount fetched",
//     "data": {
//         "rate": 0.001391,
//         "source": {
//             "currency": "USD",
//             "amount": 1.39073
//         },
//         "destination": {
//             "currency": "NGN",
//             "amount": 1000
//         }
//     }
// }
