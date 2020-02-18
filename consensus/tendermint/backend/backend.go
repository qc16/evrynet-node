package backend

import (
	"crypto/ecdsa"
	"math/big"
	"sync"
	"time"

	queue "github.com/enriquebris/goconcurrentqueue"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/backend/fixed_valset_info"
	tendermintCore "github.com/Evrynetlabs/evrynet-node/consensus/tendermint/core"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/event"
	"github.com/Evrynetlabs/evrynet-node/evrdb"
	"github.com/Evrynetlabs/evrynet-node/log"
)

const (
	fetcherID         = "tendermint"
	maxNumberMessages = 64 * 128 * 6 // 64 node * 128 round * 6 messages per round. These number are made higher than expected for safety.
	maxTrigger        = 1000         // maximum of trigger signal that dequeuing op will store.

	maxBroadcastSleepTime        = time.Minute * 5
	initialBroadcastSleepTime    = time.Millisecond * 100
	broadcastSleepTimeIncreament = time.Millisecond * 100
	inMemoryValset               = 10
)

var (
	//ErrNoBroadcaster is return when trying to access backend.Broadcaster without SetBroadcaster first
	ErrNoBroadcaster = errors.New("no broadcaster is set")
)

//Option return an optional function for backend's initial behaviour
type Option func(b *Backend) error

//WithValsetAddresses return an option to assign backend.valSetInfo to fixed valset info
//it will only do so if the input addresses set is not empty
func WithValsetAddresses(addrs []common.Address) Option {
	return func(b *Backend) error {
		if len(addrs) > 0 {
			b.valSetInfo = fixed_valset_info.NewFixedValidatorSetInfo(addrs)
		}
		return nil
	}
}

// New creates an backend for Istanbul core engine.
// The p2p communication, i.e, broadcaster is set separately by calling backend.SetBroadcaster
func New(config *tendermint.Config, privateKey *ecdsa.PrivateKey, opts ...Option) consensus.Tendermint {
	valSetCache, _ := lru.NewARC(inMemoryValset)
	be := &Backend{
		config:               config,
		tendermintEventMux:   new(event.TypeMux),
		privateKey:           privateKey,
		address:              crypto.PubkeyToAddress(privateKey.PublicKey),
		commitChs:            newCommitChannels(),
		mutex:                &sync.RWMutex{},
		storingMsgs:          queue.NewFIFO(),
		dequeueMsgTriggering: make(chan struct{}, maxTrigger),
		broadcastCh:          make(chan broadcastTask),
		controlChan:          make(chan struct{}),
		stakingContractAddr:  common.HexToAddress(config.SCAddress),
		computedValSetCache:  valSetCache,
	}
	be.core = tendermintCore.New(be, config)

	for _, opt := range opts {
		if err := opt(be); err != nil {
			log.Error("error at initialization of backend", err)
		}
	}

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
	db                 evrdb.Database
	broadcaster        consensus.Broadcaster
	address            common.Address

	//once voting finish, the block will be send for commit here
	//it is a map of blocknumber- channels with mutex
	commitChs *commitChannels

	coreStarted bool
	mutex       *sync.RWMutex
	chain       consensus.FullChainReader
	controlChan chan struct{}

	//storingMsgs is used to store msg to handler when core stopped
	storingMsgs          *queue.FIFO
	dequeueMsgTriggering chan struct{}

	currentBlock func() *types.Block

	broadcastCh chan broadcastTask

	valSetInfo          ValidatorSetInfo
	stakingContractAddr common.Address // stakingContractAddr stores the address of staking smart-contract
	computedValSetCache *lru.ARCCache  // computedValSetCache stores the valset is computed from stateDB
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
func (sb *Backend) Broadcast(valSet tendermint.ValidatorSet, blockNumber *big.Int, payload []byte) error {
	// send to others
	if err := sb.Gossip(valSet, blockNumber, payload); err != nil {
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
	Payload     []byte
	MinPeers    int
	TotalPeers  int
	Targets     map[common.Address]bool
	BlockNumber *big.Int
}

// Gossip implements tendermint.Backend.Gossip
// It sends message to its validators only, not itself.
// The validators must be able to connected through Peer.
// It will return backend.ErrNoBroadcaster if no broadcaster is set for backend
func (sb *Backend) Gossip(valSet tendermint.ValidatorSet, blockNumber *big.Int, payload []byte) error {
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
			Payload:     payload,
			MinPeers:    valSet.MinPeers(),
			Targets:     targets,
			TotalPeers:  len(targets),
			BlockNumber: blockNumber,
		}
		go func() {
			sb.broadcastCh <- task
		}()
	}
	return nil
}

func (sb *Backend) gossipLoop() {
	var (
		task        broadcastTask
		finalEvtSub = sb.EventMux().Subscribe(tendermint.FinalCommittedEvent{})
		stopEvtSub  = sb.EventMux().Subscribe(tendermint.StopCoreEvent{})
	)
	defer func() {
		finalEvtSub.Unsubscribe()
		stopEvtSub.Unsubscribe()
	}()
	for {
		select {
		case task = <-sb.broadcastCh:
		case <-finalEvtSub.Chan():
			continue // we skip finalEvtSub to avoid block the go routine send finalEvt to backend.EventMux()
		case <-stopEvtSub.Chan():
			log.Info("cancel broadcast task because core is stopped")
			return
		}

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
				// sleep and retries until success or core stop or new block event
				select {
				case finalEvt := <-finalEvtSub.Chan():
					if finalEvt.Data.(tendermint.FinalCommittedEvent).BlockNumber.Cmp(task.BlockNumber) >= 0 {
						log.Info("cancel broadcast task because of final event")
						break taskLoop
					}
				case <-stopEvtSub.Chan():
					log.Info("cancel broadcast task because core is stopped")
					return
				case <-time.After(timeSleep):
					// increase timeSleep 100ms after each epoch until timeSleep >= maxBroadcastSleepTime
					if timeSleep < maxBroadcastSleepTime {
						timeSleep += broadcastSleepTimeIncreament
					}
				}
				continue taskLoop
			}
			break taskLoop
		}
	}
}

// Multicast implements tendermint.Backend.Multicast
// Send msgs to peers in a set of address
// return err if not found peer with address or sending failed
func (sb *Backend) Multicast(targets map[common.Address]bool, payload []byte) error {
	if sb.broadcaster == nil {
		return ErrNoBroadcaster
	}
	if len(targets) == 0 {
		return nil
	}
	var (
		failed   = 0
		ps       = sb.broadcaster.FindPeers(targets)
		notFound = len(targets) - len(ps)
	)
	log.Trace("multicast", "targets", len(targets), "found", len(ps))
	for addr, peer := range ps {
		if err := peer.Send(consensus.TendermintMsg, payload); err != nil {
			failed++
			log.Debug("failed to send when multicast", "err", err, "addr", addr)
		}
	}
	if failed != 0 || notFound != 0 {
		return errors.Errorf("failed to multicast: failed to send %d address, not found %d address", failed, notFound)
	}
	return nil
}

// Validators return validator set for a block number
// TODO: revise this function once auth vote is implemented
func (sb *Backend) Validators(blockNumber *big.Int) tendermint.ValidatorSet {
	valSet, err := sb.valSetInfo.GetValSet(sb.chain, blockNumber)
	if err != nil {
		log.Error("failed to get validator set", "error", err, "block", blockNumber.Int64())
	}
	return valSet
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
	valSet, err := sb.valSetInfo.GetValSet(chain, blockNumber)
	if err != nil {
		log.Error("failed to get validator set", "error", err, "block", blockNumber.Int64())
	}
	return valSet
}
