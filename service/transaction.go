package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"

	"github.com/stellar/go/build"
)

func postTransaction(w http.ResponseWriter, r *http.Request) {
	PrintApiPath(r)

	decoder := json.NewDecoder(r.Body)
	var payment PaymentForm
	errDecode := decoder.Decode(&payment)

	if errDecode != nil {
		fmt.Print("[DECODE - ERROR]\n")
		JSONResponse(w, false, errDecode.Error(), StatusBadRequest, nil)
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
		JSONResponse(w, false, errAccount.Error(), StatusForbidden, nil)
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
			JSONResponse(w, false, errTransaction.Error(), StatusBadRequest, nil)
			return
		}

		txe, errSign := tx.Sign(source)
		if errSign != nil {
			fmt.Print("[SIGN - ERROR]\n")
			JSONResponse(w, false, errSign.Error(), StatusBadRequest, nil)
			return
		}

		txeB64, errBase64 := txe.Base64()
		if errBase64 != nil {
			fmt.Print("[BASE64 - ERROR]\n")
			JSONResponse(w, false, errBase64.Error(), StatusBadRequest, nil)
			return
		}

		resp, errSubmit := HorizonDefaultClient.SubmitTransaction(txeB64)
		if errSubmit != nil {
			fmt.Print("[SUBMIT - ERROR]\n")
			JSONResponse(w, false, errSubmit.Error(), StatusBadRequest, nil)
			return
		}

		JSONResponse(w, true, SUCCESS, StatusOK, &resp)
	}

}
