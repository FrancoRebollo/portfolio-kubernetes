package dto

type DefaultResponse struct {
	Message string `json:"message"`
}

type ExternalAPIResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type ResponseDefault struct {
	Message string `json:"message"`
}
