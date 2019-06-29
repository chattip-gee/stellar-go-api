package service

import (
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"

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
