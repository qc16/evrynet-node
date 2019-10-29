package backend

import (
	"math/big"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
)

// TendermintAPI is a user facing RPC API to dump tendermint state
type TendermintAPI struct {
	chain consensus.ChainReader
	be    *backend
}

// ProposeValidator proposes a validator
// vote is false represents for kicking the validator out of network,
// vote is true represents for adding the validator to the network
func (api *TendermintAPI) ProposeValidator(address common.Address, vote bool) (bool, error) {
	if err := api.be.proposedValidator.setProposedValidator(address, vote); err != nil {
		return false, err
	}
	return true, nil
}

// ClearPendingProposedValidator this func will be removed a candidate that in the pending status
// returns true when done
func (api *TendermintAPI) ClearPendingProposedValidator() bool {
	api.be.proposedValidator.clearPendingProposedValidator()
	return true
}

// GetPendingProposedValidator returns the pending proposed validator
func (api *TendermintAPI) GetPendingProposedValidator() map[string]interface{} {
	validator, vote := api.be.proposedValidator.getPendingProposedValidator()
	return map[string]interface{}{
		"validator": validator,
		"vote":      vote,
	}
}

// GetValidators returns the list of validators by block's number
func (api *TendermintAPI) GetValidators(number *uint64) []common.Address {
	var (
		blockNumber *big.Int
	)
	if number == nil {
		blockNumber = api.chain.CurrentHeader().Number
	} else {
		blockNumber = new(big.Int).SetUint64(*number)
	}
	valSet := api.be.Validators(blockNumber)
	validators := make([]common.Address, 0, valSet.Size())
	for _, validator := range valSet.List() {
		validators = append(validators, validator.Address())
	}
	return validators
}
