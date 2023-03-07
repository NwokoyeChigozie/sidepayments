package external_models

type WalletFundedNotificationRequest struct {
	AccountID uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type WalletDebitNotificationRequest struct {
	AccountID     uint    `json:"account_id"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
}
