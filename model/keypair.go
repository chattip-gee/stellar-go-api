package model

type KeyPair struct {
	Address string
	Seed    string
}

type KeyPairResponse struct {
	Success    bool
	Message    string
	StatusCode int
	Data       *KeyPair
}
