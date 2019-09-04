package core

import (
	"math/big"
	"sync"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
)

//getStoredState init core with the last known roundState
//if there is no state in storage, init a new state.
func (c *core) getStoredState() *roundState {
	var rs *roundState
	//TODO: Implement storage

	//if there is no stored roundState, init new one
	//TODO: init block 0
	if rs == nil {
		view := tendermint.View{
			Round:       big.NewInt(0),
			BlockNumber: big.NewInt(1),
		}
		rs = &roundState{
			view:  &view,
			mu:    &sync.RWMutex{},
			step:  RoundStepNewHeight,
			block: types.NewBlockWithHeader(&types.Header{}),
			//TODO: timeout setup
		}
	}
	return rs
}
