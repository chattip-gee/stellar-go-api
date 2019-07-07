package model

type AccountsPaymentsItem struct {
	Embedded struct {
		Records []struct {
			ID          string `json:"id"`
			PagingToken string `json:"paging_token"`
			Type        string `json:"type"`
			CreatedAt   string `json:"created_at"`
			AssetType   string `json:"asset_type"`
			AssetCode   string `json:"asset_code"`
			AssetIssuer string `json:"asset_issuer"`
			From        string `json:"from"`
			To          string `json:"to"`
			Amount      string `json:"amount"`
		} `json:"records"`
	} `json:"_embedded"`
}

type AccountsPaymentsResponse struct {
	Success    bool                  `json:"success"`
	Message    string                `json:"message"`
	StatusCode int                   `json:"statusCode"`
	Data       *AccountsPaymentsItem `json:"data"`
}
