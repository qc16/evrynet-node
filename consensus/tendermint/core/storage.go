package core

import (
	"math/big"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
)

//getStoredState init core with the last known roundState
//if there is no state in storage, init a new state.
func (c *core) getStoredState() *roundState {
	var rs *roundState
	//TODO: Implement storage

	//if there is no stored roundState, init new one
	if rs == nil {
		view := tendermint.View{
			Round:       big.NewInt(0),
			BlockNumber: big.NewInt(0),
		}
		rs = &roundState{
			view: &view,
			//TODO: timeout setup
		}
	}
	return rs
}
