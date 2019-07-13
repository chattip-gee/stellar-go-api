package model

import "github.com/stellar/go/clients/horizon"

type BalanceItem struct {
	Balances *[]horizon.Balance `json:"balances"`
}
