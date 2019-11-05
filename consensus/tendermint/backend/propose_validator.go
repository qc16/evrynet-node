package backend

import (
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
)

// ProposalValidator store imformation (about address and vote true/false)
// for a candidate that is proposed from a node via the rpc api
type ProposalValidator struct {
	address    common.Address
	vote       bool
	isStick    bool
	stickBlock int64
	mu         *sync.RWMutex
}

func newProposedValidator() *ProposalValidator {
	return &ProposalValidator{
		mu: &sync.RWMutex{},
	}
}

// setProposedValidator sets proposed validator in the block extra field
func (v *ProposalValidator) setProposedValidator(address common.Address, vote bool) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.address = address
	v.vote = vote
	v.isStick = false
	// assign stickBlock equal zero to ensure that not clear when a previous proposed-validator when done at the engine.Seal
	v.stickBlock = 0
	return nil
}

// clearPendingProposedValidator remove the pending validator
func (v *ProposalValidator) clearPendingProposedValidator() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.address = common.Address{}
	v.vote = false
	v.isStick = false
}

// getPendingProposedValidator returns pending validator and validator's status lock
func (v *ProposalValidator) getPendingProposedValidator() (common.Address, bool, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.address, v.vote, v.isStick
}

// isValidatorStick returns whether validator is stick or not
func (v *ProposalValidator) isValidatorStick() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.isStick
}

// getStickBlock return block when the proposed validator is added to header
func (v *ProposalValidator) getStickBlock() int64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.stickBlock
}

// stickValidator stick proposed validator at a specific block
func (v *ProposalValidator) stickValidator(blockNumber int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.isStick = true
	v.stickBlock = blockNumber
}

// removeStick allow worker to propose validator
func (v *ProposalValidator) removeStick() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.isStick = false
}
