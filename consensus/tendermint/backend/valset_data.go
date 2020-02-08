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

// ValSetData struct
type ValSetData struct {
	Epoch uint64
}

// NewValSetData returns new ValSetData
func NewValSetData(epochDuration uint64) *ValSetData {
	return &ValSetData{
		Epoch: epochDuration,
	}
}

//GetValSet keep tracks of validator set in associate with blockNumber
func (v *ValSetData) GetValSet(chainReader consensus.ChainReader, blockNumber *big.Int) (tendermint.ValidatorSet, error) {
	var (
		// get the checkpoint of block-number
		checkPoint    = utils.GetCheckpointNumber(v.Epoch, blockNumber.Uint64())
		valSet        = validator.NewSet([]common.Address{}, tendermint.RoundRobin, blockNumber.Int64())
		validatorAdds []common.Address
	)

	header := chainReader.GetHeaderByNumber(checkPoint)
	extra, err := types.ExtractTendermintExtra(header)
	if err != nil {
		log.Error("cannot load extra-data", "number", blockNumber, "error", err)
		return valSet, err
	}
	// The length of Validator set should be larger than 0
	if len(extra.ValidatorAdds) == 0 {
		log.Error("validator set is empty", "number", blockNumber)
		return valSet, tendermint.ErrEmptyValSet
	}
	if err = rlp.DecodeBytes(extra.ValidatorAdds, &validatorAdds); err != nil {
		log.Error("can't decode validator set from extra-data", "number", blockNumber)
		return valSet, err
	}

	return validator.NewSet(validatorAdds, tendermint.RoundRobin, blockNumber.Int64()), nil
}
