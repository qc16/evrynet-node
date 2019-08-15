package backend

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	tendermintCore "github.com/ethereum/go-ethereum/consensus/tendermint/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
)

const (
	tendermintMsg = 0x11
)

// New creates an Ethereum backend for Tendermint core engine.
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
	broadcaster        consensus.Broadcaster
	address            common.Address

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

// Broadcast implements tendermint.Backend.Broadcast
// It sends message to its validator by calling gossiping, and send message to itself by eventMux
// TODO: change AddressSet to validatorSet
func (sb *backend) Broadcast(addressSet []common.Address, payload []byte) error {
	// send to others
	if err := sb.Gossip(addressSet, payload); err != nil {
		return err
	}
	// send to self
	go func() {
		if err := sb.tendermintEventMux.Post(payload); err != nil {
			fmt.Printf("error in Post event %v", err)
		}
	}()
	return nil
}

// Gossip implements tendermint.Backend.Gossip
// It sends message to its validators only, not itself.
// The validators must be able to connected through Peer.
// TODO: change AddressSet to validatorSet
func (sb *backend) Gossip(addressSet []common.Address, payload []byte) error {
	//TODO: check for known message by lru.ARCCache

	targets := make(map[common.Address]bool)
	for _, val := range addressSet {
		if val != sb.address {
			targets[val] = true
		}
	}

	if sb.broadcaster != nil && len(targets) > 0 {
		ps := sb.broadcaster.FindPeers(targets)
		for _, p := range ps {
			//TODO: check for recent messsages using lru.ARCCache
			go func(p consensus.Peer) {
				if err := p.Send(tendermintMsg, payload); err != nil {
					fmt.Printf("Error sending message to peer %v", err)
				}
			}(p)
		}
	}
	return nil
}
