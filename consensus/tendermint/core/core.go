package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/core/types"
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

	height int64
	round  int
	block  *types.Block

	valSet tendermint.ValidatorSet // validators set

	lockedRound int          // validator's locked round
	lockedBlock *types.Block // validator's locked block

	validRound int          // last known round with PoLC for non-nil valid block
	validBlock *types.Block // last known block of PoLC above

	timeoutProposal *event.TypeMuxSubscription
	// timeoutPrevote = TimeoutPrevote or TimeoutPrecommit depends on current round step
	timeoutPrevote *event.TypeMuxSubscription

	proposalReceived   *tendermint.Proposal
	prevotesReceived   *messageSet
	precommitsReceived *messageSet

	handlerWg *sync.WaitGroup
}
