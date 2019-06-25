package model

type AssetForm struct {
	IssuerAddress string `json:"issuerAddress"`
	RecipientSeed string `json:"recipientSeed"`
	Code          string `json:"code"`
	Limit         string `json:"limit"`
}
