package core

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/rlp"
)

var (
	ErrInvalidProposalPOLRound     = errors.New("invalid proposal POL round")
	ErrInvalidProposalSignature    = errors.New("invalid proposal signature")
	ErrVoteHeightMismatch          = errors.New("vote height mismatch")
	ErrVoteInvalidValidatorAddress = errors.New("invalid validator address")
	ErrEmptyBlockProposal          = errors.New("empty block proposal")
	emptyBlockHash                 = common.Hash{}
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

	c.finalCommitted = c.backend.EventMux().Subscribe(
		tendermint.FinalCommittedEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
	c.finalCommitted.Unsubscribe()
}

// handleEvents will receive messages as well as timeout and is solely responsible for state change.
func (c *core) handleEvents() {
	var logger = c.getLogger()
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
				c.handleNewBlock(ev.Block)
			case tendermint.MessageEvent:
				//TODO: Handle ev.Payload, if got error then call c.backend.Gossip()
				var msg message
				if err := rlp.DecodeBytes(ev.Payload, &msg); err != nil {
					logger.Errorw("failed to decode msg", "error", err)
				} else {
					//log.Info("received message event", "from", msg.Address, "msg_Code", msg.Code)
					if err := c.handleMsg(msg); err != nil {
						logger.Errorw("failed to handle msg", "error", err)
					}
				}
			default:
				c.getLogger().Infow("Unknown event ", "event", ev)
			}
		case ti, ok := <-c.timeout.Chan(): //something from timeout...
			if !ok {
				return
			}
			c.handleTimeout(ti)
		case event, ok := <-c.finalCommitted.Chan():
			if !ok {
				return
			}
			switch ev := event.Data.(type) {
			case tendermint.FinalCommittedEvent:
				_ = c.handleFinalCommitted(ev.BlockNumber)
			}
		}
	}
}

// handleFinalCommitted is calling when received a final committed proposal
func (c *core) handleFinalCommitted(newHeadNumber *big.Int) error {
	var (
		state = c.CurrentState()
	)

	if state.BlockNumber().Cmp(newHeadNumber) > 0 {
		log.Warn("current state block number is ahead of new Head number. Ignore updating...",
			"current_block_number", state.BlockNumber().String(),
			"new_head_number", newHeadNumber.String())
		return nil
	}
	c.updateStateForNewblock()
	c.startRoundZero()
	return nil
}

func (c *core) handleNewBlock(block *types.Block) {
	var state = c.CurrentState()
	c.getLogger().Infow("received New Block event", "new_block_number", block.Number(), "new_block_hash", block.Hash().Hex())

	if block.Number() == nil || state.BlockNumber().Cmp(block.Number()) > 0 {
		//This is temporary to let miner come up with a newer block
		c.getLogger().Errorw("new block number is smaller than current block",
			"new_block_number", block.Number(), "state.BlockNumber", state.BlockNumber())
		//return a nil block to allow miner to send over a new one
		c.backend.Cancel(types.NewBlockWithHeader(&types.Header{
			Number: block.Number(),
		}))

		return
	}
	state.SetBlock(block)
}

//VerifyProposal validate msg & proposal when get from other nodes
func (c *core) VerifyProposal(proposal tendermint.Proposal, msg message) error {
	// Verify POLRound, which must be -1 or in range [0, proposal.Round).
	if proposal.POLRound < -1 ||
		((proposal.POLRound >= 0) && proposal.POLRound >= proposal.Round) {
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

	if proposal.Block == nil || (proposal.Block != nil && proposal.Block.Hash().Hex() == emptyBlockHash.Hex()) {
		return ErrEmptyBlockProposal
	}

	// verify the header of proposed block
	// ignore ErrEmptyCommittedSeals error because we don't have the committed seals yet
	if err := c.backend.VerifyProposalHeader(proposal.Block.Header()); err != nil && err != tendermint.ErrEmptyCommittedSeals {
		return err
	}

	// verify transaction hash & header
	if err := c.verifyTxs(proposal); err != nil {
		return err
	}

	return nil
}

func (c *core) verifyTxs(proposal tendermint.Proposal) error {
	var (
		block   = proposal.Block
		txs     = block.Transactions()
		txsHash = types.DeriveSha(txs)
	)

	// Verify txs hash
	if txsHash != block.Header().TxHash {
		return tendermint.ErrMismatchTxhashes
	}

	// Verify transaction for CoreTxPool
	if c.txPool != nil {
		for _, tx := range txs {
			if err := c.txPool.ValidateTx(tx, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *core) handlePropose(msg message) error {
	var (
		state    = c.CurrentState()
		proposal tendermint.Proposal
	)

	if err := rlp.DecodeBytes(msg.Msg, &proposal); err != nil {
		return err
	}
	logger := c.getLogger().With("proposal_round", proposal.Round, "proposal_block_hash", proposal.Block.Hash().Hex(),
		"proposal_block_number", proposal.Block.Number().String())
	logger.Infow("received a proposal", "from", msg.Address)

	// Already have one
	// TODO: possibly catch double proposals
	if state.ProposalReceived() != nil {
		return nil
	}

	// Does not apply, this is not an error but may happen due to network lattency
	if proposal.Block.Number().Cmp(state.BlockNumber()) != 0 || proposal.Round != state.Round() {
		logger.Warnw("received proposal with different height/round. Skip processing it")
		return nil
	}
	if err := c.VerifyProposal(proposal, msg); err != nil {
		return err
	}
	logger.Infow("setProposal receive...")

	state.SetProposalReceived(&proposal)
	//WARNING: THIS piece of code is experimental
	if state.Step() <= RoundStepPropose && state.IsProposalComplete() {
		log.Info("handle proposal: received proposal, proposal completed. before enterPrevote Jump to enterPrevote")
		// Move onto the next step
		c.enterPrevote(state.BlockNumber(), state.Round())
	} else if state.Step() == RoundStepCommit {
		log.Info("handle proposal: received proposal, proposal completed. at commit waiting for proposal block. Jump to finalizeCommit")

		// If we're waiting on the proposal block...
		c.finalizeCommit(proposal.Block.Number())
	} //// TODO: We can check if Proposal is for a different block as this is a sign of misbehavior!
	return nil
}

func (c *core) handlePrevote(msg message) error {
	var (
		vote  tendermint.Vote
		state = c.CurrentState()
	)
	if err := rlp.DecodeBytes(msg.Msg, &vote); err != nil {
		return err
	}

	if vote.BlockHash == nil || vote.BlockNumber == nil {
		c.getLogger().Panic("nil block hash is not allowed. Please make sure that prevote nil send an emptyBlockHash")
	}
	logger := c.getLogger().With("vote_block_number", vote.BlockNumber, "from", msg.Address, "vote_round", vote.Round, "block_hash", vote.BlockHash.Hex())

	if vote.BlockNumber.Cmp(state.BlockNumber()) != 0 {
		logger.Warnw("vote's block is different with current block", "vote_block", vote.BlockNumber, "from", msg.Address)
		if vote.BlockNumber.Cmp(state.BlockNumber()) > 0 {
			// vote from future block, save to future message queue
			logger.Infow("store prevote vote from future block", "vote_block", vote.BlockNumber, "vote_round", vote.Round, "from", msg.Address)
			if err := c.futureMessages.Enqueue(msg); err != nil {
				log.Error("failed to store future prevote message to queue", "err", err, "vote_block", vote.BlockNumber, "from", msg.Address)
			}
		}
		return nil
	}
	//log.Info("received prevote", "from", msg.Address, "round", vote.Round, "block_hash", vote.BlockHash.Hex())
	added, err := state.addPrevote(msg, &vote, c.valSet)
	if err != nil {
		return err
	}
	if !added {
		return nil
	}

	logger.Infow("added prevote vote into roundState")
	prevotes, ok := state.GetPrevotesByRound(vote.Round)
	if !ok {
		logger.Panic("expect prevotes to exist now")
	}
	//at this stage, state.PrevoteReceived[vote.Round] is guaranteed to exist.
	if blockHash, ok := prevotes.TwoThirdMajority(); ok {
		logger.Infow("got 2/3 majority on a block", "prevote_block", blockHash.Hex())
		var (
			lockedRound = state.LockedRound()
			lockedBlock = state.LockedBlock()
		)
		//if there is a lockedRound<vote.Round <= state.Round
		//and lockedBlock != nil
		if lockedRound != -1 && lockedRound < vote.Round && vote.Round <= state.Round() && lockedBlock.Hash().Hex() != blockHash.Hex() {
			logger.Infow("unlocking because of POL", "locked_round", lockedRound, "POL_round", vote.Round)
			state.Unlock()
		}

		//set valid Block if the polka is not emptyBlock
		if blockHash.Hex() != emptyBlockHash.Hex() && state.ValidRound() < vote.Round && vote.Round == state.Round() {
			if state.ProposalReceived() != nil && state.ProposalReceived().Block.Hash().Hex() == blockHash.Hex() {
				logger.Infow("updating validblock because of POL", "valid_round", state.ValidRound(), "POL_round", vote.Round)
				state.SetValidRoundAndBlock(vote.Round, state.ProposalReceived().Block)
			} else {
				logger.Infow("updating proposalBlock to nil since we received a valid block we don't know about")
				state.SetProposalReceived(nil)
			}
		}
	}
	//rebroadcast
	//note that tendermint doesn't do it, but it seems like this would speed up the process of gossiping
	//go func() {
	//	//We don't re-gossip if this is our own message
	//	if msg.Address.Hex() == c.backend.Address().Hex() {
	//		return
	//	}
	//	payload, err := rlp.EncodeToBytes(&msg)
	//	if err != nil {
	//		log.Error("failed to encode msg", "error", err)
	//		return
	//	}
	//	if err := c.backend.Gossip(c.valSet, payload); err != nil {
	//		log.Error("failed to re-gossip the vote received", "error", err)
	//	}
	//}()
	//if we receive a future roundthat come to 2/3 of prevotes on any block
	switch {
	case state.Round() < vote.Round && prevotes.HasTwoThirdAny():
		//Skip to vote.round
		c.enterNewRound(state.BlockNumber(), vote.Round)
	case state.Round() == vote.Round && RoundStepPrevote <= state.Step(): // current round
		blockHash, ok := prevotes.TwoThirdMajority()
		if ok && (state.IsProposalComplete() || blockHash.Hex() == emptyBlockHash.Hex()) {
			c.enterPrecommit(state.BlockNumber(), vote.Round)
		} else if prevotes.HasTwoThirdAny() {
			//wait till we got a majority
			c.enterPrevoteWait(state.BlockNumber(), vote.Round)
		}
	case state.ProposalReceived() != nil && 0 <= state.ProposalReceived().POLRound && state.ProposalReceived().POLRound == vote.Round:
		if state.IsProposalComplete() {
			c.enterPrevote(state.BlockNumber(), vote.Round)
		}
	}
	return nil
}

func (c *core) handlePrecommit(msg message) error {
	var (
		vote  tendermint.Vote
		state = c.CurrentState()
	)
	if err := rlp.DecodeBytes(msg.Msg, &vote); err != nil {
		return err
	}
	if vote.BlockHash == nil || vote.BlockNumber == nil {
		c.getLogger().Panic("nil block hash is not allowed. Please make sure that prevote nil send an emptyBlockHash")
	}

	logger := c.getLogger().With("vote_block", vote.BlockNumber, "vote_round", vote.Round,
		"from", msg.Address.Hex(), "block_hash", vote.BlockHash.Hex())
	if vote.BlockNumber.Cmp(state.BlockNumber()) != 0 {
		logger.Warnw("vote's block is different with current block", "vote_block", vote.BlockNumber, "from", msg.Address)
		if vote.BlockNumber.Cmp(state.BlockNumber()) > 0 {
			// vote from future block, save to future message queue
			logger.Infow("store precommit vote from future block", "vote_block", vote.BlockNumber, "vote_round", vote.Round, "from", msg.Address)
			if err := c.futureMessages.Enqueue(msg); err != nil {
				logger.Errorw("failed to store future prevote message to queue", "err", err, "vote_block", vote.BlockNumber, "from", msg.Address)
			}
		}
		logger.Warnw("vote's block is different with current block")
		return nil
	}
	//log.Info("received precommit", "from", msg.Address, "round", vote.Round, "block_hash", vote.BlockHash.Hex())
	added, err := state.addPrecommit(msg, &vote, c.valSet)
	if err != nil {
		return err
	}
	if !added {
		return nil
	}
	logger.Infow("added precommit vote into roundState")

	//TODO: revise if we need rebroadcast
	//rebroadcast
	//note that tendermint doesn't do it, but it seems like this would speed up the process of gossiping
	//go func() {
	//	//we don't re-gossip if this is our own message
	//	if msg.Address.Hex() == c.backend.Address().Hex() {
	//		return
	//	}
	//	payload, err := rlp.EncodeToBytes(&msg)
	//	if err != nil {
	//		log.Error("failed to encode msg", "error", err)
	//		return
	//	}
	//	if err := c.backend.Gossip(c.valSet, payload); err != nil {
	//		log.Error("failed to re-gossip the vote received", "error", err)
	//	}
	//}()

	precommits, ok := state.GetPrecommitsByRound(vote.Round)
	if !ok {
		panic("expect precommits to exist now")
	}
	//at this stage, state.PrevoteReceived[vote.Round] is guaranteed to exist.

	blockHash, ok := precommits.TwoThirdMajority()
	if ok {
		log.Info(" got 2/3 precommits  majority on a block", "block", blockHash)
		//this will go through the roundstep again to update core's roundState accordingly in case the vote Round is higher than core's Round
		c.enterNewRound(state.BlockNumber(), vote.Round)
		c.enterPrecommit(state.BlockNumber(), vote.Round)
		//if the precommit are not nil, enter commit
		if blockHash.Hex() != emptyBlockHash.Hex() {
			c.enterCommit(state.BlockNumber(), vote.Round)
			//TODO: if we need to skip when precommits has all votes
		} else {
			//wait for more precommit
			c.enterPrecommitWait(state.BlockNumber(), vote.Round)
		}
		return nil
	}

	//if there is no majority block
	if state.Round() <= vote.Round && precommits.HasTwoThirdAny() {
		//go through roundstep again to update round state
		c.enterNewRound(state.BlockNumber(), vote.Round)
		//wait for more precommit
		c.enterPrecommitWait(state.BlockNumber(), vote.Round)
	}
	return nil
}

func (c *core) handleMsg(msg message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch msg.Code {
	case msgPropose:
		return c.handlePropose(msg)
	case msgPrevote:
		return c.handlePrevote(msg)
	case msgPrecommit:
		return c.handlePrecommit(msg)
	default:
		return fmt.Errorf("unknown msg code %d", msg.Code)
	}
}

func (c *core) handleTimeout(ti timeoutInfo) {
	var (
		round       = c.CurrentState().Round()
		blockNumber = c.CurrentState().BlockNumber()
		step        = c.CurrentState().Step()
		logger      = c.getLogger().With("ti_block_number", ti.BlockNumber, "ti_round", ti.Round, "ti_step", ti.Step)
	)
	logger.Infow("Received timeout signal from core.timeout", "timeout", ti.Duration)
	// timeouts must be for current height, round, step
	if ti.BlockNumber.Cmp(blockNumber) != 0 || ti.Round < round || (ti.Round == round && ti.Step < step) {
		logger.Infow("Ignoring timeout because we're ahead")
		return
	}

	// the timeout will now cause a state transition
	c.mu.Lock()
	defer c.mu.Unlock()

	switch ti.Step {
	case RoundStepNewHeight:
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
		logger.Panicw("Invalid timeout step")
	}
}
