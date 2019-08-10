package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/event"
)

// New creates an Istanbul consensus core
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

	handlerWg *sync.WaitGroup
}
