// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"io"
	"math/big"
	"sync"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/rlp"
)

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState(view *tendermint.View, validatorSet tendermint.ValidatorSet,
	proposalReceived *tendermint.Proposal, block *types.Block,
	lockedRound *big.Int, lockedBlock *types.Block,
	validRound *big.Int, validBlock *types.Block) *roundState {
	return &roundState{
		view:               view,
		block:              block,
		lockedRound:        lockedRound,
		lockedBlock:        lockedBlock,
		validRound:         validRound,
		validBlock:         validBlock,
		ProposalReceived:   proposalReceived,
		PrevotesReceived:   newMessageSet(validatorSet),
		PrecommitsReceived: newMessageSet(validatorSet),
		mu:                 new(sync.RWMutex),
	}
}

// roundState stores the consensus state
type roundState struct {
	view  *tendermint.View // view contains round and height
	block *types.Block     // current proposed block

	lockedRound *big.Int     // lockedRound is latest round it is locked
	lockedBlock *types.Block // lockedBlock is block it is locked at lockedRound above

	validRound *big.Int     // validRound is last known round with PoLC for non-nil valid block
	validBlock *types.Block // validBlock last known block of PoLC above

	ProposalReceived   *tendermint.Proposal //
	PrevotesReceived   *messageSet
	PrecommitsReceived *messageSet

	//step is the enumerate Step that currently the core is at.
	//to jump to the next step, UpdateRoundStep is called.
	step RoundStepType

	mu *sync.RWMutex
}

func (s *roundState) Step() RoundStepType{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.step
}


func (s *roundState) BlockNumber() *big.Int{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.view.BlockNumber
}

func (s *roundState) Round() *big.Int{
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.view.Round
}

func (s *roundState) UpdateRoundStep(round *big.Int, step RoundStepType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.view.Round = round
	s.step = step
}

func (s *roundState) SetProposalReceived(proposalReceived *tendermint.Proposal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ProposalReceived = proposalReceived
}

func (s *roundState) SetView(v *tendermint.View) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.view = v
}

// IsProposalComplete Returns true if the proposal block is complete &&
// (if POLRound was proposed, we have +2/3 prevotes from there).
func (s *roundState) IsProposalComplete() bool {
	//TODO: implement this, it have to do with number of votes receives (in handle prevotes)
	return true
}

func (s *roundState) View() *tendermint.View {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.view
}

func (s *roundState) SetBlock(bl *types.Block) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.block = bl
}

func (s *roundState) Block() *types.Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.block
}

func (s *roundState) SetLockedRoundAndBlock(lockedR *big.Int, lockedBl *types.Block) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.lockedRound = lockedR
	s.lockedBlock = lockedBl
}

func (s *roundState) LockedRound() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lockedRound
}

func (s *roundState) LockedBlock() *types.Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lockedBlock
}

func (s *roundState) SetValidRoundAndBlock(validR *big.Int, validBl *types.Block) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.validRound = validR
	s.validBlock = validBl
}

func (s *roundState) ValidRound() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.validRound
}

func (s *roundState) ValidBlock() *types.Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.validBlock
}

// The DecodeRLP method should read one value from the given
// Stream. It is not forbidden to read less or more, but it might
// be confusing.
func (s *roundState) DecodeRLP(stream *rlp.Stream) error {
	var ss struct {
		View               *tendermint.View
		Block              *types.Block
		LockedRound        *big.Int
		LockedBlock        *types.Block
		ValidRound         *big.Int
		ValidBlock         *types.Block
		ProposalReceived   *tendermint.Proposal
		PrevotesReceived   *messageSet
		PrecommitsReceived *messageSet
	}

	if err := stream.Decode(&ss); err != nil {
		return err
	}
	s.view, s.block = ss.View, ss.Block
	s.lockedRound, s.lockedBlock = ss.LockedRound, ss.LockedBlock
	s.validRound, s.validBlock = ss.ValidRound, ss.ValidBlock
	s.ProposalReceived = ss.ProposalReceived
	s.PrevotesReceived = ss.PrevotesReceived
	s.PrecommitsReceived = ss.PrecommitsReceived
	s.mu = new(sync.RWMutex)

	return nil
}

// EncodeRLP should write the RLP encoding of its receiver to w.
// If the implementation is a pointer method, it may also be
// called for nil pointers.
//
// Implementations should generate valid RLP. The data written is
// not verified at the moment, but a future version might. It is
// recommended to write only a single value but writing multiple
// values or no value at all is also permitted.
func (s *roundState) EncodeRLP(w io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return rlp.Encode(w, []interface{}{
		s.view,
		s.block,
		s.lockedRound,
		s.lockedBlock,
		s.validRound,
		s.validBlock,
		s.ProposalReceived,
		s.PrevotesReceived,
		s.PrecommitsReceived,
	})
}
