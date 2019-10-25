package miner

import (
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
)

type validator struct {
	address     common.Address
	vote        bool
	isLock      bool
	lockedBlock int64
	mu          *sync.RWMutex
}

func newProposedValidator() *validator {
	return &validator{
		mu: &sync.RWMutex{},
	}
}

func (v *validator) Address() common.Address {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.address
}

func (v *validator) Vote() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.vote {
		return "in"
	}
	return "out"
}

// setProposedValidator sets proposed validator in the block extra field
func (v *validator) setProposedValidator(address common.Address, vote bool) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.address = address
	v.vote = vote
	v.isLock = false
	return nil
}

// clearPendingProposedValidator remove the pending validator
func (v *validator) clearPendingProposedValidator() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.address = common.Address{}
	v.vote = false
	v.isLock = false
}

// getPendingProposedValidator returns pending validator
func (v *validator) getPendingProposedValidator() (validator common.Address, vote bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.address, v.vote
}

// isValidatorLocked returns whether validator is locked or not
func (v *validator) isValidatorLocked() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.isLock
}

// getLockBlock return block when the proposed validator is added to header
func (v *validator) getLockBlock() int64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.lockedBlock
}

// lockValidator lock proposed validator at a specific block
func (v *validator) lockValidator(blockNumber int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.isLock = true
	v.lockedBlock = blockNumber
}

// removeLock allow worker to propose validator
func (v *validator) removeLock() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.isLock = false
}
