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
type PaymentInvoiceNotificationRequest struct {
	Reference                 string  `json:"reference"`
	PaymentID                 string  `json:"payment_id"`
	TransactionType           string  `json:"transaction_type"`
	TransactionID             string  `json:"transaction_id"`
	Buyer                     int     `json:"buyer"`
	Seller                    int     `json:"seller"`
	InspectionPeriodFormatted string  `json:"inspection_period_formatted"`
	ExpectedDelivery          string  `json:"expected_delivery"`
	Title                     string  `json:"title"`
	Currency                  string  `json:"currency"`
	Amount                    float64 `json:"amount"`
	EscrowCharge              float64 `json:"escrow_charge"`
	BrokerCharge              float64 `json:"broker_charge"`
}
