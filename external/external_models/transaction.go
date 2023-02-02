package external_models

type PartyAccessLevel struct {
	CanView    bool `json:"can_view"`
	CanReceive bool `json:"can_receive"`
	MarkAsDone bool `json:"mark_as_done"`
	Approve    bool `json:"approve"`
}

type PartyResponse struct {
	PartyID     int              `json:"party_id"`
	AccountID   int              `json:"account_id"`
	AccountName string           `json:"account_name"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	Role        string           `json:"role"`
	Status      string           `json:"status"`
	AccessLevel PartyAccessLevel `json:"access_level"`
}

type MilestonesResponse struct {
	Index            int                           `json:"index"`
	MilestoneID      string                        `json:"milestone_id"`
	Title            string                        `json:"title"`
	Amount           float64                       `json:"amount"`
	Status           string                        `json:"status"`
	InspectionPeriod string                        `json:"inspection_period"`
	DueDate          string                        `json:"due_date"`
	Recipients       []MilestonesRecipientResponse `json:"recipients"`
}

type MilestonesRecipientResponse struct {
	AccountID   int     `json:"title"`
	AccountName string  `json:"amount"`
	Email       string  `json:"status"`
	PhoneNumber string  `json:"inspection_period"`
	Amount      float64 `json:"due_date"`
}
