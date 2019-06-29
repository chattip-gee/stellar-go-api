package service

import (
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	. "github.com/chattip-gee/stellar-go-api/model"
	"github.com/gorilla/mux"
)

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
