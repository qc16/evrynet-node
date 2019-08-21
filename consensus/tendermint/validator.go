package tendermint

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Validator interface {
	// Address returns address
	Address() common.Address

	// String representation of Validator
	String() string
}

// ----------------------------------------------------------------------------

// Validators type is list of Validator
type Validators []Validator

// Len must be implemented for sort.Sort()
func (slice Validators) Len() int {
	return len(slice)
}

// Less must be implemented for sort.Sort()
func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].String(), slice[j].String()) < 0
}

// Swap must be implemented for sort.Sort()
func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ----------------------------------------------------------------------------

// ValidatorSet interface handles validator, proposer for defaultSet
type ValidatorSet interface {
	// Return the validator size
	Size() int
	// Return the validator array
	List() []Validator
	// Get validator by index
	GetByIndex(i uint64) Validator
	// Get validator by given address
	GetByAddress(addr common.Address) (int, Validator)
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, uint64) Validator

// View includes a round number and a height of block we want to commit
type View struct {
	Round  *big.Int
	Height *big.Int
}
