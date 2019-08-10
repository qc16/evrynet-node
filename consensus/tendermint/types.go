package tendermint

import (
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// Proposal supports retrieving height and serialized block to be used during Istanbul consensus.
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
