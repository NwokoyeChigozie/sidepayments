package models

type DebitWalletRequest struct {
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency" validate:"required"`
	BusinessID    int     `json:"business_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	EscrowWallet  string  `json:"escrow_wallet" validate:"required,oneof=yes no"`
	MorWallet     string  `json:"mor_wallet" validate:"required,oneof=yes no"`
	TransactionID string  `json:"transaction_id" pgvalidate:"exists=transaction$transactions$transaction_id"`
}

type CreditWalletRequest struct {
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency" validate:"required"`
	BusinessID    int     `json:"business_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	IsRefund      bool    `json:"is_refund"`
	EscrowWallet  string  `json:"escrow_wallet" validate:"required,oneof=yes no"`
	MorWallet     string  `json:"mor_wallet" validate:"required,oneof=yes no"`
	TransactionID string  `json:"transaction_id" pgvalidate:"exists=transaction$transactions$transaction_id"`
}
