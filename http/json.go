package http

import (
	"encoding/json"
	"net/http"

	. "github.com/chattip-gee/stellar-go-api/model"
)

func JSONEncode(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

func JSONDecode(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func JSONError(w http.ResponseWriter, message string, code int) {
	response := Response{
		Success:    false,
		Message:    message,
		StatusCode: code,
	}
	JSONEncode(w, response)
}
