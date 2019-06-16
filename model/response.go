package model

import "github.com/stellar/go/clients/horizon"

type Any struct{}

type Response struct {
	Success    bool
	Message    string
	StatusCode int
	Data       Any
}

type BalanceResponse struct {
	Success    bool
	Message    string
	StatusCode int
	Data       *[]horizon.Balance
}
