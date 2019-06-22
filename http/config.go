package http

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

const LOCALHOST_PORT = "8080"

var HorizonDefaultClient = horizon.DefaultTestNetClient
var BuildNetwork = build.TestNetwork
