package backend

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/consensus"
	tendermintCore "github.com/ethereum/go-ethereum/consensus/tendermint/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
)

// New creates an Ethereum backend for Istanbul core engine.
func New(privateKey *ecdsa.PrivateKey) consensus.Tendermint {
	backend := &backend{
		tendermintEventMux: new(event.TypeMux),
		privateKey:         privateKey,
	}
	backend.core = tendermintCore.New(backend)
	return backend
}

// ----------------------------------------------------------------------------
type backend struct {
	tendermintEventMux *event.TypeMux
	privateKey         *ecdsa.PrivateKey
	core               tendermintCore.Engine

	coreStarted bool
	coreMu      sync.RWMutex
}

// EventMux implements tendermint.Backend.EventMux
func (sb *backend) EventMux() *event.TypeMux {
	return sb.tendermintEventMux
}

// Sign implements tendermint.Backend.Sign
func (sb *backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, sb.privateKey)
}
