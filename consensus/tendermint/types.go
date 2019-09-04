package tendermint

import (
	"io"
	"math/big"

	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/rlp"
)

//Proposal represent a propose message to be sent in the case of the node is a proposer
//for its Round.
type Proposal struct {
	Block    *types.Block
	Round    *big.Int
	POLRound *big.Int
	//TODO: check if we need block Height
}

func (p *Proposal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		p.Block,
		p.Round,
		p.POLRound,
	})
}
