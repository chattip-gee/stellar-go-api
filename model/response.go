package model

type Any struct{}

type Response struct {
	Success    bool
	Message    string
	StatusCode int
	Data       Any
}
