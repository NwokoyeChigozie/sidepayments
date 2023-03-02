package external_models

type MonnifyLoginResponse struct {
	RequestSuccessful bool                     `json:"requestSuccessful"`
	ResponseMessage   string                   `json:"responseMessage"`
	ResponseCode      string                   `json:"responseCode"`
	ResponseBody      MonnifyLoginResponseBody `json:"responseBody"`
}

type MonnifyLoginResponseBody struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

type MonnifyMatchBvnDetailsReq struct {
	Bvn         string `json:"bvn"`
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
	MobileNo    string `json:"mobileNo"`
}

type MonnifyMatchBvnDetailsResponse struct {
	RequestSuccessful bool                               `json:"requestSuccessful"`
	ResponseMessage   string                             `json:"responseMessage"`
	ResponseCode      string                             `json:"responseCode"`
	ResponseBody      MonnifyMatchBvnDetailsResponseBody `json:"responseBody"`
}

type MonnifyMatchBvnDetailsResponseBody struct {
	Bvn         string                                 `json:"bvn"`
	Name        MonnifyMatchBvnDetailsResponseBodyName `json:"name"`
	DateOfBirth string                                 `json:"dateOfBirth"`
	MobileNo    string                                 `json:"mobileNo"`
}

type MonnifyMatchBvnDetailsResponseBodyName struct {
	MatchStatus     string `json:"matchStatus"`
	MatchPercentage int    `json:"matchPercentage"`
}

type MonnifyInitPaymentRequest struct {
	Amount             float64 `json:"amount"`
	CustomerName       string  `json:"customerName"`
	CustomerEmail      string  `json:"customerEmail"`
	PaymentReference   string  `json:"paymentReference"`
	PaymentDescription string  `json:"paymentDescription"`
	CurrencyCode       string  `json:"currencyCode"`
	ContractCode       string  `json:"contractCode"`
	RedirectUrl        string  `json:"redirectUrl"`
}
type MonnifyInitPaymentResponse struct {
	RequestSuccessful bool                           `json:"requestSuccessful"`
	ResponseMessage   string                         `json:"responseMessage"`
	ResponseCode      string                         `json:"responseCode"`
	ResponseBody      MonnifyInitPaymentResponseBody `json:"responseBody"`
}
type MonnifyInitPaymentResponseBody struct {
	TransactionReference string   `json:"transactionReference"`
	PaymentReference     string   `json:"paymentReference"`
	MerchantName         string   `json:"merchantName"`
	ApiKey               string   `json:"apiKey"`
	RedirectUrl          string   `json:"redirectUrl"`
	EnabledPaymentMethod []string `json:"enabledPaymentMethod"`
	CheckoutUrl          string   `json:"checkoutUrl"`
}
