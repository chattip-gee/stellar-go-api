package service

import (
	"encoding/json"
	"log"
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"
	"github.com/gorilla/mux"

	"github.com/stellar/go/build"
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

func getBalances(w http.ResponseWriter, r *http.Request) {
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
	if _, err := horizon.DefaultTestNetClient.LoadAccount(destination); err != nil {
		panic(err)
	}

	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{AddressOrSeed: source},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
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

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(source)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
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
