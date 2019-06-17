package model

import "github.com/stellar/go/clients/horizon"

type BalanceResponse struct {
	Success    bool
	Message    string
	StatusCode int
	Data       *[]horizon.Balance
}
