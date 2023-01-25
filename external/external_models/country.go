package external_models

type Country struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country_code"`
	CurrencyCode string `json:"currency_code"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type GetCountryModel struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country_code"`
	CurrencyCode string `json:"currency_code"`
}

type GetCountryResponse struct {
	Status  string  `json:"status"`
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    Country `json:"data"`
}
