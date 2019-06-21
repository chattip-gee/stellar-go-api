package model

import "github.com/stellar/go/clients/horizon"

type BalanceItem struct {
	Balances *[]horizon.Balance `json:"balances"`
}

type BalanceResponse struct {
	Success    bool         `json:"success"`
	Message    string       `json:"message"`
	StatusCode int          `json:"statusCode"`
	Data       *BalanceItem `json:"data"`
}
