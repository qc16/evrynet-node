package core

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/log"
)

// ----------------------------------------------------------------------------

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		tendermint.NewBlockEvent{},
		tendermint.MessageEvent{},
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
		select {
		case event, ok := <-c.events.Chan(): //backend sending something...
			if !ok {
				return
			}
			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case tendermint.NewBlockEvent:
				c.currentState.SetBlock(ev.Block)
			case tendermint.MessageEvent:
				fmt.Printf("--- Type of event.Data: %+v\n", reflect.TypeOf(ev))
				fmt.Printf("--- Value of event.Data: %+v\n", ev.Payload)
				//TODO: Handle ev.Payload, if got error then call c.backend.Gossip()
			default:
				fmt.Printf("--- Unknow event :%v", ev)
			}
		case ti := <-c.timeout.Chan(): //something from timeout...
			c.handleTimeout(ti)
		}
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
	if ti.BlockNumber.Cmp(blockNumber) != 0 || ti.Round.Cmp(round) < 0 || (ti.Round.Cmp(round) == 0 && ti.Step < step) {
		log.Debug("Ignoring timeout because we're ahead", "block_number", blockNumber, "round", round, "step", step)
		return
	}

	// the timeout will now cause a state transition
	c.currentState.mu.Lock()
	defer c.currentState.mu.Unlock()

	switch ti.Step {
	case RoundStepNewHeight:
		// NewRound event fired from enterNewRound.
		c.enterNewRound(ti.BlockNumber, big.NewInt(0))
	case RoundStepNewRound:
		c.enterPropose(ti.BlockNumber, big.NewInt(0))
	case RoundStepPropose:
		c.enterPrevote(ti.BlockNumber, ti.Round)
	case RoundStepPrevoteWait:
		c.enterPrecommit(ti.BlockNumber, ti.Round)
	case RoundStepPrecommitWait:
		c.enterPrecommit(ti.BlockNumber, ti.Round)
		c.enterNewRound(ti.BlockNumber, big.NewInt(0).Add(ti.Round, big.NewInt(1)))
	default:
		panic(fmt.Sprintf("Invalid timeout step: %v", ti.Step))
	}

}
