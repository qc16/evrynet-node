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
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evrynet-official/evrynet-client/common/math"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"golang.org/x/crypto/ed25519"
)

var (
	// MaxSignatureSize is a maximum allowed signature size for the Proposal
	// and Vote.
	// XXX: secp256k1 does not have Size nor MaxSize defined.
	MaxSignatureSize = math.MaxInt(ed25519.SignatureSize, 64)
)

// Proposal defines a block proposal for the consensus.
// It refers to the block by blockHash field.
// It must be signed by the correct proposer for the given Height/Round
// to be considered valid. It may depend on votes from a previous round,
// a so-called Proof-of-Lock (POL) round, as noted in the POLRound.
// If POLRound >= 0, then blockHash corresponds to the block that is locked in POLRound.
type Proposal struct {
	Type      SignedMsgType
	Height    *big.Int    `json:"height"`
	Round     *big.Int    `json:"round"`
	POLRound  *big.Int    `json:"pol_round"` // -1 if null.
	BlockHash common.Hash `json:"block_hash"`
	Timestamp time.Time   `json:"timestamp"`
	Signature []byte      `json:"signature"`
}

// NewProposal returns a new Proposal.
// If there is no POLRound, polRound should be -1.
func NewProposal(height *big.Int, round *big.Int, polRound *big.Int, blockHash common.Hash) *Proposal {
	return &Proposal{
		Type:      ProposalType,
		Height:    height,
		Round:     round,
		BlockHash: blockHash,
		POLRound:  polRound,
		Timestamp: time.Now(),
	}
}

// ValidateBasic performs basic validation.
func (p *Proposal) ValidateBasic() error {
	if p.Type != ProposalType {
		return errors.New("Invalid Type")
	}

	if p.Height.Cmp(big.NewInt(0)) < 0 {
		return errors.New("Negative Height")
	}
	if p.Round.Cmp(big.NewInt(0)) < 0 {
		return errors.New("Negative Round")
	}
	if p.POLRound.Cmp(big.NewInt(-1)) < 0 {
		return errors.New("Negative POLRound (exception: -1)")
	}
	// BlockHash must have size equal sha256.Size (tendermint's hash size)
	if len(p.BlockHash) != sha256.Size {
		return errors.New("Wrong BlockHash")
	}

	// NOTE: Timestamp validation is subtle and handled elsewhere.

	if len(p.Signature) == 0 {
		return errors.New("Signature is missing")
	}
	if len(p.Signature) > MaxSignatureSize {
		return errors.New("Signature is too big")
	}
	return nil
}

func (c *core) sendProposal(request *tendermint.Request) {
	// It is the proposer and it has the same height with the proposal
	if c.currentState.Height().Cmp(request.Proposal.Number()) == 0 && c.IsProposer() {
		type pp struct {
			View     *tendermint.View
			Proposal tendermint.Proposal
		}

		proposal, err := tendermint.Encode(&pp{
			View:     c.currentState.View(),
			Proposal: request.Proposal,
		})

		if err != nil {
			fmt.Errorf("Failed to encode: ", "view", c.currentState.View())
			return
		}

		c.broadcast(&message{
			Code: msgPropose,
			Msg:  proposal,
		})
	}
}

func (c *core) handleProposalReceived(msg *message, src tendermint.Validator) error {

	// Decode message
	var proposal struct {
		View     *tendermint.View
		Proposal tendermint.Proposal
	}

	err := msg.Decode(&proposal)
	if err != nil {
		return tendermint.ErrFailedDecodeProposal
	}

	// Make sure we have same view with the proposal message
	if err := c.checkMessage(msgPropose, proposal.View); err != nil {
		if err == tendermint.ErrOldMessage {
			// TODO: Consider to broadcast the COMMIT if received an older message
		}
		return err
	}

	// Check if message is sent from correct proposer
	if !c.valSet.IsProposer(src.Address()) {
		return tendermint.ErrIncorrectProposer
	}

	// Verify the proposal we've received
	if duration, err := c.backend.Verify(proposal.Proposal); err != nil {
		fmt.Errorf("Failed to verify proposal", "err", err, "duration", duration)
		// if it is a future block, should handle it again after the duration
		if err == consensus.ErrFutureBlock {
			// TODO:
			// Handle it again in the future after duration
		} else {
			// TODO:
			// Go to next round with round + 1
		}
		return err
	}

	// Accept the proposal
	// TODO: Need to check if we are at the step to accept the proposal (i.e waiting for proposal)

	return nil
}
