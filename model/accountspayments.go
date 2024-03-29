package model

type AccountsPaymentsItem struct {
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Embedded *struct {
		Records []struct {
			ID                    string `json:"id"`
			PagingToken           string `json:"paging_token"`
			TransactionSuccessful bool   `json:"transaction_successful"`
			Type                  string `json:"type"`
			CreatedAt             string `json:"created_at"`
			TransactionHash       string `json:"transaction_hash"`
			StartingBalance       string `json:"starting_balance,omitempty"`
			Funder                string `json:"funder,omitempty"`
			Account               string `json:"account,omitempty"`
			AssetType             string `json:"asset_type,omitempty"`
			AssetCode             string `json:"asset_code,omitempty"`
			AssetIssuer           string `json:"asset_issuer,omitempty"`
			From                  string `json:"from,omitempty,omitempty"`
			To                    string `json:"to,omitempty"`
			Amount                string `json:"amount,omitempty"`
		} `json:"records"`
	} `json:"_embedded,omitempty"`
}
