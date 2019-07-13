package service

import (
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	"github.com/gorilla/mux"
)

func getFriendbot(w http.ResponseWriter, r *http.Request) {
	PrintApiPath(r)

	vars := mux.Vars(r)
	if friendBotResp, err := http.Get(HORIZON_FRIENDBOT_URL + vars[ADDR]); err != nil {
		JSONResponse(w, false, err.Error(), StatusBadRequest, nil)
	} else {
		var message = Status{
			Code:   friendBotResp.StatusCode,
			Detail: friendBotResp.Status,
		}
		JSONResponse(w, friendBotResp.StatusCode == StatusOK, GetMessage(&message), friendBotResp.StatusCode, nil)

		defer friendBotResp.Body.Close()
	}

}
