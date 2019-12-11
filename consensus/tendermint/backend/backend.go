package backend

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"
	"time"

	queue "github.com/enriquebris/goconcurrentqueue"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/log"
)

const (
	fetcherID         = "tendermint"
	maxNumberMessages = 64 * 128 * 6 // 64 node * 128 round * 6 messages per round. These number are made higher than expected for safety.
	maxTrigger        = 1000         // maximum of trigger signal that dequeuing op will store.

	maxBroadcastSleepTime        = time.Minute * 5
	initialBroadcastSleepTime    = time.Millisecond * 100
	broadcastSleepTimeIncreament = time.Millisecond * 100
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

// New creates an backend for Istanbul core engine.
// The p2p communication, i.e, broadcaster is set separately by calling backend.SetBroadcaster
func New(config *tendermint.Config, privateKey *ecdsa.PrivateKey, opts ...Option) consensus.Tendermint {
	be := &Backend{
		config:               config,
		tendermintEventMux:   new(event.TypeMux),
		privateKey:           privateKey,
		address:              crypto.PubkeyToAddress(privateKey.PublicKey),
		commitChs:            newCommitChannels(),
		mutex:                &sync.RWMutex{},
		storingMsgs:          queue.NewFIFO(),
		proposedValidator:    newProposedValidator(),
		dequeueMsgTriggering: make(chan struct{}, maxTrigger),
		broadcastCh:          make(chan broadcastTask),
		controlChan:          make(chan struct{}),
	}
	be.core = tendermintCore.New(be, config)

	for _, opt := range opts {
		if err := opt(be); err != nil {
			log.Error("error at initialization of backend", err)
		}
	}

	go be.gossipLoop()
	go be.dequeueMsgLoop()
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

	//once voting finish, the block will be send for commit here
	//it is a map of blocknumber- channels with mutex
	commitChs *commitChannels

	coreStarted bool
	mutex       *sync.RWMutex
	chain       consensus.ChainReader
	controlChan chan struct{}

	//storingMsgs is used to store msg to handler when core stopped
	storingMsgs          *queue.FIFO
	dequeueMsgTriggering chan struct{}

	currentBlock func() *types.Block

	proposedValidator *ProposalValidator

	broadcastCh chan broadcastTask
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

// SetTxPool define a method to allow Injecting a txpool
func (sb *Backend) SetTxPool(txpool *core.TxPool) {
	sb.core.SetTxPool(txpool)
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

type broadcastTask struct {
	Payload    []byte
	MinPeers   int
	TotalPeers int
	Targets    map[common.Address]bool
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
		task := broadcastTask{
			Payload:    payload,
			MinPeers:   valSet.F() * 2,
			Targets:    targets,
			TotalPeers: len(targets),
		}
		select {
		case sb.broadcastCh <- task:
		default:
			go func() {
				sb.broadcastCh <- task
			}()
		}
	}
	return nil
}

func (sb *Backend) gossipLoop() {
	for {
		task := <-sb.broadcastCh
		var (
			timeSleep   = initialBroadcastSleepTime
			successSent = 0
			mu          sync.Mutex
		)

	taskLoop:
		for {
			ps := sb.broadcaster.FindPeers(task.Targets)
			log.Info("find peers", "len", len(ps), "min", task.MinPeers, "success", successSent)

			var wg sync.WaitGroup
			for addr, p := range ps {
				wg.Add(1)
				//TODO: check for recent messsages using lru.ARCCache
				go func(p consensus.Peer, addr common.Address) {
					defer wg.Done()
					if err := p.Send(consensus.TendermintMsg, task.Payload); err != nil {
						log.Error("failed to send message to peer", "error", err)
						return
					}
					mu.Lock()
					delete(task.Targets, addr)
					successSent += 1
					mu.Unlock()
				}(p, addr)
			}
			wg.Wait()

			if successSent < task.MinPeers {
				log.Info("failed to sent to peer, sleeping", "min", task.MinPeers, "success", successSent,
					"time_sleep", timeSleep)
				// increase timeSleep 100ms after each epoch until timeSleep >= maxBroadcastSleepTime
				// if receive new task then reset the timer
				<-time.After(timeSleep)
				if timeSleep < maxBroadcastSleepTime {
					timeSleep += broadcastSleepTimeIncreament
				}
				continue taskLoop
			}
			break taskLoop
		}
	}
}

// Validators return validator set for a block number
// TODO: revise this function once auth vote is implemented
func (sb *Backend) Validators(blockNumber *big.Int) tendermint.ValidatorSet {
	return sb.getValSet(sb.chain, blockNumber)
}

// FindExistingPeers check validator peers exist or not by address
func (sb *Backend) FindExistingPeers(valSet tendermint.ValidatorSet) map[common.Address]consensus.Peer {
	targets := make(map[common.Address]bool)
	for _, val := range valSet.List() {
		if val.Address() != sb.Address() {
			targets[val.Address()] = true
		}
	}
	return sb.broadcaster.FindPeers(targets)
}

//Commit implement tendermint.Backend.Commit()
func (sb *Backend) Commit(block *types.Block) {
	sb.commitChs.sendBlock(block)
	// if node is not proposer, EnqueueBlock for downloading
	if block.Coinbase() != sb.address {
		sb.EnqueueBlock(block)
	}
}

func (sb *Backend) Cancel(block *types.Block) {
	sb.commitChs.sendBlock(block)
}

// EnqueueBlock adds a block returned from consensus into fetcher queue
func (sb *Backend) EnqueueBlock(block *types.Block) {
	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}
}

func (sb *Backend) CurrentHeadBlock() *types.Block {
	return sb.currentBlock()
}

// ValidatorsByChainReader returns val-set from snapshot
func (sb *Backend) ValidatorsByChainReader(blockNumber *big.Int, chain consensus.ChainReader) tendermint.ValidatorSet {
	return sb.getValSet(chain, blockNumber)
}
