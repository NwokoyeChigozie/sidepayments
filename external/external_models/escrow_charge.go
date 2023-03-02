package external_models

type GetEscrowChargeRequest struct {
	BusinessID int     `json:"business_id"`
	Amount     float64 `json:"amount"`
}

type GetEscrowChargeResponse struct {
	Status  string                      `json:"status"`
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Data    GetEscrowChargeResponseData `json:"data"`
}
type GetEscrowChargeResponseData struct {
	Amount float64 `json:"amount"`
	Charge float64 `json:"charge"`
}
