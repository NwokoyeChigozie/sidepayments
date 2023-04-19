package external_models

import "time"

type TransactionByID struct {
	ID                  uint                        `json:"id"`
	TransactionID       string                      `json:"transaction_id"`
	PartiesID           string                      `json:"parties_id"`
	MilestoneID         string                      `json:"milestone_id"`
	BrokerID            string                      `json:"broker_id"`
	Title               string                      `json:"title"`
	Type                string                      `json:"type"`
	Description         string                      `json:"description"`
	Amount              float64                     `json:"amount"`
	Status              string                      `json:"status"`
	Quantity            int                         `json:"quantity"`
	InspectionPeriod    string                      `json:"inspection_period"`
	DueDate             string                      `json:"due_date"`
	ShippingFee         float64                     `json:"shipping_fee"`
	GracePeriod         string                      `json:"grace_period"`
	Currency            string                      `json:"currency"`
	DeletedAt           string                      `json:"deleted_at"`
	CreatedAt           string                      `json:"created_at"`
	UpdatedAt           string                      `json:"updated_at"`
	BusinessID          int                         `json:"business_id"`
	IsPaylinked         bool                        `json:"is_paylinked"`
	Source              string                      `json:"source"`
	TransUssdCode       int                         `json:"trans_ussd_code"`
	Recipients          []MileStoneRecipient        `json:"recipients"`
	DisputeHandler      string                      `json:"dispute_handler"`
	AmountPaid          float64                     `json:"amount_paid"`
	EscrowCharge        float64                     `json:"escrow_charge"`
	EscrowWallet        string                      `json:"escrow_wallet"`
	Products            []ProductTransaction        `json:"products"`
	Parties             map[string]TransactionParty `json:"parties"`
	Members             []PartyResponse             `json:"members"`
	Files               []TransactionFile           `json:"files"`
	TotalAmount         float64                     `json:"total_amount"`
	Milestones          []MilestonesResponse        `json:"milestones"`
	Broker              TransactionBroker           `json:"broker"`
	Activities          []ActivityLog               `json:"activities"`
	Country             Country                     `json:"country"`
	DueDateFormatted    string                      `json:"due_date_formatted"`
	TransactionClosedAt time.Time                   `json:"transaction_closed_at"`
	IsDisputed          bool                        `json:"is_disputed"`
}

type ListTransactionsByIDResponse struct {
	Status  string          `json:"status"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    TransactionByID `json:"data"`
}

type ListTransactionsRequestMid struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Filter     string `json:"filter" validate:"oneof=day week month"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type ListTransactionsRequest struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Filter     string `json:"filter" validate:"oneof=day week month"`
}

type ListTransactionsResponse struct {
	Status     string            `json:"status"`
	Code       int               `json:"code"`
	Message    string            `json:"message"`
	Data       []TransactionByID `json:"data"`
	Pagination interface{}       `json:"pagination"`
}
