package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const localhost = "8080"

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = localhost
		fmt.Println("No Port In Heroku " + port)
	}
	return ":" + port
}

func HandleRequest() {
	r := mux.NewRouter()
	r.Schemes("https")

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/keypair", getKeyPair).Methods("GET")
	s.HandleFunc("/friendbot/{addr}", getFriendbot).Methods("GET")
	s.HandleFunc("/account/balances/{addr}", getBalances).Methods("GET")
	s.HandleFunc("/transaction/payment", postTransaction).Methods("POST")

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	http.ListenAndServe(getPort(), r)
}
