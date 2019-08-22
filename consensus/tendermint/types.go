package tendermint

import (
	"io"
	"math/big"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/rlp"
)

// Proposal supports retrieving height and serialized block to be used during Tendermint consensus.
type Proposal interface {
	// Number retrieves the sequence number of this proposal.
	Number() *big.Int
	// Hash retrieves the hash of this proposal.
	Hash() common.Hash
	EncodeRLP(w io.Writer) error
	DecodeRLP(s *rlp.Stream) error
	String() string
}

type Request struct {
	Proposal Proposal
}

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
