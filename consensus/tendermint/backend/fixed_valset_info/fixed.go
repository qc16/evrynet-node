package fixed_valset_info

import (
	"fmt"
	"math/big"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/validator"
)

type FixedValidatorSetInfo struct {
	addresses []common.Address
}

func NewFixedValidatorSetInfo(addrs []common.Address) *FixedValidatorSetInfo {
	return &FixedValidatorSetInfo{
		addresses: addrs,
	}
}

//GetValSet keep tracks of validator set in associate with blockNumber
func (mvi *FixedValidatorSetInfo) GetValSet(chainReader consensus.ChainReader, blockNumber *big.Int) (tendermint.ValidatorSet, error) {
	fmt.Println("FixedValidatorSetInfo.GetValSet", "blockNumber", blockNumber, "validatorAdds", mvi.addresses)
	return validator.NewSet(mvi.addresses, tendermint.RoundRobin, blockNumber.Int64()), nil
}
