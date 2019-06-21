package model

type Any struct{}

type Response struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       Any    `json:"data"`
}
