package service

import (
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"
	"github.com/gorilla/mux"
)

func getAccountsPayments(w http.ResponseWriter, r *http.Request) {
	PrintApiPath(r)

	vars := mux.Vars(r)
	accountsPayments := new(AccountsPaymentsItem)
	if err := JSONDecode(HORIZON_RECEIVE_PAYMENTS_URL+vars[ADDR]+PAYMENTS_PART, accountsPayments); err != nil {
		JSONResponse(w, false, err.Error(), StatusBadRequest, nil)
	} else {
		if accountsPayments.Detail == "" {
			JSONResponse(w, true, SUCCESS, StatusOK, &accountsPayments)
		} else {
			JSONResponse(w, false, FAILED, accountsPayments.Status, nil)
		}
	}

}
