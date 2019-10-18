package core

import (
	"math/big"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
)

//getStoredState init core with the last known roundState
//if there is no state in storage, init a new state.
func (c *Core) getStoredState() *roundState {
	var (
		//these are default config at block 1 round 0
		rs *roundState

		prevotesReceived  = make(map[int64]*messageSet)
		precommitReceived = make(map[int64]*messageSet)
		block             = types.NewBlockWithHeader(&types.Header{})
		view              = tendermint.View{
			Round:       0,
			BlockNumber: big.NewInt(1),
		}
		lockedRound      int64 = -1
		lockedBlock      *types.Block
		validRound       int64 = -1
		validBlock       *types.Block
		proposalReceived *tendermint.Proposal
		step             = RoundStepNewHeight
	)

	//to continue from a stored State, get the last known block height
	lastKnownHeight := c.backend.CurrentHeadBlock().Number()

	// Increase block number to 1 block
	view.BlockNumber = new(big.Int).Add(lastKnownHeight, big.NewInt(1))

	rs = newRoundState(&view, prevotesReceived, precommitReceived, block,
		lockedRound, lockedBlock,
		validRound, validBlock,
		proposalReceived,
		step,
	)

	//TODO: timeout setup
	return rs
}
