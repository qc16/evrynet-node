
package backend

import (
	"bytes"
	
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
)

const (
	dbKeySnapshotPrefix = "tendermint-snapshot"
)

// Snapshot is the state of the authorization voting at a given point in time.
type Snapshot struct {
	Epoch uint64 // The number of blocks after which to checkpoint and reset the pending votes
	Number uint64                   // Block number where the snapshot was created
	Hash   common.Hash              // Block hash where the snapshot was created
	ValSet tendermint.ValidatorSet    // Set of authorized validators at this moment
}

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		Epoch:  s.Epoch,
		Number: s.Number,
		Hash:   s.Hash,
		ValSet: s.ValSet.Copy(),
	}
	//TODO: implement copy votes
	return cpy
}

// apply creates a new authorization snapshot by applying the given headers to
// the original one.
func (s *Snapshot) apply(headers []*types.Header) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(headers) == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(headers)-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, errInvalidVotingChain
		}
	}
	if headers[0].Number.Uint64() != s.Number+1 {
		return nil, errInvalidVotingChain
	}
	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	// TODO: Research & Implement more

	snap.Number += uint64(len(headers))
	snap.Hash = headers[len(headers)-1].Hash()

	return snap, nil
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator.Address())
	}
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			if bytes.Compare(validators[i][:], validators[j][:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}
