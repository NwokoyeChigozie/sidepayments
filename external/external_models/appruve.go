package external_models

type AppruveReqModelFirst struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
	Gender       string `json:"gender"`
	Phone_number string `json:"phone_number"`
	DateOfBirth  string `json:"date_of_birth"`
	CountryCode  string `json:"country_code"`
	Endpoint     string `json:"endpoint"`
}
type AppruveReqModelMain struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
	Gender       string `json:"gender"`
	Phone_number string `json:"phone_number"`
	DateOfBirth  string `json:"date_of_birth"`
}

type AppruveReponse struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
	Gender       string `json:"gender"`
	Phone_number string `json:"phone_number"`
	DateOfBirth  string `json:"date_of_birth"`
}
