package external_models

type User struct {
	ID                    uint   `json:"id"`
	AccountID             uint   `json:"account_id"`
	AccountType           string `json:"account_type"`
	Firstname             string `json:"firstname"`
	Lastname              string `json:"lastname"`
	EmailAddress          string `json:"email_address"`
	PhoneNumber           string `json:"phone_number"`
	Username              string `json:"username"`
	TierType              int    `json:"tier_type"`
	DeletedAt             string `json:"deleted_at"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	BusinessId            int    `json:"business_id"`
	Middlename            string `json:"middlename"`
	HasSeenTour           bool   `json:"has_seen_tour"`
	AuthorizationRequired bool   `json:"authorization_required"`
	Meta                  string `json:"meta"`
	ThePeerReference      string `json:"the_peer_reference"`
}

type GetUserModel struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    GetUserData `json:"data"`
}

type GetUserData struct {
	User User `json:"user"`
}

type GetUserRequestModel struct {
	ID           uint   `json:"id"`
	AccountID    uint   `json:"account_id"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
	Username     string `json:"username"`
}
type SetUserAuthorizationRequiredStatusModel struct {
	AccountID uint `json:"account_id"`
	Status    bool `json:"status"`
}

type SetUserAuthorizationRequiredStatusResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    bool   `json:"data"`
}
