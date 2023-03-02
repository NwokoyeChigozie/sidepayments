package external_models

type BusinessCharge struct {
	ID                  uint                    `json:"id"`
	BusinessId          int                     `json:"account_id"`
	Country             string                  `json:"country"`
	Currency            string                  `json:"currency"`
	BusinessCharge      string                  `json:"business_charge"`
	VesicashCharge      string                  `json:"vesicash_charge"`
	ProcessingFee       string                  `json:"processing_fee"`
	CancellationFee     string                  `json:"cancellation_fee"`
	DisbursementCharge  string                  `json:"disbursement_charge"`
	PaymentGateway      string                  `json:"payment_gateway"`
	DisbursementGateway string                  `json:"disbursement_gateway"`
	ChargeMin           *map[string]interface{} `json:"charge_min"`
	ChargeMid           *map[string]interface{} `json:"charge_mid"`
	ChargeMax           *map[string]interface{} `json:"charge_max"`
	ProcessingFeeMode   string                  `json:"processing_fee_mode"`
	DeletedAt           string                  `json:"deleted_at"`
	CreatedAt           string                  `json:"created_at"`
	UpdatedAt           string                  `json:"updated_at"`
}

type InitBusinessChargeModel struct {
	BusinessID uint   `json:"business_id"`
	Currency   string `json:"currency"`
}

type GetBusinessChargeResponse struct {
	Status  string         `json:"status"`
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    BusinessCharge `json:"data"`
}

type GetBusinessChargeModel struct {
	ID         uint   `json:"id"`
	BusinessID uint   `json:"business_id"`
	Country    string `json:"country"`
	Currency   string `json:"currency"`
}
