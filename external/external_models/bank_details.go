package external_models

type BankDetail struct {
	ID                  uint   `json:"id"`
	AccountID           int    `json:"account_id"`
	BankID              int    `json:"bank_id"`
	AccountName         string `json:"account_name"`
	AccountNo           string `json:"account_no"`
	MobileMoneyOperator string `json:"mobile_money_operator"`
	SwiftCode           string `json:"swift_code"`
	SortCode            string `json:"sort_code"`
	BankAddress         string `json:"bank_address"`
	BankName            string `json:"bank_name"`
	MobileMoneyNumber   string `json:"mobile_money_number"`
	Country             string `json:"country"`
	Currency            string `json:"currency"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type GetBankDetailModel struct {
	ID        uint `json:"id"`
	AccountID uint `json:"account_id"`
}

type GetBankDetailResponse struct {
	Status  string     `json:"status"`
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    BankDetail `json:"data"`
}
