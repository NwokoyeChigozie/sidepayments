package external_models

type VerificationFailedModel struct {
	AccountID uint   `json:"account_id"`
	Type      string `json:"type"`
}

type VerificationSuccessfulModel struct {
	AccountID uint   `json:"account_id"`
	Type      string `json:"type"`
}

type CheckVerificationRequest struct {
	AccountID uint   `json:"account_id"`
	Type      string `json:"type"`
}

type CheckVerificationResponse struct {
	Status  string                        `json:"status"`
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    CheckVerificationResponseData `json:"data"`
}

type CheckVerificationResponseData struct {
	Verified        bool            `json:"verified"`
	VerificationDoc VerificationDoc `json:"verification_doc"`
}
