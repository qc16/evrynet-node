package backend

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
)

// ----------------------------------------------------------------------------

type backend struct {
	istanbulEventMux *event.TypeMux
	privateKey       *ecdsa.PrivateKey

	// the channels for tendermint engine notifications
	commitCh          chan *types.Block
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
}

// EventMux implements istanbul.Backend.EventMux
func (sb *backend) EventMux() *event.TypeMux {
	return sb.istanbulEventMux
}

// Sign implements tendermint.Backend.Sign
func (sb *backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, sb.privateKey)
}
