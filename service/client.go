package service

import (
	"encoding/json"
	"fmt"
	"html"
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
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[ERROR]: "+err.Error())

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
		Message:    SUCCESS,
		StatusCode: StatusOK,
		Data:       &data,
	}
	JSONEncode(w, response)
}

func getFriendbot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if friendBotResp, err := http.Get(HORIZON_FRIENDBOT_URL + vars[ADDR]); err != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[ERROR]: "+err.Error())
		errResponse := Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: StatusBadRequest,
		}
		JSONEncode(w, errResponse)
	} else {
		response := Response{
			Success:    friendBotResp.StatusCode == StatusOK,
			Message:    friendBotResp.Status,
			StatusCode: friendBotResp.StatusCode,
		}
		JSONEncode(w, response)

		defer friendBotResp.Body.Close()
	}
}

func getBalances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if account, err := HorizonDefaultClient.LoadAccount(vars[ADDR]); err != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[ERROR]: "+err.Error())
		errResponse := Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: StatusForbidden,
		}
		JSONEncode(w, errResponse)
	} else {
		balancesItem := BalanceItem{Balances: &account.Balances}
		response := BalanceResponse{
			Success:    true,
			Message:    SUCCESS,
			StatusCode: StatusOK,
			Data:       &balancesItem,
		}
		JSONEncode(w, response)
	}
}

func postTransaction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payment PaymentForm
	errDecode := decoder.Decode(&payment)

	if errDecode != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[DECODE - ERROR]: "+errDecode.Error())
		errResponse := Response{
			Success:    false,
			Message:    errDecode.Error(),
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

	if _, errAccount := HorizonDefaultClient.LoadAccount(destination); errAccount != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[ACCOUNT - ERROR]: "+errAccount.Error())
		accountError := Response{
			Success:    false,
			Message:    errAccount.Error(),
			StatusCode: StatusForbidden,
		}
		JSONEncode(w, accountError)

	} else {
		tx, errTransaction := build.Transaction(
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

		if errTransaction != nil {
			fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[TRANSACTION - ERROR]: "+errTransaction.Error())
			errResponse := Response{
				Success:    false,
				Message:    errTransaction.Error(),
				StatusCode: StatusBadRequest,
			}
			JSONEncode(w, errResponse)

			return
		}

		txe, errSign := tx.Sign(source)
		if errSign != nil {
			fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[SIGN - ERROR]: "+errSign.Error())
			errResponse := Response{
				Success:    false,
				Message:    errSign.Error(),
				StatusCode: StatusBadRequest,
			}
			JSONEncode(w, errResponse)

			return
		}

		txeB64, errBase64 := txe.Base64()
		if errBase64 != nil {
			fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[BASE64 - ERROR]: "+errBase64.Error())
			errResponse := Response{
				Success:    false,
				Message:    errBase64.Error(),
				StatusCode: StatusBadRequest,
			}
			JSONEncode(w, errResponse)

			return
		}

		resp, errSubmit := HorizonDefaultClient.SubmitTransaction(txeB64)
		if errSubmit != nil {
			fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[SUBMIT - ERROR]: "+errSubmit.Error())
			errResponse := Response{
				Success:    false,
				Message:    errSubmit.Error(),
				StatusCode: StatusBadRequest,
			}
			JSONEncode(w, errResponse)

			return
		}

		response := TransactionResponse{
			Success:    true,
			Message:    SUCCESS,
			StatusCode: StatusOK,
			Data:       &resp,
		}
		JSONEncode(w, response)
	}

}
