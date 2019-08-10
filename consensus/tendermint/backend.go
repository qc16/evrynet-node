package tendermint

import (
	"github.com/ethereum/go-ethereum/event"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {
	// EventMux returns the event mux in backend
	EventMux() *event.TypeMux

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)
}
