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
		JSONResponse(w, false, errDecode.Error(), StatusBadRequest, nil)
		return
	}

	issuerAddress := assetInfo.IssuerAddress
	recipientSeed := assetInfo.RecipientSeed
	code := assetInfo.Code
	limit := assetInfo.Limit

	recipient, errKeyParse := keypair.Parse(recipientSeed)
	if errKeyParse != nil {
		fmt.Print("[KEY_PARSE - ERROR]\n")
		JSONResponse(w, false, errKeyParse.Error(), StatusBadRequest, nil)
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
		JSONResponse(w, false, errBuildTransaction.Error(), StatusBadRequest, nil)
		return
	}

	trustTxe, errSign := trustTx.Sign(recipientSeed)
	if errSign != nil {
		fmt.Print("[SIGN - ERROR]\n")
		JSONResponse(w, false, errSign.Error(), StatusBadRequest, nil)
		return
	}

	trustTxeB64, errTrustTxeB64 := trustTxe.Base64()
	if errTrustTxeB64 != nil {
		fmt.Print("[TRUST_TXE_BASE64 - ERROR]\n")
		JSONResponse(w, false, errTrustTxeB64.Error(), StatusBadRequest, nil)
		return
	}

	submitResponse, errSubmit := horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if errSubmit != nil {
		fmt.Print("[SUBMIT_TRANSACTION - ERROR]\n")
		JSONResponse(w, false, errSubmit.Error(), StatusBadRequest, nil)
		return
	}

	JSONResponse(w, true, SUCCESS, StatusOK, &submitResponse)

}
