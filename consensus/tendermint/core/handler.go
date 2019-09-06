package core

import (
	"errors"
	"fmt"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/rlp"
)

var (
	ErrInvalidProposalPOLRound  = errors.New("Error invalid proposal POL round")
	ErrInvalidProposalSignature = errors.New("Error invalid proposal signature")
)

// ----------------------------------------------------------------------------

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		tendermint.NewBlockEvent{},
		tendermint.MessageEvent{},
		tendermint.Proposal{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
}

// handleEvents will receive messages as well as timeout and is solely responsible for state change.
func (c *core) handleEvents() {
	// Clear state
	defer func() {
		c.handlerWg.Done()
	}()

	c.handlerWg.Add(1)

	for {
		log.Debug("core's handling is running...")
		select {
		case event, ok := <-c.events.Chan(): //backend sending something...
			if !ok {
				return
			}
			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case tendermint.NewBlockEvent:
				log.Debug("Received New Block event", "event", ev)
				c.currentState.SetBlock(ev.Block)
			case tendermint.MessageEvent:

				log.Debug("Received Message event", "message", ev)
				//TODO: Handle ev.Payload, if got error then call c.backend.Gossip()
				var msg message
				if err := rlp.DecodeBytes(ev.Payload, &msg); err != nil {
					log.Error("failed to decode msg", "error", err)
				} else {
					if err := c.handleMsg(msg); err != nil {
						log.Error("failed decode msg", "error", err)
					}
				}
			default:
				log.Debug("Unknown event ", "event", ev)
			}
		case ti := <-c.timeout.Chan(): //something from timeout...
			c.handleTimeout(ti)
		}
	}
}

func (c *core) verifyProposal(proposal tendermint.Proposal, msg message) error {

	// Verify POLRound, which must be nil or in range [0, proposal.Round).
	if proposal.POLRound != -1 &&
		(proposal.POLRound >= 0) && proposal.POLRound >= proposal.Round {
		return ErrInvalidProposalPOLRound
	}

	// Verify signature
	signer, err := msg.GetAddressFromSignature()
	if err != nil {
		return err
	}

	// signature must come from Proposer of this round
	if c.valSet.GetProposer().Address() != signer {
		return ErrInvalidProposalSignature
	}
	return nil
}

func (c *core) handlePropose(msg message) error {
	var (
		state    = c.currentState
		proposal tendermint.Proposal
	)
	if err := rlp.DecodeBytes(msg.Msg, &proposal); err != nil {
		return err
	}
	// Already have one
	// TODO: possibly catch double proposals
	if state.ProposalReceived() != nil {
		return nil
	}
	// Does not apply, this is not an error but may happen due to network lattency
	if proposal.Block.Number().Cmp(state.BlockNumber()) != 0 || proposal.Round != state.Round() {
		log.Debug("Received proposal with different height/round")
		return nil
	}
	if err := c.verifyProposal(proposal, msg); err != nil {
		return err
	}
	state.SetProposalReceived(&proposal)
	//// TODO: We can check if Proposal is for a different block as this is a sign of misbehavior!
	log.Info("Received proposal", "proposal", proposal)
	return nil
}

func (c *core) handleMsg(msg message) error {
	switch msg.Code {
	case msgPropose:
		return c.handlePropose(msg)
	default:
		return fmt.Errorf("unknown msg code %d", msg.Code)
	}
}

func (c *core) handleTimeout(ti timeoutInfo) {
	log.Debug("Received timeout signal from core.timeout", "timeout", ti.Duration, "block_number", ti.BlockNumber, "round", ti.Round, "step", ti.Step)
	var (
		round       = c.currentState.Round()
		blockNumber = c.currentState.BlockNumber()
		step        = c.currentState.Step()
	)
	// timeouts must be for current height, round, step
	if ti.BlockNumber.Cmp(blockNumber) != 0 || ti.Round < round || (ti.Round == round && ti.Step < step) {
		log.Debug("Ignoring timeout because we're ahead", "block_number", blockNumber, "round", round, "step", step)
		return
	}

	// the timeout will now cause a state transition
	c.currentState.mu.Lock()
	defer c.currentState.mu.Unlock()

	switch ti.Step {
	case RoundStepNewHeight:
		// NewRound event fired from enterNewRound.
		c.enterNewRound(ti.BlockNumber, 0)
	case RoundStepNewRound:
		c.enterPropose(ti.BlockNumber, 0)
	case RoundStepPropose:
		c.enterPrevote(ti.BlockNumber, ti.Round)
	case RoundStepPrevoteWait:
		c.enterPrecommit(ti.BlockNumber, ti.Round)
	case RoundStepPrecommitWait:
		c.enterPrecommit(ti.BlockNumber, ti.Round)
		c.enterNewRound(ti.BlockNumber, ti.Round+1)
	default:
		panic(fmt.Sprintf("Invalid timeout step: %v", ti.Step))
	}

}
