package external_models

type Transaction struct {
	ID               uint    `json:"id"`
	TransactionID    string  `json:"transaction_id"`
	PartiesID        string  `json:"parties_id"`
	MilestoneID      string  `json:"milestone_id"`
	BrokerID         string  `json:"broker_id"`
	Title            string  `json:"title"`
	Type             string  `json:"type"`
	Description      string  `json:"description"`
	Amount           float64 `json:"amount"`
	Status           string  `json:"status"`
	Quantity         int     `json:"quantity"`
	InspectionPeriod string  `json:"inspection_period"`
	DueDate          string  `json:"due_date"`
	ShippingFee      float64 `json:"shipping_fee"`
	GracePeriod      string  `json:"grace_period"`
	Currency         string  `json:"currency"`
	DeletedAt        string  `json:"deleted_at"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	BusinessID       int     `json:"business_id"`
	IsPaylinked      bool    `json:"is_paylinked"`
	Country          string  `json:"country"`
	Source           string  `json:"source"`
	TransUssdCode    int     `json:"trans_ussd_code"`
	Recipients       string  `json:"recipients"`
	DisputeHandler   string  `json:"dispute_handler"`
	AmountPaid       float64 `json:"amount_paid"`
	EscrowCharge     float64 `json:"escrow_charge"`
	EscrowWallet     string  `json:"escrow_wallet"`
}
type PartyAccessLevel struct {
	CanView    bool `json:"can_view"`
	CanReceive bool `json:"can_receive"`
	MarkAsDone bool `json:"mark_as_done"`
	Approve    bool `json:"approve"`
}

type PartyResponse struct {
	PartyID     int              `json:"party_id"`
	AccountID   int              `json:"account_id"`
	AccountName string           `json:"account_name"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	Role        string           `json:"role"`
	Status      string           `json:"status"`
	AccessLevel PartyAccessLevel `json:"access_level"`
}

type MilestonesResponse struct {
	Index            int                           `json:"index"`
	MilestoneID      string                        `json:"milestone_id"`
	Title            string                        `json:"title"`
	Amount           float64                       `json:"amount"`
	Status           string                        `json:"status"`
	InspectionPeriod string                        `json:"inspection_period"`
	DueDate          string                        `json:"due_date"`
	Recipients       []MilestonesRecipientResponse `json:"recipients"`
}

type MilestonesRecipientResponse struct {
	AccountID   int     `json:"title"`
	AccountName string  `json:"amount"`
	Email       string  `json:"status"`
	PhoneNumber string  `json:"inspection_period"`
	Amount      float64 `json:"due_date"`
}
type MileStoneRecipient struct {
	AccountID    int     `json:"account_id"`
	Amount       float64 `json:"amount"`
	EmailAddress string  `json:"email_address"`
	PhoneNumber  string  `json:"phone_number"`
}

type UpdateTransactionAmountPaidRequest struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Action        string  `json:"action"`
}

type UpdateTransactionAmountPaidResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    Transaction `json:"data"`
}

type CreateActivityLogRequest struct {
	TransactionID string `json:"transaction_id"`
	Description   string `json:"description"`
}

type UpdateTransactionStatusRequest struct {
	AccountID     int    `json:"account_id"`
	TransactionID string `json:"transaction_id"`
	MilestoneID   string `json:"milestone_id"`
	Status        string `json:"status"`
}

type OnlyTransactionIDRequiredRequest struct {
	TransactionID string `json:"transaction_id"`
}
type CreateExchangeTransactionRequest struct {
	AccountID     int     `json:"account_id"`
	InitialAmount float64 `json:"initial_amount"`
	FinalAmount   float64 `json:"final_amount"`
	RateID        int     `json:"rate_id"`
	Status        string  `json:"status"`
}

type Rate struct {
	ID           int64   `json:"id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	From_symbol  string  `json:"from_symbol"`
	ToSymbol     string  `json:"to_symbol"`
	Amount       float64 `json:"amount"`
}
type RateResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Rate   `json:"data"`
}
