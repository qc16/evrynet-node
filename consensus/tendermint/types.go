package tendermint

import (
	"math/big"

	"github.com/evrynet-official/evrynet-client/core/types"
)


//Proposal represent a propose message to be sent in the case of the node is a proposer
//for its Round.
type Proposal struct {
	Block    *types.Block
	Round    *big.Int
	POLRound *big.Int
	//TODO: check if we need block Height
}
