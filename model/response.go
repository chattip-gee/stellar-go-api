package model

type Response struct {
	Success    bool
	Message    string
	StatusCode int
	Data       *KeyPair
}
