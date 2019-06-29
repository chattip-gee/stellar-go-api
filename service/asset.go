package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

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
