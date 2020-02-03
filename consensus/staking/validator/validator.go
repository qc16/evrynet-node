package staking

import (
	"strings"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi"
	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
)

// ValidatorCaller is an auto generated read-only Go binding around an Evrynet contract.
type ValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NewValidatorCaller creates a new read-only instance of Validator, bound to a specific deployed contract.
func NewValidatorCaller(address common.Address, caller bind.ContractCaller, valABI string) (*ValidatorCaller, error) {
	abiParser, err := abi.JSON(strings.NewReader(valABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, abiParser, caller, nil, nil)
	return &ValidatorCaller{contract: contract}, nil
}

// GetValidators is a free data retrieval call binding the contract method.
//
// Solidity: function getCandidates(uint256) constant returns(address[])
func (val *ValidatorCaller) GetValidators(opts *bind.CallOpts, blockNumber uint64) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	// call to smart-contract to get validators by getValidators method
	err := val.contract.Call(opts, out, "getValidators", blockNumber)
	return *ret0, err
}
