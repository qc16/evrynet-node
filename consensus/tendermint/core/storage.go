package core

import (
	"math/big"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
)

//getStoredState init core with the last known roundState
//if there is no state in storage, init a new state.
func (c *core) getStoredState() *roundState {
	var (
		//these are default config at block 1 round 0
		rs                *roundState
		prevotesReceived  = make(map[int64]*messageSet)
		precommitReceived = make(map[int64]*messageSet)
		block             = types.NewBlockWithHeader(&types.Header{})
		view              = tendermint.View{
			Round:       0,
			BlockNumber: big.NewInt(1),
		}
		lockedRound      int64 = -1
		lockedBlock      types.Block
		validRound       int64 = -1
		validBlock       types.Block
		proposalReceived tendermint.Proposal
	)
	//TODO: Implement storage

	//if there is no stored roundState, init new one
	//TODO: init block 0
	if rs == nil {

		rs = newRoundState(&view, prevotesReceived, precommitReceived, block,
			lockedRound, &lockedBlock,
			validRound, &validBlock,
			&proposalReceived,
		)

		//TODO: timeout setup
	}
	return rs
}
