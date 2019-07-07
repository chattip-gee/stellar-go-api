package http

import (
	"fmt"
	"html"
	"net/http"
)

const PAYMENTS_PART = "/payments"
const FRIENDBOT_PART = "/friendbot"

//Horizon
const HORIZON_HOST_URL = "https://horizon-testnet.stellar.org"
const HORIZON_FRIENDBOT_URL = HORIZON_HOST_URL + FRIENDBOT_PART + "?addr="
const HORIZON_RECEIVE_PAYMENTS_URL = HORIZON_HOST_URL + "/accounts/"

//Prefix
const API_PREFIX = "/api"

//API URL
const KEYPAIR_PATH = "/keypair"
const FIRENDBOT_PATH = "/friendbot/{addr}"
const ACCOUNT_BALANCES_PATH = "/account/balances/{addr}"
const ACCOUNT_ADD_ASSET_PATH = "/account/addasset"
const TRANSACTION_PAYMENT_PATH = "/transaction/payment"
const ACCOUNTS_PAYMENTS_PATH = "/accounts/{addr}" + PAYMENTS_PART

func PrintApiPath(r *http.Request) {
	fmt.Printf("%q \n", "[API URL]: "+html.EscapeString(r.URL.Path))
}
