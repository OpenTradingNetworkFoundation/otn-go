package otn

import (
	"os"

	"github.com/opentradingnetworkfoundation/otn-go/api"
)

// Microservice interface represents a basic contract for every
// service that wants to perform in OTN infrastructure
type Microservice interface {
	Start(api api.BitsharesAPI)
	Stop()
	SignalHandler(s os.Signal)
}
