package external_models

type BusinessProfile struct {
	ID                                uint                   `json:"id"`
	AccountID                         int                    `json:"account_id"`
	BusinessName                      string                 `json:"business_name"`
	BusinessType                      string                 `json:"business_type"`
	LogoUri                           string                 `json:"logo_uri"`
	Website                           string                 `json:"website"`
	Country                           string                 `json:"country"`
	BusinessAddress                   string                 `json:"business_address"`
	PaymentGateway                    string                 `json:"payment_gateway"`
	EscrowChargeOld                   float32                `json:"escrow_charge_old"`
	DisbursementGateway               string                 `json:"disbursement_gateway"`
	AutoTransactionStatusSettings     bool                   `json:"auto_transaction_status_settings"`
	DisbursementSettings              string                 `json:"disbursement_settings"`
	State                             string                 `json:"state"`
	City                              string                 `json:"city"`
	Webhook_uri                       string                 `json:"webhook_uri"`
	Currency                          string                 `json:"currency"`
	IsRegistered                      bool                   `json:"is_registered"`
	DefaultDeliveryPeriod             string                 `json:"default_delivery_period"`
	BusinessIgnoredNotifications      string                 `json:"business_ignored_notifications"`
	BusinessCancellationFee           string                 `json:"business_cancellation_fee"`
	BusinessProcessingFee             string                 `json:"business_processing_fee"`
	AutoAggregateTransactionsSettings bool                   `json:"auto_aggregate_transactions_settings"`
	DefaultChargeBearer               string                 `json:"default_charge_bearer"`
	IsVerificationWaved               bool                   `json:"is_verification_waved"`
	BusinessGivenNotifications        string                 `json:"business_given_notifications"`
	Units                             float32                `json:"units"`
	RedirectUrl                       string                 `json:"redirect_url"`
	BusinessDisabledNotifications     bool                   `json:"business_disabled_notifications"`
	IsBankTransferFeeWaved            bool                   `json:"is_bank_transfer_fee_waved"`
	EscrowCharge                      map[string]interface{} `json:"escrow_charge"`
	Bio                               string                 `json:"bio"`
	BusinessEmail                     string                 `json:"business_email"`
	DeletedAt                         string                 `json:"deleted_at"`
	CreatedAt                         string                 `json:"created_at"`
	UpdatedAt                         string                 `json:"updated_at"`
}

type GetBusinessProfileModel struct {
	ID        uint `json:"id"`
	AccountID uint `json:"account_id"`
}

type GetBusinessProfileResponse struct {
	Status  string          `json:"status"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    BusinessProfile `json:"data"`
}
