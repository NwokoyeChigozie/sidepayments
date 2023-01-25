package external_models

type ValidateOnDBReq struct {
	Table string      `json:"table"`
	Type  string      `json:"type"`
	Query string      `json:"query"`
	Value interface{} `json:"value"`
}

type ValidateOnDBReqModel struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    bool   `json:"data"`
}

type ValidateAuthorizationReq struct {
	Type               string `validate:"required" json:"type"`
	AuthorizationToken string `json:"authorization-token"`
	VApp               string `json:"v-app"`
	VPrivateKey        string `json:"v-private-key"`
	VPublicKey         string `json:"v-public-key"`
}

type ValidateAuthorizationModel struct {
	Status  string                         `json:"status"`
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Data    ValidateAuthorizationDataModel `json:"data"`
}
type ValidateAuthorizationDataModel struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}
