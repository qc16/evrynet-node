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
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/eth/transaction"
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
type Option func(b *Backend) error

//WithDB return an option to set backend's db
func WithDB(db ethdb.Database) Option {
	return func(b *Backend) error {
		b.db = db
		return nil
	}
}

//WithTxPoolOpts return an option to set backend's txpool
func WithTxPoolOpts(txPoolOpts *transaction.TxPoolOpts) Option {
	return func(b *Backend) error {
		b.txPool = txPoolOpts
		return nil
	}
}

// New creates an backend for Istanbul core engine.
// The p2p communication, i.e, broadcaster is set separately by calling backend.SetBroadcaster
func New(config *tendermint.Config, privateKey *ecdsa.PrivateKey, opts ...Option) consensus.Tendermint {
	be := &Backend{
		config:             config,
		tendermintEventMux: new(event.TypeMux),
		privateKey:         privateKey,
		address:            crypto.PubkeyToAddress(privateKey.PublicKey),
		commitChs:          make(map[string]chan *types.Block),
	}
	be.core = tendermintCore.New(be, tendermint.DefaultConfig)

	for _, opt := range opts {
		if err := opt(be); err != nil {
			log.Error("error at initialization of backend", err)
		}
	}
	return be
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *Backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.broadcaster = broadcaster
}

// ----------------------------------------------------------------------------
type Backend struct {
	config             *tendermint.Config
	tendermintEventMux *event.TypeMux
	privateKey         *ecdsa.PrivateKey
	core               tendermintCore.Engine
	db                 ethdb.Database
	broadcaster        consensus.Broadcaster
	address            common.Address
	txPool             *transaction.TxPoolOpts

	//once voting finish, the block will be send for commit here
	//it is a map of
	commitChs map[string]chan *types.Block

	coreStarted bool
	coreMu      sync.RWMutex
	chain       consensus.ChainReader

	currentBlock func() *types.Block
}

// EventMux implements tendermint.Backend.EventMux
func (sb *Backend) EventMux() *event.TypeMux {
	return sb.tendermintEventMux
}

// Sign implements tendermint.Backend.Sign
func (sb *Backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, sb.privateKey)
}

// Address implements tendermint.Backend.Address
func (sb *Backend) Address() common.Address {
	return sb.address
}

// Broadcast implements tendermint.Backend.Broadcast
// It sends message to its validator by calling gossiping, and send message to itself by eventMux
func (sb *Backend) Broadcast(valSet tendermint.ValidatorSet, payload []byte) error {
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
func (sb *Backend) Gossip(valSet tendermint.ValidatorSet, payload []byte) error {
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
		log.Info("prepare to send message to peers", "total_peers", len(ps))
		for _, p := range ps {
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

// Validators return validator set for a block number
// TODO: revise this function once auth vote is implemented
func (sb *Backend) Validators(blockNumber *big.Int) tendermint.ValidatorSet {
	var (
		previousBlock uint64
		header        *types.Header
		err           error
		snap          *Snapshot
	)
	// check if blockNumber is zero
	if blockNumber.Cmp(big.NewInt(0)) == 0 {
		previousBlock = 0
	} else {
		previousBlock = uint64(blockNumber.Int64() - 1)
	}
	header = sb.chain.GetHeaderByNumber(previousBlock)
	if header == nil {
		log.Error("cannot get valSet since previousBlock is not available", "block_number", blockNumber)
	}
	snap, err = sb.Snapshot(sb.chain, previousBlock, header.Hash(), nil)
	if err != nil {
		log.Error("cannot load snapshot", "error", err)
	}
	if err == nil {
		return snap.ValSet
	}
	return validator.NewSet(nil, sb.config.ProposerPolicy, int64(0))
}

func (sb *Backend) FindPeers(valSet tendermint.ValidatorSet) bool {
	targets := make(map[common.Address]bool)
	for _, val := range valSet.List() {
		if val.Address() != sb.Address() {
			targets[val.Address()] = true
		}
	}

	rs := sb.broadcaster.FindPeers(targets)
	if len(rs) > valSet.F() {
		return true
	}
	return false
}

//Commit implement tendermint.Backend.Commit()
func (sb *Backend) Commit(block *types.Block) {
	ch, ok := sb.commitChs[block.Number().String()]
	if !ok {
		log.Error("no commit channel available", "block_number", block.Number().String())
		return
	}
	ch <- block
}

func (sb *Backend) CurrentHeadBlock() *types.Block {
	return sb.currentBlock()
}

//TxPool return transaction pool
func (sb *Backend) TxPool() *transaction.TxPoolOpts {
	return sb.txPool
}

//Chain return chain
func (sb *Backend) Chain() consensus.ChainReader {
	return sb.chain
}

//Core return Core
func (sb *Backend) Core() *tendermintCore.Core {
	return sb.core.Core()
}

//
// // Verify implements tendermint.Backend.Verify
// func (sb *backend) Verify(proposal tendermint.Proposal) error {
// 	var (
// 		block   = proposal.Block
// 		txs     = block.Transactions()
// 		txnHash = types.DeriveSha(txs)
// 	)
//
// 	// check block body
// 	if txnHash != block.Header().TxHash {
// 		return errMismatchTxhashes
// 	}
//
// 	// Verify transaction for CoreTxPool
// 	if sb.txPool != nil && sb.txPool.CoreTxPool != nil {
// 		for _, t := range txs {
// 			if err := sb.txPool.CoreTxPool.ValidateTx(t, false); err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	// verify the header of proposed block
// 	err := sb.VerifyHeader(sb.chain, block.Header(), false)
// 	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
// 	if err == nil || err == errEmptyCommittedSeals {
// 		return nil
// 	}
// 	return err
// }
