package external_models

type ProductTransaction struct {
	ID                   int64   `json:"id"`
	TransactionID        string  `json:"transaction_id"`
	ProductTransactionID string  `json:"product_transaction_id"`
	Title                string  `json:"title"`
	Quantity             int64   `json:"quantity"`
	Photo                string  `json:"photo"`
	Amount               float64 `json:"amount"`
	DeletedAt            string  `json:"deleted_at"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}
