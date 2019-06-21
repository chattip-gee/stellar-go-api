package model

import (
	. "github.com/stellar/go/clients/horizon"
)

type PaymentForm struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Amount      string `json:"amount"`
	Memo        string `json:"memo"`
	Basefee     uint64 `json:"baseFee"`
}

type TransactionResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	StatusCode int                 `json:"statusCode"`
	Data       *TransactionSuccess `json:"data"`
}
