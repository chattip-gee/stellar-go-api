package model

type PaymentForm struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Amount      string `json:"amount"`
	Code        string `json:"code"`
	Memo        string `json:"memo"`
	Basefee     uint64 `json:"baseFee"`
}
