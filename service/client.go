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
	"github.com/stellar/go/clients/horizon"
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
	code := payment.Code
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
		defaultPayment := build.Payment(
			build.Destination{AddressOrSeed: destination},
			build.NativeAmount{Amount: amount},
		)
		if code != NATIVE {
			defaultPayment = build.Payment(
				build.Destination{AddressOrSeed: destination},
				build.CreditAmount{Code: code, Issuer: source, Amount: amount},
			)
		}

		tx, errTransaction := build.Transaction(
			BuildNetwork,
			build.SourceAccount{AddressOrSeed: source},
			build.AutoSequence{SequenceProvider: HorizonDefaultClient},
			build.MemoText{Value: memo},
			build.BaseFee{Amount: baseFee},
			defaultPayment,
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

func postAddAsset(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var assetInfo AssetForm
	errDecode := decoder.Decode(&assetInfo)

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

	issuerAddress := assetInfo.IssuerAddress
	recipientSeed := assetInfo.RecipientSeed
	code := assetInfo.Code
	limit := assetInfo.Limit

	recipient, errKeyParse := keypair.Parse(recipientSeed)
	if errKeyParse != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[KEY_PARSE - ERROR]: "+errKeyParse.Error())
		errResponse := Response{
			Success:    false,
			Message:    errKeyParse.Error(),
			StatusCode: StatusBadRequest,
		}
		JSONEncode(w, errResponse)

		return
	}

	assetName := build.CreditAsset(code, issuerAddress)

	trustTx, errBuildTransaction := build.Transaction(
		build.SourceAccount{AddressOrSeed: recipient.Address()},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.Trust(assetName.Code, assetName.Issuer, build.Limit(limit)),
	)
	if errBuildTransaction != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[BUILD_TRANSACTION - ERROR]: "+errBuildTransaction.Error())
		errResponse := Response{
			Success:    false,
			Message:    errBuildTransaction.Error(),
			StatusCode: StatusBadRequest,
		}
		JSONEncode(w, errResponse)

		return
	}
	trustTxe, errSign := trustTx.Sign(recipientSeed)
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
	trustTxeB64, errTrustTxeB64 := trustTxe.Base64()
	if errTrustTxeB64 != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[TRUST_TXE_BASE64 - ERROR]: "+errTrustTxeB64.Error())
		errResponse := Response{
			Success:    false,
			Message:    errTrustTxeB64.Error(),
			StatusCode: StatusBadRequest,
		}
		JSONEncode(w, errResponse)

		return
	}
	submitResponse, errSubmit := horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if errSubmit != nil {
		fmt.Printf("%q \n %s \n", "[API URL]: "+html.EscapeString(r.URL.Path), "[SUBMIT_TRANSACTION - ERROR]: "+errSubmit.Error())
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
		Data:       &submitResponse,
	}
	JSONEncode(w, response)
}
