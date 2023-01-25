package external_models

type UsersCredential struct {
	ID                 uint   `json:"id"`
	AccountID          int    `json:"account_id"`
	Bvn                string `json:"bvn"`
	IdentificationType string `json:"identification_type"`
	IdentificationData string `json:"identification_data"`
	DeletedAt          string `json:"deleted_at"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}
type GetUserCredentialModel struct {
	ID                 uint   `json:"id"`
	AccountID          uint   `json:"account_id"`
	IdentificationType string `json:"identification_type"`
}

type CreateUserCredentialModel struct {
	AccountID          uint   `json:"account_id"`
	Bvn                string `json:"bvn"`
	IdentificationType string `json:"identification_type"`
	IdentificationData string `json:"identification_data"`
}
type UpdateUserCredentialModel struct {
	ID                 uint   `json:"id"`
	AccountID          uint   `json:"account_id"`
	IdentificationType string `json:"identification_type"`
	Bvn                string `json:"bvn"`
	IdentificationData string `json:"identification_data"`
}

type GetUserCredentialResponse struct {
	Status  string          `json:"status"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    UsersCredential `json:"data"`
}
