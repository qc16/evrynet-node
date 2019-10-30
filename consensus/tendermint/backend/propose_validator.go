package backend

import (
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
)

// ProposalValidator
type ProposalValidator struct {
	address     common.Address
	vote        bool
	isLock      bool
	lockedBlock int64
	mu          *sync.RWMutex
}

func newProposedValidator() *ProposalValidator {
	return &ProposalValidator{
		mu: &sync.RWMutex{},
	}
}

func (v *ProposalValidator) Address() common.Address {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.address
}

func (v *ProposalValidator) Vote() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.vote {
		return "in"
	}
	return "out"
}

// setProposedValidator sets proposed validator in the block extra field
func (v *ProposalValidator) setProposedValidator(address common.Address, vote bool) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.address = address
	v.vote = vote
	v.isLock = false
	return nil
}

// clearPendingProposedValidator remove the pending validator
func (v *ProposalValidator) clearPendingProposedValidator() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.address = common.Address{}
	v.vote = false
	v.isLock = false
}

// getPendingProposedValidator returns pending validator
func (v *ProposalValidator) getPendingProposedValidator() (validator common.Address, vote bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.address, v.vote
}

// isValidatorLocked returns whether validator is locked or not
func (v *ProposalValidator) isValidatorLocked() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.isLock
}

// getLockBlock return block when the proposed validator is added to header
func (v *ProposalValidator) getLockBlock() int64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.lockedBlock
}

// lockValidator lock proposed validator at a specific block
func (v *ProposalValidator) lockValidator(blockNumber int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.isLock = true
	v.lockedBlock = blockNumber
}

// removeLock allow worker to propose validator
func (v *ProposalValidator) removeLock() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.isLock = false
}
