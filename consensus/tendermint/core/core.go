package core

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/event"
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
	handlerWg  *sync.WaitGroup

	valSet       tendermint.ValidatorSet // validators set
	currentState *roundState

	timeoutProposal *event.TypeMuxSubscription
	// timeoutPrevote = TimeoutPrevote or TimeoutPrecommit depends on current round step
	timeoutPrevote *event.TypeMuxSubscription
}

func (c *core) IsProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.backend.Address()

	// Add proof of consensus
	msg.CommittedSeal = []byte{}
	// Assign the CommittedSeal if it's a COMMIT message and proposal is not nil
	if msg.Code == msgPrecommit && c.currentState.ProposalReceived != nil {
		seal := PrepareCommittedSeal((*c.currentState.ProposalReceived).Hash())
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			return nil, err
		}
	}

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}

	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(msg *message) {
	payload, err := c.finalizeMessage(msg)
	if err != nil {
		fmt.Errorf("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	if err = c.backend.Broadcast(c.valSet, payload); err != nil {
		fmt.Errorf("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

// checkMessage checks the message state
// return ErrInvalidMessage if the message is invalid
// return ErrFutureMessage if the message view is larger than current view
// return ErrOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode uint64, view *tendermint.View) error {
	if view == nil || view.Height == nil || view.Round == nil {
		return tendermint.ErrInvalidMessage
	}

	if view.Cmp(c.currentState.View()) > 0 {
		return tendermint.ErrFutureMessage
	}

	if view.Cmp(c.currentState.View()) < 0 {
		return tendermint.ErrOldMessage
	}

	return nil
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgPrecommit)})
	return buf.Bytes()
}
