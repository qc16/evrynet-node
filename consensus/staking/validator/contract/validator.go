// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
        "math/big"
        "strings"

        ethereum "github.com/Evrynetlabs/evrynet-node"
        "github.com/Evrynetlabs/evrynet-node/accounts/abi"
        "github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
        "github.com/Evrynetlabs/evrynet-node/common"
        "github.com/Evrynetlabs/evrynet-node/core/types"
        "github.com/Evrynetlabs/evrynet-node/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
        _ = big.NewInt
        _ = strings.NewReader
        _ = ethereum.NotFound
        _ = abi.U256
        _ = bind.Bind
        _ = common.Big1
        _ = types.BloomLookup
        _ = event.NewSubscription
)

// ValidatorABI is the input ABI used to generate the binding from.
const ValidatorABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_validators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ValidatorBin is the compiled bytecode used for deploying new contracts.
const ValidatorBin = `0x608060405234801561001057600080fd5b506040516102bc3803806102bc8339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825186602082028301116401000000008211171561008557600080fd5b82525081516020918201928201910280838360005b838110156100b257818101518382015260200161009a565b5050505090500160405250505060008090505b81518110156101215760008282815181106100dc57fe5b60209081029190910181015182546001808201855560009485529290932090920180546001600160a01b0319166001600160a01b0390931692909217909155016100c5565b505061018a806101326000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806335aa2e441461003b578063b7ab4db514610074575b600080fd5b6100586004803603602081101561005157600080fd5b50356100cc565b604080516001600160a01b039092168252519081900360200190f35b61007c6100f3565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156100b85781810151838201526020016100a0565b505050509050019250505060405180910390f35b600081815481106100d957fe5b6000918252602090912001546001600160a01b0316905081565b6060600080548060200260200160405190810160405280929190818152602001828054801561014b57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161012d575b505050505090509056fea265627a7a723158208e8bbb32e84ca59dc7092ec9d2bf5b15df35c055b1de13f97817f456c200c2fb64736f6c63430005100032`

// DeployValidator deploys a new Evrynet contract, binding an instance of Validator to it.
func DeployValidator(auth *bind.TransactOpts, backend bind.ContractBackend, _validators []common.Address) (common.Address, *types.Transaction, *Validator, error) {
        parsed, err := abi.JSON(strings.NewReader(ValidatorABI))
        if err != nil {
                return common.Address{}, nil, nil, err
        }
        address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ValidatorBin), backend, _validators)
        if err != nil {
                return common.Address{}, nil, nil, err
        }
        return address, tx, &Validator{ValidatorCaller: ValidatorCaller{contract: contract}, ValidatorTransactor: ValidatorTransactor{contract: contract}, ValidatorFilterer: ValidatorFilterer{contract: contract}}, nil
}

// Validator is an auto generated Go binding around an Evrynet contract.
type Validator struct {
        ValidatorCaller     // Read-only binding to the contract
        ValidatorTransactor // Write-only binding to the contract
        ValidatorFilterer   // Log filterer for contract events
}

// ValidatorCaller is an auto generated read-only Go binding around an Evrynet contract.
type ValidatorCaller struct {
        contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorTransactor is an auto generated write-only Go binding around an Evrynet contract.
type ValidatorTransactor struct {
        contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorFilterer is an auto generated log filtering Go binding around an Evrynet contract events.
type ValidatorFilterer struct {
        contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSession is an auto generated Go binding around an Evrynet contract,
// with pre-set call and transact options.
type ValidatorSession struct {
        Contract     *Validator        // Generic contract binding to set the session for
        CallOpts     bind.CallOpts     // Call options to use throughout this session
        TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidatorCallerSession is an auto generated read-only Go binding around an Evrynet contract,
// with pre-set call options.
type ValidatorCallerSession struct {
        Contract *ValidatorCaller // Generic contract caller binding to set the session for
        CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ValidatorTransactorSession is an auto generated write-only Go binding around an Evrynet contract,
// with pre-set transact options.
type ValidatorTransactorSession struct {
        Contract     *ValidatorTransactor // Generic contract transactor binding to set the session for
        TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ValidatorRaw is an auto generated low-level Go binding around an Evrynet contract.
type ValidatorRaw struct {
        Contract *Validator // Generic contract binding to access the raw methods on
}

// ValidatorCallerRaw is an auto generated low-level read-only Go binding around an Evrynet contract.
type ValidatorCallerRaw struct {
        Contract *ValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Evrynet contract.
type ValidatorTransactorRaw struct {
        Contract *ValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidator creates a new instance of Validator, bound to a specific deployed contract.
func NewValidator(address common.Address, backend bind.ContractBackend) (*Validator, error) {
        contract, err := bindValidator(address, backend, backend, backend)
        if err != nil {
                return nil, err
        }
        return &Validator{ValidatorCaller: ValidatorCaller{contract: contract}, ValidatorTransactor: ValidatorTransactor{contract: contract}, ValidatorFilterer: ValidatorFilterer{contract: contract}}, nil
}

// NewValidatorCaller creates a new read-only instance of Validator, bound to a specific deployed contract.
func NewValidatorCaller(address common.Address, caller bind.ContractCaller) (*ValidatorCaller, error) {
        contract, err := bindValidator(address, caller, nil, nil)
        if err != nil {
                return nil, err
        }
        return &ValidatorCaller{contract: contract}, nil
}

// NewValidatorTransactor creates a new write-only instance of Validator, bound to a specific deployed contract.
func NewValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorTransactor, error) {
        contract, err := bindValidator(address, nil, transactor, nil)
        if err != nil {
                return nil, err
        }
        return &ValidatorTransactor{contract: contract}, nil
}

// NewValidatorFilterer creates a new log filterer instance of Validator, bound to a specific deployed contract.
func NewValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorFilterer, error) {
        contract, err := bindValidator(address, nil, nil, filterer)
        if err != nil {
                return nil, err
        }
        return &ValidatorFilterer{contract: contract}, nil
}

// bindValidator binds a generic wrapper to an already deployed contract.
func bindValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
        parsed, err := abi.JSON(strings.NewReader(ValidatorABI))
        if err != nil {
                return nil, err
        }
        return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validator *ValidatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
        return _Validator.Contract.ValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validator *ValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
        return _Validator.Contract.ValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validator *ValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
        return _Validator.Contract.ValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validator *ValidatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
        return _Validator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validator *ValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
        return _Validator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validator *ValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
        return _Validator.Contract.contract.Transact(opts, method, params...)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() constant returns(address[])
func (_Validator *ValidatorCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
        var (
                ret0 = new([]common.Address)
        )
        out := ret0
        err := _Validator.contract.Call(opts, out, "getValidators")
        return *ret0, err
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() constant returns(address[])
func (_Validator *ValidatorSession) GetValidators() ([]common.Address, error) {
        return _Validator.Contract.GetValidators(&_Validator.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() constant returns(address[])
func (_Validator *ValidatorCallerSession) GetValidators() ([]common.Address, error) {
        return _Validator.Contract.GetValidators(&_Validator.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Validator *ValidatorCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
        var (
                ret0 = new(common.Address)
        )
        out := ret0
        err := _Validator.contract.Call(opts, out, "validators", arg0)
        return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Validator *ValidatorSession) Validators(arg0 *big.Int) (common.Address, error) {
        return _Validator.Contract.Validators(&_Validator.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Validator *ValidatorCallerSession) Validators(arg0 *big.Int) (common.Address, error) {
        return _Validator.Contract.Validators(&_Validator.CallOpts, arg0)
}