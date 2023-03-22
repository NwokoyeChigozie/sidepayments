package models

type RaveWebhookRequest struct {
	Event    string                      `json:"event"`
	Data     *RaveWebhookRequestData     `json:"data"`
	Transfer *RaveWebhookRequestTransfer `json:"transfer"`
}

type RaveWebhookRequestData struct {
	ID                *int                            `json:"id"`
	TxRef             *string                         `json:"tx_ref"`
	FlwRef            *string                         `json:"flw_ref"`
	DeviceFingerprint *string                         `json:"device_fingerprint"`
	Amount            *float64                        `json:"amount"`
	Currency          *string                         `json:"currency"`
	ChargedAmount     *float64                        `json:"charged_amount"`
	AppFee            *float64                        `json:"app_fee"`
	MerchantFee       *float64                        `json:"merchant_fee"`
	ProcessorResponse *string                         `json:"processor_response"`
	AuthModel         *string                         `json:"auth_model"`
	IP                *string                         `json:"ip"`
	Narration         *string                         `json:"narration"`
	Status            *string                         `json:"status"`
	PaymentType       *string                         `json:"payment_type"`
	CreatedAt         *string                         `json:"created_at"`
	AccountNumber     *string                         `json:"account_number"`
	BankName          *string                         `json:"bank_name"`
	BankCode          *string                         `json:"bank_code"`
	Fullname          *string                         `json:"fullname"`
	DebitCurrency     *string                         `json:"debit_currency"`
	Fee               *int                            `json:"fee"`
	Reference         *string                         `json:"reference"`
	Meta              *int                            `json:"meta"`
	Approver          interface{}                     `json:"approver"`
	CompleteMessage   *string                         `json:"complete_message"`
	RequiresApproval  *int                            `json:"requires_approval"`
	IsApproved        *int                            `json:"is_approved"`
	Customer          *RaveWebhookRequestDataCustomer `json:"customer"`
	Card              *RaveWebhookRequestDataCard     `json:"card"`
}

type RaveWebhookRequestTransfer struct {
	ID                *int                            `json:"id"`
	AccountNumber     *string                         `json:"account_number"`
	BankCode          *string                         `json:"bank_code"`
	Fullname          *string                         `json:"fullname"`
	DateCreated       *string                         `json:"date_created"`
	Currency          *string                         `json:"currency"`
	DebitCurrency     *string                         `json:"debit_currency"`
	Amount            *float64                        `json:"amount"`
	Fee               *float64                        `json:"fee"`
	Status            *string                         `json:"status"`
	Reference         *string                         `json:"reference"`
	Meta              *RaveWebhookRequestTransferMeta `json:"meta"`
	Narration         *string                         `json:"narration"`
	Approver          *string                         `json:"approver"`
	CompleteMessage   *string                         `json:"complete_message"`
	Requires_approval *int                            `json:"requires_approval"`
	IsApproved        *int                            `json:"is_approved"`
	BankName          *string                         `json:"bank_name"`
}

type RaveWebhookRequestTransferMeta struct {
	AccountId  *int    `json:"AccountId"`
	MerchantId *string `json:"merchant_id"`
}

type RaveWebhookRequestDataCustomer struct {
	ID          *int    `json:"id"`
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	CreatedAt   *string `json:"created_at"`
}

type RaveWebhookRequestDataCard struct {
	First6digits *string `json:"first_6digits"`
	Last4digits  *string `json:"last_4digits"`
	Issuer       *string `json:"issuer"`
	Country      *string `json:"country"`
	Type         *string `json:"type"`
	Expiry       *string `json:"expiry"`
}
