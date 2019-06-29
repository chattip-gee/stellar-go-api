package service

import (
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"
	"github.com/gorilla/mux"
)

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
