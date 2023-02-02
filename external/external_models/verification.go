package external_models

type VerificationFailedModel struct {
	AccountID uint   `json:"account_id"`
	Type      string `json:"type"`
}

type VerificationSuccessfulModel struct {
	AccountID uint   `json:"account_id"`
	Type      string `json:"type"`
}
