package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	. "github.com/chattip-gee/stellar-go-api/constant"
	. "github.com/chattip-gee/stellar-go-api/http"
	"github.com/gorilla/mux"
)

func getPort() string {
	var port = os.Getenv(PORT)
	if port == "" {
		port = LOCALHOST_PORT
		fmt.Println(INFO_NO_PORT_IN_HEROKU + port)
	}
	return ":" + port
}

func HandleRequest() {
	r := mux.NewRouter()
	r.Schemes(HTTPS)

	s := r.PathPrefix(API_PREFIX).Subrouter()
	s.HandleFunc(KEYPAIR_PATH, getKeyPair).Methods(GET)
	s.HandleFunc(FIRENDBOT_PATH, getFriendbot).Methods(GET)
	s.HandleFunc(ACCOUNT_BALANCES_PATH, getBalances).Methods(GET)
	s.HandleFunc(TRANSACTION_PAYMENT_PATH, postTransaction).Methods(POST)

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println(INFO_ROUTE, pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println(INFO_PATH_REGEXP, pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println(INFO_QURIES_TEMPLATES, strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println(INFO_QURIES_REGEXPS, strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println(INFO_METHODS, strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	http.ListenAndServe(getPort(), r)
}
