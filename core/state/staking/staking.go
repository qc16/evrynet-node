package staking

import (
	"github.com/pkg/errors"

	"github.com/Evrynetlabs/evrynet-node/common"
)

var (
	// ErrEmptyValidatorSet returns when valset from smart-contract is empty
	ErrEmptyValidatorSet = errors.New("empty validator set")
	// ErrLengthOfCandidatesAndStakesMisMatch returns when lengths stakes and candidates are not match
	ErrLengthOfCandidatesAndStakesMisMatch        = errors.New("length of stakes is not equal to length of candidates")
	maxGasGetValSet                        uint64 = 500000000
)

type StakingCaller interface {
	GetValidators(common.Address) ([]common.Address, error)
}
