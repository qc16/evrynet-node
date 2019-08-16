package tendermint

import (
	"github.com/ethereum/go-ethereum/event"
)

// Backend provides application specific functions for Tendermint core
type Backend interface {
	// EventMux returns the event mux in backend
	EventMux() *event.TypeMux

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	// Gossip sends a message to all validators (exclude self)
	Gossip(valSet ValidatorSet, payload []byte) error
}
