package backend

import (
	"math/big"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/utils"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/validator"
	"github.com/Evrynetlabs/evrynet-node/core/types"
	"github.com/Evrynetlabs/evrynet-node/log"
	"github.com/Evrynetlabs/evrynet-node/rlp"
)

// StakingValidatorInfo is implementation of ValidatorSetInfo
type StakingValidatorInfo struct {
}

// NewStakingValidatorInfo returns new StakingValidatorInfo
func NewStakingValidatorInfo() *StakingValidatorInfo {
	return &StakingValidatorInfo{}
}

//GetValSet keep tracks of validator set in associate with blockNumber
func (v *StakingValidatorInfo) GetValSet(chainReader consensus.ChainReader, number *big.Int) (tendermint.ValidatorSet, error) {
	var (
		// get the checkpoint of block-number
		blockNumber   = number.Int64()
		checkPoint    = utils.GetCheckpointNumber(chainReader.Config().Tendermint.Epoch, number.Uint64())
		valSet        = validator.NewSet([]common.Address{}, tendermint.RoundRobin, blockNumber)
		validatorAdds []common.Address
	)

	header := chainReader.GetHeaderByNumber(checkPoint)
	extra, err := types.ExtractTendermintExtra(header)
	if err != nil {
		log.Error("cannot load extra-data", "number", blockNumber, "error", err)
		return valSet, err
	}
	// The length of Validator's address should be larger than 0
	if len(extra.ValidatorAdds) == 0 {
		log.Error("validator' address is empty", "number", blockNumber)
		return valSet, tendermint.ErrEmptyValSet
	}
	if err = rlp.DecodeBytes(extra.ValidatorAdds, &validatorAdds); err != nil {
		log.Error("can't decode validator set from extra-data", "number", blockNumber)
		return valSet, err
	}

	return validator.NewSet(validatorAdds, tendermint.RoundRobin, blockNumber), nil
}
