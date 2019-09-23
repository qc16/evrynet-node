package tendermint

import (
	"math/big"
	"strings"

	"github.com/evrynet-official/evrynet-client/common"
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
	GetByIndex(i int64) Validator
	// Get validator by given address
	// If the address does not exist in the validator set, it will return -1
	GetByAddress(addr common.Address) (int, Validator)
	// RemoveValidator remove the input validator from a list. It return false if the validator exist and is removed.
	// If the validator is not in the set, this function will return false
	RemoveValidator(address common.Address) bool
	// Copy validator set
	Copy() ValidatorSet
	// Get the maximum number of faulty nodes
	F() int
	// Get proposer policy
	Policy() ProposerPolicy
	// Check whether the validator with given address is a proposer
	IsProposer(address common.Address) bool
	// CalcProposer return the proposer for the different of round number indicated
	CalcProposer(lastProposer common.Address, roundDiff int64)
	// GetProposer return the current proposer
	GetProposer() Validator
	// Height return block height when valSet is init
	Height() int64
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, int64) Validator

// View includes a round number and a height of block we want to commit
type View struct {
	Round       int64
	BlockNumber *big.Int
}
