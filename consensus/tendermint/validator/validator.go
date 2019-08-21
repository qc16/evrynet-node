package validator

import (
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/log"
)

// New will create new validator
func New(addr common.Address) tendermint.Validator {
	return &defaultValidator{
		address: addr,
	}
}

// NewSet will create new validator set by address list & policy
func NewSet(addrs []common.Address, policy tendermint.ProposerPolicy) tendermint.ValidatorSet {
	return newDefaultSet(addrs, policy)
}

// IsProposer will be checking whether the validator with given address is a proposer
func (valSet *defaultSet) IsProposer(address common.Address) bool {
	log.Warn("validator.IsProposer: implement me")
	//TODO: implement for this function to check is proposer
	return false
}

// ExtractValidators will extract extra data to address list
func ExtractValidators(extraData []byte) []common.Address {
	// get the validator addresses
	addrs := make([]common.Address, len(extraData)/common.AddressLength)
	for i := 0; i < len(addrs); i++ {
		copy(addrs[i][:], extraData[i*common.AddressLength:])
	}

	return addrs
}
