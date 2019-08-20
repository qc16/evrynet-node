package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
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
	backend    tendermint.Backend
	events     *event.TypeMuxSubscription
	timeoutSub *event.TypeMuxSubscription

	db ethdb.Database

	valSet       tendermint.ValidatorSet // validators set
	currentState *roundState

	timeoutProposal *event.TypeMuxSubscription
	// timeoutPrevote = TimeoutPrevote or TimeoutPrecommit depends on current round step
	timeoutPrevote *event.TypeMuxSubscription

	handlerWg *sync.WaitGroup
}
