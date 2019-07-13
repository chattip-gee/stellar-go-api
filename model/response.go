package model

type Response struct {
	Success    bool         `json:"success"`
	Message    string       `json:"message"`
	StatusCode int          `json:"statusCode"`
	Data       *interface{} `json:"data,omitempty"`
}
