package external_models

type UploadFileRequest struct {
	PlaceHolderName string `json:"place_holder_name"`
	File            []byte
}

type UploadFileResponse struct {
	Status  string                   `json:"status"`
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []UploadFileResponseData `json:"data"`
}

type UploadFileResponseData struct {
	OriginalName string `json:"original_name"`
	FileUrl      string `json:"file_url"`
}
