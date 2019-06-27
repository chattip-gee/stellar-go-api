package service

import (
	"encoding/json"
	"fmt"
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
	PrintApiPath(r)

	pair, err := keypair.Random()

	if err != nil {
		JSONError(w, err.Error(), StatusInternalServerError)
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
	PrintApiPath(r)

	vars := mux.Vars(r)
	if friendBotResp, err := http.Get(HORIZON_FRIENDBOT_URL + vars[ADDR]); err != nil {
		JSONError(w, err.Error(), StatusBadRequest)
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
	PrintApiPath(r)

	vars := mux.Vars(r)
	if account, err := HorizonDefaultClient.LoadAccount(vars[ADDR]); err != nil {
		JSONError(w, err.Error(), StatusForbidden)
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
	PrintApiPath(r)

	decoder := json.NewDecoder(r.Body)
	var payment PaymentForm
	errDecode := decoder.Decode(&payment)

	if errDecode != nil {
		fmt.Print("[DECODE - ERROR]\n")
		JSONError(w, errDecode.Error(), StatusBadRequest)
		return
	}

	source := payment.Source
	destination := payment.Destination
	amount := payment.Amount
	code := payment.Code
	memo := payment.Memo
	baseFee := payment.Basefee

	if _, errAccount := HorizonDefaultClient.LoadAccount(destination); errAccount != nil {
		fmt.Print("[ACCOUNT - ERROR]\n")
		JSONError(w, errAccount.Error(), StatusForbidden)
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
			fmt.Print("[TRANSACTION - ERROR]\n")
			JSONError(w, errTransaction.Error(), StatusBadRequest)
			return
		}

		txe, errSign := tx.Sign(source)
		if errSign != nil {
			fmt.Print("[SIGN - ERROR]\n")
			JSONError(w, errSign.Error(), StatusBadRequest)
			return
		}

		txeB64, errBase64 := txe.Base64()
		if errBase64 != nil {
			fmt.Print("[BASE64 - ERROR]\n")
			JSONError(w, errBase64.Error(), StatusBadRequest)
			return
		}

		resp, errSubmit := HorizonDefaultClient.SubmitTransaction(txeB64)
		if errSubmit != nil {
			fmt.Print("[SUBMIT - ERROR]\n")
			JSONError(w, errSubmit.Error(), StatusBadRequest)
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
	PrintApiPath(r)

	decoder := json.NewDecoder(r.Body)
	var assetInfo AssetForm
	errDecode := decoder.Decode(&assetInfo)

	if errDecode != nil {
		fmt.Print("[DECODE - ERROR]\n")
		JSONError(w, errDecode.Error(), StatusBadRequest)
		return
	}

	issuerAddress := assetInfo.IssuerAddress
	recipientSeed := assetInfo.RecipientSeed
	code := assetInfo.Code
	limit := assetInfo.Limit

	recipient, errKeyParse := keypair.Parse(recipientSeed)
	if errKeyParse != nil {
		fmt.Print("[KEY_PARSE - ERROR]\n")
		JSONError(w, errKeyParse.Error(), StatusBadRequest)
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
		fmt.Print("[BUILD_TRANSACTION - ERROR]\n")
		JSONError(w, errBuildTransaction.Error(), StatusBadRequest)
		return
	}

	trustTxe, errSign := trustTx.Sign(recipientSeed)
	if errSign != nil {
		fmt.Print("[SIGN - ERROR]\n")
		JSONError(w, errSign.Error(), StatusBadRequest)
		return
	}

	trustTxeB64, errTrustTxeB64 := trustTxe.Base64()
	if errTrustTxeB64 != nil {
		fmt.Print("[TRUST_TXE_BASE64 - ERROR]\n")
		JSONError(w, errTrustTxeB64.Error(), StatusBadRequest)
		return
	}

	submitResponse, errSubmit := horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if errSubmit != nil {
		fmt.Print("[SUBMIT_TRANSACTION - ERROR]\n")
		JSONError(w, errSubmit.Error(), StatusBadRequest)
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
