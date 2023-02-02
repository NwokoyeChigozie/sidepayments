package external_models

type TransactionParty struct {
	ID                   uint                    `json:"id"`
	TransactionPartiesID string                  `json:"transaction_parties_id"`
	TransactionID        string                  `json:"transaction_id"`
	AccountID            int                     `json:"account_id"`
	Role                 string                  `json:"role"`
	DeletedAt            string                  `json:"deleted_at"`
	CreatedAt            string                  `json:"created_at"`
	UpdatedAt            string                  `json:"updated_at"`
	RoleCapabilities     *map[string]interface{} `json:"role_capabilities"`
	RoleDescription      string                  `json:"role_description"`
	Status               string                  `json:"status"`
}
