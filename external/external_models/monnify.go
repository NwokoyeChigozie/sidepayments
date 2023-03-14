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
type MonnifyVerifyByReferenceResponse struct {
	RequestSuccessful bool                                 `json:"requestSuccessful"`
	ResponseMessage   string                               `json:"responseMessage"`
	ResponseCode      string                               `json:"responseCode"`
	ResponseBody      MonnifyVerifyByReferenceResponseBody `json:"responseBody"`
}
type MonnifyVerifyByReferenceResponseBody struct {
	CreatedOn            string  `json:"createdOn"`
	Amount               float64 `json:"amount"`
	CurrencyCode         string  `json:"currencyCode"`
	CustomerName         string  `json:"customerName"`
	CustomerEmail        string  `json:"customerEmail"`
	PaymentDescription   string  `json:"paymentDescription"`
	PaymentStatus        string  `json:"paymentStatus"`
	TransactionReference string  `json:"transactionReference"`
	PaymentReference     string  `json:"paymentReference"`
}

type MonnifyReserveAccountRequest struct {
	AccountReference string `json:"accountReference"`
	AccountName      string `json:"accountName"`
	CurrencyCode     string `json:"currencyCode"`
	ContractCode     string `json:"contractCode"`
	CustomerEmail    string `json:"customerEmail"`
}

type MonnifyReserveAccountResponse struct {
	RequestSuccessful bool                              `json:"requestSuccessful"`
	ResponseMessage   string                            `json:"responseMessage"`
	ResponseCode      string                            `json:"responseCode"`
	ResponseBody      MonnifyReserveAccountResponseBody `json:"responseBody"`
}

type MonnifyReserveAccountResponseBody struct {
	ContractCode          string        `json:"contractCode"`
	AccountReference      string        `json:"accountReference"`
	AccountName           string        `json:"accountName"`
	CurrencyCode          string        `json:"currencyCode"`
	CustomerEmail         string        `json:"customerEmail"`
	CustomerName          string        `json:"customerName"`
	AccountNumber         string        `json:"accountNumber"`
	BankName              string        `json:"bankName"`
	BankCode              string        `json:"bankCode"`
	CollectionChannel     string        `json:"collectionChannel"`
	ReservationReference  string        `json:"reservationReference"`
	ReservedAccountType   string        `json:"reservedAccountType"`
	Status                string        `json:"status"`
	CreatedOn             string        `json:"createdOn"`
	IncomeSplitConfig     []interface{} `json:"incomeSplitConfig"`
	RestrictPaymentSource bool          `json:"restrictPaymentSource"`
}

type GetMonnifyReserveAccountTransactionsResponse struct {
	RequestSuccessful bool                                             `json:"requestSuccessful"`
	ResponseMessage   string                                           `json:"responseMessage"`
	ResponseCode      string                                           `json:"responseCode"`
	ResponseBody      GetMonnifyReserveAccountTransactionsResponseBody `json:"responseBody"`
}
type GetMonnifyReserveAccountTransactionsResponseBody struct {
	Content  []GetMonnifyReserveAccountTransactionsResponseBodyContent `json:"content"`
	Pageable struct {
		Sort struct {
			Sorted   bool `json:"sorted"`
			Unsorted bool `json:"unsorted"`
			Empty    bool `json:"empty"`
		} `json:"sort"`
		PageSize   int  `json:"pageSize"`
		PageNumber int  `json:"pageNumber"`
		Offset     int  `json:"offset"`
		Unpaged    bool `json:"unpaged"`
		Paged      bool `json:"paged"`
	} `json:"pageable"`
	TotalElements int  `json:"totalElements"`
	TotalPages    int  `json:"totalPages"`
	Last          bool `json:"last"`
	Sort          struct {
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
		Empty    bool `json:"empty"`
	} `json:"sort"`
	First            bool `json:"first"`
	NumberOfElements int  `json:"numberOfElements"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
	Empty            bool `json:"empty"`
}
type GetMonnifyReserveAccountTransactionsResponseBodyContent struct {
	CustomerDTO struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		MerchantCode string `json:"merchantCode"`
	} `json:"customerDTO"`
	ProviderAmount       float64 `json:"providerAmount"`
	PaymentMethod        string  `json:"paymentMethod"`
	CreatedOn            string  `json:"createdOn"`
	Amount               float64 `json:"amount"`
	Flagged              bool    `json:"flagged"`
	ProviderCode         string  `json:"providerCode"`
	Fee                  float64 `json:"fee"`
	CurrencyCode         string  `json:"currencyCode"`
	CompletedOn          string  `json:"completedOn"`
	PaymentDescription   string  `json:"paymentDescription"`
	PaymentStatus        string  `json:"paymentStatus"`
	TransactionReference string  `json:"transactionReference"`
	PaymentReference     string  `json:"paymentReference"`
	MerchantCode         string  `json:"merchantCode"`
	MerchantName         string  `json:"merchantName"`
	PayableAmount        float64 `json:"payableAmount"`
	AmountPaid           float64 `json:"amountPaid"`
	Completed            bool    `json:"completed"`
}
