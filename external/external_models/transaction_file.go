package external_models

type TransactionFile struct {
	ID            uint   `json:"id"`
	TransactionID string `json:"transaction_id"`
	AccountID     int    `json:"account_id"`
	FileType      string `json:"file_type"`
	FileUrl       string `json:"file_url"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
