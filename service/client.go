package service

import (
	"log"
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"
	"github.com/gorilla/mux"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

func getKeyPair(w http.ResponseWriter, r *http.Request) {
	pair, err := keypair.Random()

	if err != nil {
		log.Fatal(err)

		errResponse := Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: StatusInternalServerError,
		}
		JSONEncode(w, errResponse)

		return
	}

	data := KeyPair{Address: pair.Address(), Seed: pair.Seed()}
	response := KeyPairResponse{
		Success:    true,
		Message:    Success,
		StatusCode: StatusOK,
		Data:       &data,
	}
	JSONEncode(w, response)
}

func getFriendbot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	friendBotResp, err := http.Get("https://horizon-testnet.stellar.org/friendbot?addr=" + vars["addr"])
	if err != nil {
		log.Fatal(err)
		errResponse := Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: StatusBadRequest,
		}
		JSONEncode(w, errResponse)

		return
	}

	var message = Status{
		Code:   friendBotResp.StatusCode,
		Detail: friendBotResp.Status,
	}
	response := Response{
		Success:    friendBotResp.StatusCode == StatusOK,
		Message:    GetMessage(&message),
		StatusCode: friendBotResp.StatusCode,
	}
	JSONEncode(w, response)

	defer friendBotResp.Body.Close()
}

type BalanceResponse struct {
	Success    bool
	Message    string
	StatusCode int
	Data       *[]horizon.Balance
}

func getAccountDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	account, err := horizon.DefaultTestNetClient.LoadAccount(vars["addr"])
	if err != nil {
		log.Fatal(err)
		errResponse := Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: StatusBadRequest,
		}
		JSONEncode(w, errResponse)

		return
	}

	response := BalanceResponse{
		Success:    true,
		Message:    Success,
		StatusCode: StatusOK,
		Data:       &account.Balances,
	}
	JSONEncode(w, response)
}
