package external_models

type TransactionBroker struct {
	ID                  uint   `json:"id"`
	TransactionBrokerID string `json:"transaction_broker_id"`
	TransactionID       string `json:"transaction_id"`
	BrokerCharge        string `json:"broker_charge"`
	BrokerChargeBearer  string `json:"broker_charge_bearer"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	BrokerChargeType    string `json:"broker_charge_type"`
	IsSellerAccepted    bool   `json:"is_seller_accepted"`
	IsBuyerAccepted     bool   `json:"is_buyer_accepted"`
}
