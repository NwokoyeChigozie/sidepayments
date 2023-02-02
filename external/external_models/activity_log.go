package external_models

type ActivityLog struct {
	ID            uint   `json:"id"`
	TransactionID string `json:"transaction_id"`
	Description   string `json:"description"`
	DeletedAt     string `json:"deleted_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
