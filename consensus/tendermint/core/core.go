package core

import (
	"bytes"
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/rlp"
)

const (
	msgCommit uint64 = iota
)

// New creates an Tendermint consensus core
func New(backend tendermint.Backend) Engine {
	c := &core{
		handlerWg: new(sync.WaitGroup),
		backend:   backend,
		timeout:   NewTimeoutTicker(),
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
	// timeoutPrevote = channel  or TimeoutPrecommit depends on current round step
	timeoutPrevote *event.TypeMuxSubscription
	//timeout will schedule all timeout requirement and fire the timeout event once it's finished.
	timeout TimeoutTicker
	//config store the config of the chain
	config tendermint.Config
}

// Start implements core.Engine.Start

func (c *core) Start() error {
	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	if err := c.timeout.Start(); err != nil {
		return err
	}
	go c.handleEvents()
	c.startRoundZero()
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

func (c *core) FinalizeMsg(msg message) ([]byte, error) {
	msg.Address = c.backend.Address()
	msgPayLoadWithoutSignature, err := rlp.EncodeToBytes(message{
		Code:    msg.Code,
		Address: msg.Address,
		Msg:     msg.Msg,
	})
	if err != nil {
		return nil, err
	}

	return c.backend.Sign(msgPayLoadWithoutSignature)
}

//SendPropose will Finalize the Proposal in term of signature and
//Gossip it to other nodes
func (c *core) SendPropose(propose *tendermint.Proposal) {
	msgData, err := rlp.EncodeToBytes(propose)
	if err != nil {
		log.Error("Failed to encode Proposal to bytes", "error", err)
	}
	payload, err := c.FinalizeMsg(message{
		Code: msgPropose,
		Msg:  msgData,
	})
	if err != nil {
		log.Error("Failed to Finalize Message", "error", err)
	}

	if err := c.backend.Gossip(c.valSet, payload); err != nil {
		log.Error("Failed to Gossip proposal", "error", err)
	}
}
