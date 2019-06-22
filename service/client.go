package service

import (
	"encoding/json"
	"log"
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"
	"github.com/gorilla/mux"

	"github.com/stellar/go/build"
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

	data := KeyPairItem{Address: pair.Address(), Seed: pair.Seed()}
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
	friendBotResp, err := http.Get(HORIZON_FRIENDBOT_URL + vars[ADDR])
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

func getBalances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	account, err := HorizonDefaultClient.LoadAccount(vars[ADDR])
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

	balancesItem := BalanceItem{Balances: &account.Balances}

	response := BalanceResponse{
		Success:    true,
		Message:    Success,
		StatusCode: StatusOK,
		Data:       &balancesItem,
	}
	JSONEncode(w, response)
}

func postTransaction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var payment PaymentForm
	err := decoder.Decode(&payment)

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

	source := payment.Source
	destination := payment.Destination
	amount := payment.Amount
	memo := payment.Memo
	baseFee := payment.Basefee

	// Make sure destination account exists
	if _, err := HorizonDefaultClient.LoadAccount(destination); err != nil {
		panic(err)
	}

	tx, err := build.Transaction(
		BuildNetwork,
		build.SourceAccount{AddressOrSeed: source},
		build.AutoSequence{SequenceProvider: HorizonDefaultClient},
		build.MemoText{Value: memo},
		build.BaseFee{Amount: baseFee},
		build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.NativeAmount{Amount: amount},
		),
	)

	if err != nil {
		panic(err)
	}

	txe, err := tx.Sign(source)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	resp, err := HorizonDefaultClient.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	response := TransactionResponse{
		Success:    true,
		Message:    Success,
		StatusCode: StatusOK,
		Data:       &resp,
	}
	JSONEncode(w, response)
}
