package external_models

type SMSToPhoneNotificationRequest struct {
	AccountId   uint   `json:"account_id"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
}
