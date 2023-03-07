package external_models

type WalletBalance struct {
	ID        uint    `json:"id"`
	AccountID int     `json:"account_id"`
	Available float64 `json:"available"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Currency  string  `json:"currency"`
}
type WalletBalanceResponse struct {
	Status  string        `json:"status"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    WalletBalance `json:"data"`
}

type GetWalletRequest struct {
	AccountID uint   `json:"account_id"`
	Currency  string `json:"currency"`
}
type CreateWalletRequest struct {
	AccountID uint    `json:"account_id"`
	Currency  string  `json:"currency"`
	Available float64 `json:"available"`
}
type UpdateWalletRequest struct {
	ID        uint    `json:"id"`
	Available float64 `json:"available"`
}
