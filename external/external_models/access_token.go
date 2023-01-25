package external_models

type AccessToken struct {
	ID            uint   `json:"id"`
	AccountID     int    `json:"account_id"`
	PublicKey     string `json:"public_key"`
	PrivateKey    string `json:"private_key"`
	IsLive        bool   `json:"is_live"`
	IsTermsAgreed bool   `json:"is_terms_agreed"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type GetAccessTokenModel struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    AccessToken `json:"data"`
}
