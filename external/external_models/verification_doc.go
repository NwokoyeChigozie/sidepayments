package external_models

type VerificationDoc struct {
	ID        uint   `json:"id"`
	AccountID int    `json:"account_id"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	DeletedAt string `json:"deleted_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Meta      string `json:"meta"`
}
