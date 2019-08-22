package core

import (
	"bytes"
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/event"
)

const (
	msgCommit uint64 = iota
)

// New creates an Tendermint consensus core
func New(backend tendermint.Backend) Engine {
	c := &core{
		handlerWg: new(sync.WaitGroup),
		backend:   backend,
	}
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	//backend implement tendermint.Backend
	//this component will send/receive data to other nodes and other components
	backend tendermint.Backend
	//events is the channel to receives 2 types of event:
	//- RequestEvent: when there is a new composed block from Tx_pool
	//- MessageEvent: when there is a new message from other validators/ peers
	events *event.TypeMuxSubscription
	//handleWg will help core stop gracefully, i.e, core will wait till handlingEvents done before reutrning.
	handlerWg *sync.WaitGroup

	//valSet keep track of the current core's validator set.
	valSet tendermint.ValidatorSet // validators set
	//currentState store the state of current consensus
	//it contain round/ block number as well as how many votes this machine has received.
	currentState *roundState

	//timeoutProposal is the channel to receive proposal timeout.
	//TODO: check if the timeout can be done without relating to the current state of core.
	timeoutProposal *event.TypeMuxSubscription
	// timeoutPrevote = channe  or TimeoutPrecommit depends on current round step
	timeoutPrevote *event.TypeMuxSubscription
}

// Start implements core.Engine.Start
func (c *core) Start() error {
	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	go c.handleEvents()

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	c.unsubscribeEvents()
	c.handlerWg.Wait()
	return nil
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgCommit)})
	return buf.Bytes()
}
