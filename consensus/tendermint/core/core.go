package core

import (
	"sync"
	"time"

	"github.com/Workiva/go-datastructures/queue"
	"go.uber.org/zap"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	evrynetCore "github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/rlp"
)

// New creates an Tendermint consensus core
func New(backend tendermint.Backend, config *tendermint.Config) Engine {
	c := &core{
		handlerWg:      new(sync.WaitGroup),
		backend:        backend,
		timeout:        NewTimeoutTicker(),
		config:         config,
		mu:             &sync.RWMutex{},
		blockFinalize:  new(event.TypeMux),
		futureMessages: queue.NewPriorityQueue(0, true),
	}
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	//backend implement tendermint.Backend
	//this component will send/receive data to other nodes and other components
	backend tendermint.Backend
	//events is the channel to receives 2 types of event:
	//- NewBlockEvent: when there is a new composed block from Tx_pool
	//- MessageEvent: when there is a new message from other validators/ peers
	events *event.TypeMuxSubscription

	// finalCommitted events is the chanel to raise when committed and run to new round
	finalCommitted *event.TypeMuxSubscription

	//BlockFinalizeEvent
	blockFinalize *event.TypeMux
	//handleWg will help core stop gracefully, i.e, core will wait till handlingEvents done before reutrning.
	handlerWg *sync.WaitGroup

	//valSet keep track of the current core's validator set.
	valSet tendermint.ValidatorSet // validators set
	//currentState store the state of current consensus
	//it contain round/ block number as well as how many votes this machine has received.
	currentState *roundState

	//timeout will schedule all timeout requirement and fire the timeout event once it's finished.
	timeout TimeoutTicker
	//config store the config of the chain
	config *tendermint.Config
	//mutex mark critical section of core which should not be accessed parallel
	mu *sync.RWMutex

	//proposeStart mark the time core enter propose. This is purely use for metrics
	proposeStart time.Time

	// futureMessages stores future messages (prevote and precommit) fromo other peers
	// and handle them later when we jump to that block number
	// futureMessages only accepts msgItem
	futureMessages *queue.PriorityQueue

	txPool *evrynetCore.TxPool
}

// Start implements core.Engine.Start
// Note: this function is not thread-safe
func (c *core) Start() error {
	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.getLogger().Infow("starting Tendermint's core...")
	if c.currentState == nil {
		c.currentState = c.getStoredState()
	}
	c.subscribeEvents()

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	if err := c.timeout.Start(); err != nil {
		return err
	}
	c.startRoundZero()
	go c.handleEvents()

	return nil
}

// Stop implements core.Engine.Stop
// Note: this function is not thread-safe
func (c *core) Stop() error {
	c.getLogger().Infow("stopping Tendermint's timeout core...")
	err := c.timeout.Stop()
	c.unsubscribeEvents()
	c.handlerWg.Wait()
	c.clearCoreState()
	c.getLogger().Infow("Tendermint's timeout core stopped")
	return err
}

//FinalizeMsg set address, signature and encode msg to bytes
func (c *core) FinalizeMsg(msg *message) ([]byte, error) {
	msg.Address = c.backend.Address()
	msgPayLoadWithoutSignature, err := msg.PayLoadWithoutSignature()
	if err != nil {
		return nil, err
	}
	signature, err := c.backend.Sign(msgPayLoadWithoutSignature)
	if err != nil {
		return nil, err
	}
	msg.Signature = signature
	return rlp.EncodeToBytes(msg)
}

//SendPropose will Finalize the Proposal in term of signature and
//Gossip it to other nodes
func (c *core) SendPropose(propose *tendermint.Proposal) {
	logger := c.getLogger().With("propose_round", propose.Round,
		"propose_block_number", propose.Block.Number(), "propose_block_hash", propose.Block.Hash())

	msgData, err := rlp.EncodeToBytes(propose)
	if err != nil {
		logger.Errorw("Failed to encode Proposal to bytes", "error", err)
		return
	}
	payload, err := c.FinalizeMsg(&message{
		Code: msgPropose,
		Msg:  msgData,
	})
	if err != nil {
		logger.Errorw("Failed to Finalize Proposal", "error", err)
		return
	}

	if err := c.backend.Broadcast(c.valSet, payload); err != nil {
		c.getLogger().Errorw("Failed to Broadcast proposal", "error", err)
		return
	}
	//TODO: remove this log in production
	logger.Infow("sent proposal")
}

//SetBlockForProposal define a method to allow Injecting a Block for testing purpose
func (c *core) SetBlockForProposal(b *types.Block) {
	c.CurrentState().SetBlock(b)
}

//SetTxPool define a method to allow Injecting a txpool
func (c *core) SetTxPool(txPool *evrynetCore.TxPool) {
	c.txPool = txPool
}

//SendVote send broadcast its vote to the network
//it only accept 2 voteType: msgPrevote and msgcommit
func (c *core) SendVote(voteType uint64, block *types.Block, round int64) {
	logger := c.getLogger().With("send_vote_type", voteType, "send_vote_round", round)
	//This should never happen, but it is a safe guard
	if i, _ := c.valSet.GetByAddress(c.backend.Address()); i == -1 {
		logger.Warnw("this node is not a validator of this round, skipping vote", "address", c.backend.Address())
		return
	}
	if voteType != msgPrevote && voteType != msgPrecommit {
		logger.Errorw("vote type is invalid")
		return
	}
	var (
		blockHash = emptyBlockHash
		seal      []byte
	)
	if block != nil {
		var err error
		commitHash := utils.PrepareCommittedSeal(block.Header().Hash())
		seal, err = c.backend.Sign(commitHash)
		if err != nil {
			logger.Errorw("failed to sign seal", err, "err")
			return
		}
		blockHash = block.Hash()
	}
	vote := &tendermint.Vote{
		BlockHash:   &blockHash,
		Round:       round,
		BlockNumber: c.CurrentState().BlockNumber(),
		Seal:        seal,
	}
	msgData, err := rlp.EncodeToBytes(vote)
	if err != nil {
		logger.Errorw("Failed to encode Vote to bytes", "error", err)
		return
	}
	payload, err := c.FinalizeMsg(&message{
		Code: voteType,
		Msg:  msgData,
	})
	if err != nil {
		logger.Errorw("Failed to Finalize Vote", "error", err)
		return
	}
	if err := c.backend.Broadcast(c.valSet, payload); err != nil {
		logger.Errorw("Failed to Broadcast vote", "error", err)
		return
	}
	logger.Infow("sent vote", "vote_round", vote.Round, "vote_block_number", vote.BlockNumber, "vote_block_hash", vote.BlockHash.Hex())
}

func (c *core) CurrentState() *roundState {
	return c.currentState
}

// clearn core's state when the core stop and node in step propose a block
func (c *core) clearCoreState() {
	var (
		state  = c.CurrentState()
		logger = c.getLogger()
	)

	if c.currentState.step == RoundStepPropose {
		// if the currentState in step propose when stop core
		// we have to clear state of core
		logger.Infow("Core's state is cleaning...")
		state.UpdateRoundStep(0, RoundStepNewHeight)

		if state.commitTime.IsZero() {
			// "Now" makes it easier to sync up dev nodes.
			// We add timeoutCommit to allow transactions
			// to be gathered for the first block.
			// And alternative solution that relies on clocks:
			state.startTime = c.config.Commit(time.Now())
		} else {
			state.startTime = c.config.Commit(state.commitTime)
		}
		state.SetLockedRoundAndBlock(-1, nil)
		state.SetValidRoundAndBlock(-1, nil)
		state.SetProposalReceived(nil)

		state.commitRound = -1
		state.PrevotesReceived = make(map[int64]*messageSet)
		state.PrecommitsReceived = make(map[int64]*messageSet)
		state.PrecommitWaited = false

		c.currentState = state
		logger.Infow("Core's state is cleaned")
	}
}

// getLogger returns a zap logger with state info
func (c *core) getLogger() *zap.SugaredLogger {
	if c.currentState == nil {
		return zap.S()
	}
	return zap.L().With(
		zap.Stringer("block", c.currentState.BlockNumber()),
		zap.Int64("round", c.currentState.Round()),
		zap.Stringer("step", c.currentState.Step())).Sugar()
}
