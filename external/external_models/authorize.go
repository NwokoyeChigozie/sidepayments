package external_models

type Authorize struct {
	ID           uint   `json:"id"`
	AccountID    int    `json:"account_id"`
	Authorized   bool   `json:"authorized"`
	Token        string `json:"token"`
	IpAddress    string `json:"ip_address"`
	Browser      string `json:"browser"`
	Os           string `json:"os"`
	Location     string `json:"location"`
	Attempt      int    `json:"attempt"`
	AuthorizedAt string `json:"authorized_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}

type GetAuthorizeModel struct {
	ID         uint   `json:"id"`
	AccountID  uint   `json:"account_id"`
	Authorized bool   `json:"authorized"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
}

type CreateAuthorizeModel struct {
	AccountID  uint   `json:"account_id"`
	Authorized bool   `json:"authorized"`
	Token      string `json:"token"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
	Os         string `json:"os"`
	Location   string `json:"location"`
}

type UpdateAuthorizeModel struct {
	ID         uint   `json:"id"`
	AccountID  uint   `json:"account_id"`
	Authorized bool   `json:"authorized"`
	Token      string `json:"token"`
	IpAddress  string `json:"ip_address"`
	Browser    string `json:"browser"`
	Os         string `json:"os"`
	Location   string `json:"location"`
	Attempt    int    `json:"attempt"`
}

type GetAuthorizeResponse struct {
	Status  string    `json:"status"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    Authorize `json:"data"`
}

type AuthorizeNotificationRequest struct {
	AccountID int    `json:"account_id"`
	Token     string `json:"token"`
	Ip        string `json:"ip"`
	Location  string `json:"location"`
	Device    string `json:"device"`
}
