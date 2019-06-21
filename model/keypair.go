package model

type KeyPairItem struct {
	Address string `json:"address"`
	Seed    string `json:"seed"`
}

type KeyPairResponse struct {
	Success    bool         `json:"success"`
	Message    string       `json:"message"`
	StatusCode int          `json:"statusCode"`
	Data       *KeyPairItem `json:"data"`
}
