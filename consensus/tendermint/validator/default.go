package validator

import (
	"math"
	"sort"
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
)

type defaultValidator struct {
	address common.Address
}

// Address will return address of defaultValidator
func (val *defaultValidator) Address() common.Address {
	return val.address
}

// String will parse address of defaultValidator to string and return it
func (val *defaultValidator) String() string {
	return val.Address().String()
}

// ----------------------------------------------------------------------------

// defaultSet stores list of validator,
// proposer, proposer policy, proposer selector for voting
// and validator mutex to handle the conflict of the reader/writer between goroutines
type defaultSet struct {
	validators tendermint.Validators
	policy     tendermint.ProposerPolicy

	proposer    tendermint.Validator
	validatorMu sync.RWMutex
	selector    tendermint.ProposalSelector
}

func newDefaultSet(addrs []common.Address, policy tendermint.ProposerPolicy) *defaultSet {
	valSet := &defaultSet{}

	valSet.policy = policy
	// init validators
	valSet.validators = make([]tendermint.Validator, len(addrs))
	for i, addr := range addrs {
		valSet.validators[i] = New(addr)
	}
	// sort validator
	sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer = valSet.GetByIndex(0)
	}
	if policy == tendermint.Sticky {
		valSet.selector = stickyProposer
	} else {
		valSet.selector = roundRobinProposer
	}

	return valSet
}

// Size will return the length of validators in defaultSet
func (valSet *defaultSet) Size() int {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return len(valSet.validators)
}

// List will return validators in defaultSet
func (valSet *defaultSet) List() []tendermint.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

// GetByIndex will return validator by index in defaultSet
// If the index >= size of validators, return nil
func (valSet *defaultSet) GetByIndex(i int64) tendermint.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	if i < int64(valSet.Size()) {
		return valSet.validators[i]
	}
	return nil
}

// GetByAddress will return the validator & its index of the validator set by the address
// If the address does not exist in the validator set, it will return -1
func (valSet *defaultSet) GetByAddress(addr common.Address) (int, tendermint.Validator) {
	for i, val := range valSet.List() {
		if addr == val.Address() {
			return i, val
		}
	}
	return -1, nil
}

func calcSeed(valSet tendermint.ValidatorSet, proposer common.Address, round int64) int64 {
	offset := 0
	if idx, val := valSet.GetByAddress(proposer); val != nil {
		offset = idx
	}
	return int64(offset) + round
}

func emptyAddress(addr common.Address) bool {
	return addr == common.Address{}
}

func roundRobinProposer(valSet tendermint.ValidatorSet, proposer common.Address, round int64) tendermint.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	seed := int64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round) + 1
	}
	pick := seed % int64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func stickyProposer(valSet tendermint.ValidatorSet, proposer common.Address, round int64) tendermint.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	seed := int64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round)
	}
	pick := seed % int64(valSet.Size())
	return valSet.GetByIndex(pick)
}

// RemoveValidator will remove a validator from validatorset
func (valSet *defaultSet) RemoveValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for i, v := range valSet.validators {
		if v.Address() == address {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

// Copy allows copy all items from A to B
func (valSet *defaultSet) Copy() tendermint.ValidatorSet {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	addresses := make([]common.Address, 0, len(valSet.validators))
	for _, v := range valSet.validators {
		addresses = append(addresses, v.Address())
	}
	return NewSet(addresses, valSet.policy)
}

// F get the maximum number of faulty nodes
func (valSet *defaultSet) F() int { return int(math.Ceil(float64(valSet.Size())/3)) - 1 }

// Policy get proposal policy
func (valSet *defaultSet) Policy() tendermint.ProposerPolicy { return valSet.policy }

//CalcProposer implement valSet.CalcProposer. Based on the proposer selection scheme,
//it will set valSet.proposer to the address of the pre-determined round.
func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round int64) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	valSet.proposer = valSet.selector(valSet, lastProposer, round)
}

//GetProposer return the current proposer of this valSet
func (valSet *defaultSet) GetProposer() tendermint.Validator {
	return valSet.proposer
}
