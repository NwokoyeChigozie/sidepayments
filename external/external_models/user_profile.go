package external_models

type UserProfile struct {
	ID         uint   `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID  int    `gorm:"column:account_id; type:int; unique; not null" json:"account_id"`
	Address    string `gorm:"column:address; type:varchar(250)" json:"address"`
	State      string `gorm:"column:state; type:varchar(250)" json:"state"`
	City       string `gorm:"column:city; type:varchar(250)" json:"city"`
	Country    string `gorm:"column:country; type:varchar(250)" json:"country"`
	Dob        string `gorm:"column:dob; type:varchar(250)" json:"dob"`
	Currency   string `gorm:"column:currency; type:varchar(250);default:'USD'; not null" json:"currency"`
	IpAddress  string `gorm:"column:ip_address; type:varchar(250)" json:"ip_address"`
	Sex        string `gorm:"column:sex; type:varchar(250)" json:"sex"`
	Profession string `gorm:"column:profession; type:varchar(250)" json:"profession"`
	Age        uint   `gorm:"column:age; type:int" json:"age"`
	Bio        string `gorm:"column:bio; type:text" json:"bio"`
	DeletedAt  string `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt  string `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt  string `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type GetUserProfileModel struct {
	ID        uint `json:"id"`
	AccountID uint `json:"account_id"`
}

type GetUserProfileResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    UserProfile `json:"data"`
}
