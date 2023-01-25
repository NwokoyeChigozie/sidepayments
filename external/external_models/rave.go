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
