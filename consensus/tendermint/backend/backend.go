package backend

import (
	"crypto/ecdsa"
	"errors"

	"math/big"
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/log"
)

const (
	tendermintMsg = 0x11
)

var (
	//ErrNoBroadcaster is return when trying to access backend.Broadcaster without SetBroadcaster first
	ErrNoBroadcaster = errors.New("no broadcaster is set")
)

//Option return an optional function for backend's initial behaviour
type Option func(b *backend) error

//WithDB return an option to set backend's db
func WithDB(db ethdb.Database) Option {
	return func(b *backend) error {
		b.db = db
		return nil
	}
}

// New creates an backend for Istanbul core engine.
// The p2p communication, i.e, broadcaster is set separately by calling backend.SetBroadcaster
func New(config *tendermint.Config, privateKey *ecdsa.PrivateKey, opts ...Option) consensus.Tendermint {
	be := &backend{
		config:             config,
		tendermintEventMux: new(event.TypeMux),
		privateKey:         privateKey,
		address:            crypto.PubkeyToAddress(privateKey.PublicKey),
	}
	be.core = tendermintCore.New(be, tendermint.DefaultConfig)
	for _, opt := range opts {
		if err := opt(be); err != nil {
			log.Error("error at initialization of backend",
				err)
		}
	}
	return be
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.broadcaster = broadcaster
}

// ----------------------------------------------------------------------------
type backend struct {
	config             *tendermint.Config
	tendermintEventMux *event.TypeMux
	privateKey         *ecdsa.PrivateKey
	core               tendermintCore.Engine
	db                 ethdb.Database
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

// Address implements tendermint.Backend.Address
func (sb *backend) Address() common.Address {
	return sb.address
}

// Broadcast implements tendermint.Backend.Broadcast
// It sends message to its validator by calling gossiping, and send message to itself by eventMux
func (sb *backend) Broadcast(valSet tendermint.ValidatorSet, payload []byte) error {
	// send to others
	if err := sb.Gossip(valSet, payload); err != nil {
		return err
	}
	// send to self
	go func() {
		if err := sb.EventMux().Post(tendermint.MessageEvent{
			Payload: payload,
		}); err != nil {
			log.Error("failed to post event to self", "error", err)
		}
	}()
	return nil
}

// Gossip implements tendermint.Backend.Gossip
// It sends message to its validators only, not itself.
// The validators must be able to connected through Peer.
// It will return backend.ErrNoBroadcaster if no broadcaster is set for backend
func (sb *backend) Gossip(valSet tendermint.ValidatorSet, payload []byte) error {
	//TODO: check for known message by lru.ARCCache

	targets := make(map[common.Address]bool)

	for _, val := range valSet.List() {
		if val.Address() != sb.address {
			targets[val.Address()] = true
		}
	}
	if sb.broadcaster == nil {
		return ErrNoBroadcaster
	}
	if len(targets) > 0 {
		ps := sb.broadcaster.FindPeers(targets)
		for _, p := range ps {
			//TODO: remove these logs in production
			log.Info("sending msg", "from", sb.address.Hex(), "to", p.Address().Hex())
			//TODO: check for recent messsages using lru.ARCCache
			go func(p consensus.Peer) {
				if err := p.Send(tendermintMsg, payload); err != nil {
					log.Error("failed to send message to peer", "error", err)
				}
			}(p)
		}
	}
	return nil
}

func (sb *backend) Validators(blockNumber *big.Int) tendermint.ValidatorSet {
	//TODO: implement this with snapshot
	return validator.NewSet([]common.Address{
		common.HexToAddress("0xb61F4c3E676cE9f4FbF7f5597A303eEeC3AE531B"),
		common.HexToAddress("0xF837F945733a512A087866462B2a4b82FED11146"),
	}, tendermint.RoundRobin)
}
