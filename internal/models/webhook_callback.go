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

// //////// Monnify ////////////////////////////////////////////////////////////////

type MonnifyWebhookRequest struct {
	EventType string                          `json:"eventType"`
	EventData *MonnifyWebhookRequestEventData `json:"eventData"`
}
type MonnifyWebhookRequestEventData struct {
	Product                       *MonnifyWebhookRequestEventDataProduct                       `json:"product"`
	TransactionReference          *string                                                      `json:"transactionReference"`
	PaymentReference              *string                                                      `json:"paymentReference"`
	PaidOn                        *string                                                      `json:"paidOn"`
	PaymentDescription            *string                                                      `json:"paymentDescription"`
	MetaData                      *MonnifyWebhookRequestEventDataMetaData                      `json:"metaData"`
	PaymentSourceInformation      *[]MonnifyWebhookRequestEventDataPaymentSourceInformation    `json:"paymentSourceInformation"`
	DestinationAccountInformation *MonnifyWebhookRequestEventDataDestinationAccountInformation `json:"destinationAccountInformation"`
	AmountPaid                    *float64                                                     `json:"amountPaid"`
	TotalPayable                  *float64                                                     `json:"totalPayable"`
	CardDetails                   *MonnifyWebhookRequestEventDataCardDetails                   `json:"cardDetails"`
	PaymentMethod                 *string                                                      `json:"paymentMethod"`
	Currency                      *string                                                      `json:"currency"`
	SettlementAmount              *string                                                      `json:"settlementAmount"`
	PaymentStatus                 *string                                                      `json:"paymentStatus"`
	Customer                      *MonnifyWebhookRequestEventDataCustomer                      `json:"customer"`
	Amount                        *float64                                                     `json:"amount"`
	Fee                           *float64                                                     `json:"fee"`
	TransactionDescription        *string                                                      `json:"transactionDescription"`
	DestinationAccountNumber      *string                                                      `json:"destinationAccountNumber"`
	SessionId                     *string                                                      `json:"sessionId"`
	CreatedOn                     *string                                                      `json:"createdOn"`
	DestinationAccountName        *string                                                      `json:"destinationAccountName"`
	Reference                     *string                                                      `json:"reference"`
	DestinationBankCode           *string                                                      `json:"destinationBankCode"`
	CompletedOn                   *string                                                      `json:"completedOn"`
	Narration                     *string                                                      `json:"narration"`
	DestinationBankName           *string                                                      `json:"destinationBankName"`
	Status                        *string                                                      `json:"status"`
	MerchantReason                *string                                                      `json:"merchantReason"`
	RefundStatus                  *string                                                      `json:"refundStatus"`
	CustomerNote                  *string                                                      `json:"customerNote"`
	RefundReference               *string                                                      `json:"refundReference"`
	RefundAmount                  *float64                                                     `json:"refundAmount"`
}

type MonnifyWebhookRequestEventDataProduct struct {
	Reference *string `json:"reference"`
	Type      *string `json:"type"`
}
type MonnifyWebhookRequestEventDataMetaData struct {
	Name *string `json:"name"`
	Age  *string `json:"age"`
}
type MonnifyWebhookRequestEventDataPaymentSourceInformation struct {
	BankCode      *string  `json:"bankCode"`
	AmountPaid    *float64 `json:"amountPaid"`
	AccountName   *string  `json:"accountName"`
	SessionId     *string  `json:"sessionId"`
	AccountNumber *string  `json:"accountNumber"`
}
type MonnifyWebhookRequestEventDataDestinationAccountInformation struct {
	BankCode      *string `json:"bankCode"`
	BankName      *string `json:"bankName"`
	AccountNumber *string `json:"accountNumber"`
}
type MonnifyWebhookRequestEventDataCardDetails struct {
	Last4    *string `json:"last4"`
	ExpMonth *string `json:"expMonth"`
	ExpYear  *string `json:"expYear"`
	Bin      *string `json:"bin"`
	Reusable *bool   `json:"reusable"`
}
type MonnifyWebhookRequestEventDataCustomer struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
