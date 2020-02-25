// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking_contracts

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

// StakingContractsABI is the input ABI used to generate the binding from.
const StakingContractsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"minValidatorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unvote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMaxValidatorSize\",\"type\":\"uint256\"}],\"name\":\"updateMaxValidatorSize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"getTotalStakes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"}],\"name\":\"getVoterStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"getCandidateData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_isCandidate\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_latestTotalStakes\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxValidatorSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAllCandidates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"candidates\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newCap\",\"type\":\"uint256\"}],\"name\":\"updateMinVoteCap\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getVoterLatestStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"startBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochSize\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getListCandidates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_stakes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"_maxValSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_epochSize\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"vote\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"initCandidates\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawalCap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"candidateData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isCandidate\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"latestTotalStakes\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"resign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newCap\",\"type\":\"uint256\"}],\"name\":\"updateMinValidateStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minVoterCap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"candidatesOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxValidatorSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minValidatorStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minVoteCap\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unvoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"}]"

// StakingContractsBin is the compiled bytecode used for deploying new contracts.
const StakingContractsBin = `608060405260016000553480156200001657600080fd5b50604051620037c4380380620037c4833981810160405260e08110156200003c57600080fd5b81019080805160405193929190846401000000008211156200005d57600080fd5b838201915060208201858111156200007457600080fd5b82518660208202830111640100000000821117156200009257600080fd5b8083526020830192505050908051906020019060200280838360005b83811015620000cb578082015181840152602081019050620000ae565b50505050905001604052602001805190602001909291908051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050600085116200012557600080fd5b8460068190555083600781905550826008819055508160098190555086518410156200015057600080fd5b8660029080519060200190620001689291906200049e565b5060008090505b87518110156200042f5760405180606001604052806001151581526020018581526020018873ffffffffffffffffffffffffffffffffffffffff16815250600160008a8481518110620001be57fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548160ff0219169083151502179055506020820151816002015560408201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505083600160008a84815181106200028657fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008081526020019081526020016000208190555083600160008a84815181106200032f57fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555083600160008a8481518110620003c757fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160008081526020019081526020016000208190555080806001019150506200016f565b508660039080519060200190620004489291906200049e565b5080600a60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550436005819055505050505050505062000573565b8280548282559060005260206000209081019282156200051a579160200282015b82811115620005195782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190620004bf565b5b5090506200052991906200052d565b5090565b6200057091905b808211156200056c57600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690555060010162000534565b5090565b90565b61324180620005836000396000f3fe60806040526004361061019c5760003560e01c8063572d356e116100ec578063aa6773541161008a578063b5b7a18411610064578063b5b7a18414610a80578063b97dd9e214610aab578063f851a44014610ad6578063f8ac9dd514610b2d5761019c565b8063aa67735414610983578063ae6e43f5146109f4578063b2c76f1014610a455761019c565b806375829def116100c657806375829def146107a55780638106d590146107f6578063909b40531461087157806391a9634f146108e05761019c565b8063572d356e14610650578063690ff8a1146106875780636dd7d8ea146107615761019c565b80632de7dd5f116101595780633477ee2e116101335780633477ee2e146104ea5780633a1d8c5a14610565578063432fc981146105a057806348cd4cb1146106255761019c565b80632de7dd5f146104005780632e1a7d4d1461042b5780632e6997fe1461047e5761019c565b8063017ddd351461019e57806302aa9be2146101c95780630619624f14610224578063118ea1af1461025f5780631fad0a71146102ce5780632a466ac71461035d575b005b3480156101aa57600080fd5b506101b3610b58565b6040518082815260200191505060405180910390f35b3480156101d557600080fd5b50610222600480360360408110156101ec57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610b5e565b005b34801561023057600080fd5b5061025d6004803603602081101561024757600080fd5b81019080803590602001909291905050506111b4565b005b34801561026b57600080fd5b506102b86004803603604081101561028257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611218565b6040518082815260200191505060405180910390f35b3480156102da57600080fd5b50610347600480360360608110156102f157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611276565b6040518082815260200191505060405180910390f35b34801561036957600080fd5b506103ac6004803603602081101561038057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611312565b60405180841515151581526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390f35b34801561040c57600080fd5b5061041561141a565b6040518082815260200191505060405180910390f35b34801561043757600080fd5b506104646004803603602081101561044e57600080fd5b8101908080359060200190929190505050611420565b604051808215151515815260200191505060405180910390f35b34801561048a57600080fd5b50610493611625565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156104d65780820151818401526020810190506104bb565b505050509050019250505060405180910390f35b3480156104f657600080fd5b506105236004803603602081101561050d57600080fd5b81019080803590602001909291905050506116b3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561057157600080fd5b5061059e6004803603602081101561058857600080fd5b81019080803590602001909291905050506116ef565b005b3480156105ac57600080fd5b5061060f600480360360408110156105c357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611753565b6040518082815260200191505060405180910390f35b34801561063157600080fd5b5061063a6117dd565b6040518082815260200191505060405180910390f35b34801561065c57600080fd5b506106656117e3565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b34801561069357600080fd5b5061069c6117f2565b6040518080602001806020018563ffffffff1663ffffffff1681526020018463ffffffff1663ffffffff168152602001838103835287818151815260200191508051906020019060200280838360005b838110156107075780820151818401526020810190506106ec565b50505050905001838103825286818151815260200191508051906020019060200280838360005b8381101561074957808201518184015260208101905061072e565b50505050905001965050505050505060405180910390f35b6107a36004803603602081101561077757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611e27565b005b3480156107b157600080fd5b506107f4600480360360208110156107c857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506122d4565b005b34801561080257600080fd5b5061082f6004803603602081101561081957600080fd5b81019080803590602001909291905050506123ac565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561087d57600080fd5b506108ca6004803603604081101561089457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506123e8565b6040518082815260200191505060405180910390f35b3480156108ec57600080fd5b5061092f6004803603602081101561090357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061240d565b60405180841515151581526020018381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001935050505060405180910390f35b34801561098f57600080fd5b506109f2600480360360408110156109a657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612464565b005b348015610a0057600080fd5b50610a4360048036036020811015610a1757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061284f565b005b348015610a5157600080fd5b50610a7e60048036036020811015610a6857600080fd5b8101908080359060200190929190505050612fea565b005b348015610a8c57600080fd5b50610a9561304e565b6040518082815260200191505060405180910390f35b348015610ab757600080fd5b50610ac0613054565b6040518082815260200191505060405180910390f35b348015610ae257600080fd5b50610aeb613084565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b348015610b3957600080fd5b50610b426130aa565b6040518082815260200191505060405180910390f35b60085481565b6001600080828254019250508190555060008054905060008211610bea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f776974686472617720616d6f756e74206d75737420626520706f73697469766581525060200191505060405180910390fd5b6000610bf4613054565b905060003390506000600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905084811015610cf5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f616d6f756e7420746f6f2062696720746f20776974686472617700000000000081525060200191505060405180910390fd5b6000610d0a86836130b090919063ffffffff16565b9050600160008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610e0257600854811015610dfd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806131b46021913960400191505060405180910390fd5b610e69565b6000811480610e1357506009548110155b610e68576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260388152602001806131d56038913960400191505060405180910390fd5b5b80600160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555080600160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600086815260200190815260200160002081905550610fd786600160008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201546130b090919063ffffffff16565b600160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020181905550600160008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020154600160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001016000868152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff166108fc879081150290604051600060405180830381858888f193505050501580156110fd573d6000803e3d6000fd5b507f7958395da8e26969cc7c09ee58e9507a2601574c3bd232617e2d6354224ff836838888604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15050505060005481146111af57600080fd5b505050565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461120e57600080fd5b8060078190555050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101600083815260200190815260200160002054905092915050565b6000600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008381526020019081526020016000205490509392505050565b6000806000600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff169250600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169150600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206002015490509193909250565b60075481565b6000600160008082825401925050819055506000805490506000611442613054565b90508381101561149d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806131936021913960400191505060405180910390fd5b60003390506000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008781526020019081526020016000205490506000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600088815260200190815260200160002081905550600081116115c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f776974686472617720636170206973203000000000000000000000000000000081525060200191505060405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015611609573d6000803e3d6000fd5b5060019450505050600054811461161f57600080fd5b50919050565b606060028054806020026020016040519081016040528092919081815260200182805480156116a957602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001906001019080831161165f575b5050505050905090565b600281815481106116c057fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461174957600080fd5b8060098190555050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60055481565b60006117ed613054565b905090565b6060806000806007549150611805613054565b90506000611811613054565b90506000811415611956576003805490506040519080825280602002602001820160405280156118505781602001602082028038833980820191505090505b5094506003805490506040519080825280602002602001820160405280156118875781602001602082028038833980820191505090505b50935060008090505b60038054905081101561194357600381815481106118aa57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168682815181106118e157fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060085485828151811061192a57fe5b6020026020010181815250508080600101915050611890565b5084848484945094509450945050611e21565b6000600182039050600080600090505b600280549050811015611af8576000600160006002848154811061198657fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506008546001600060028581548110611a2557fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008681526020019081526020016000205410611aea5782806001019350505b508080600101915050611966565b5080604051908082528060200260200182016040528015611b285781602001602082028038833980820191505090505b50965080604051908082528060200260200182016040528015611b5a5781602001602082028038833980820191505090505b509550600080905060008090505b600280549050811015611e1b5760006001600060028481548110611b8857fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506008546001600060028581548110611c2757fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008781526020019081526020016000205410611e0d5760028281548110611cf057fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168a8481518110611d2757fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250506001600060028481548110611d7257fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101600086815260200190815260200160002054898481518110611df857fe5b60200260200101818152505082806001019350505b508080600101915050611b68565b50505050505b90919293565b600954341015611e3657600080fd5b8060011515600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff16151514611e9757600080fd5b60003490506000339050611ef682600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201546130ca90919063ffffffff16565b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201819055506000611f46613054565b9050611fae83600160008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001016000848152602001908152602001600020546130ca90919063ffffffff16565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160008381526020019081526020016000208190555061209783600160008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546130ca90919063ffffffff16565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000838152602001908152602001600020819055507f174ba19ba3c8bb5c679c87e51db79fff7c3f04bb84c1fd55b7dacb470b674aa6828685604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15050505050565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461232e57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561236857600080fd5b80600a60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600381815481106123b957fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6004602052816000526040600020602052806000526040600020600091509150505481565b60016020528060005260406000206000915090508060000160009054906101000a900460ff16908060020154908060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146124be57600080fd5b8160001515600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff1615151461251f57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156125c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f5f63616e6469646174652061646472657373206973206d697373696e6700000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415612665576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f5f6f776e65722061646472657373206973206d697373696e670000000000000081525060200191505060405180910390fd5b6040518060600160405280600115158152602001600081526020018373ffffffffffffffffffffffffffffffffffffffff16815250600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548160ff0219169083151502179055506020820151816002015560408201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505060028390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550507f0a31ee9d46a828884b81003c8498156ea6aa15b9b54bdd0ef0b533d9eba57e558383604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a1505050565b600160008082825401925050819055506000805490508160011515600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff161515146128c657600080fd5b60003390508073ffffffffffffffffffffffffffffffffffffffff16600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461296557600080fd5b60008090505b600280549050811015612ae2578473ffffffffffffffffffffffffffffffffffffffff166002828154811061299c57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415612ad5576002600160028054905003815481106129f857fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660028281548110612a3057fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600260016002805490500381548110612a8d57fe5b9060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556002805480919060019003612acf9190613141565b50612ae2565b808060010191505061296b565b506000600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a81548160ff0219169083151502179055506000612b48613054565b90506000600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490506000612be5600854836130e990919063ffffffff16565b90506000612bfc82846130b090919063ffffffff16565b9050612c5383600160008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201546130b090919063ffffffff16565b600160008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020181905550600160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020154600160008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001016000868152602001908152602001600020819055506000600160008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000600160008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000868152602001908152602001600020819055506000612e646002866130ca90919063ffffffff16565b9050612ec983600460008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000848152602001908152602001600020546130ca90919063ffffffff16565b600460008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000838152602001908152602001600020819055506000821115612f6e578573ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015612f6c573d6000803e3d6000fd5b505b7fa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d8689604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1505050505050506000548114612fe657600080fd5b5050565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461304457600080fd5b8060088190555050565b60065481565b600061307f600654613071600554436130b090919063ffffffff16565b61310290919063ffffffff16565b905090565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60095481565b6000828211156130bf57600080fd5b818303905092915050565b6000808284019050838110156130df57600080fd5b8091505092915050565b60008183106130f857816130fa565b825b905092915050565b600080821161311057600080fd5b600082848161311b57fe5b04905082848161312757fe5b0681840201841461313757600080fd5b8091505092915050565b81548183558181111561316857818360005260206000209182019101613167919061316d565b5b505050565b61318f91905b8082111561318b576000816000905550600101613173565b5090565b9056fe63616e206e6f7420776974686472617720666f72206675747572652065706f636872656d61696e20616d6f756e74206f66206f776e657220697320746f6f206c6f7772656d61696e20616d6f756e74206d757374206265206569746865722030206f72206174206c65617374206d696e20766f74657220636170a265627a7a7231582061199783fd26b2309a726f5e316ef2a5803572944e0c339f63cd9bd624ae735864736f6c634300050b0032`

// DeployStakingContracts deploys a new Evrynet contract, binding an instance of StakingContracts to it.
func DeployStakingContracts(auth *bind.TransactOpts, backend bind.ContractBackend, _candidates []common.Address, candidatesOwner common.Address, _epochPeriod *big.Int, _maxValidatorSize *big.Int, _minValidatorStake *big.Int, _minVoteCap *big.Int, _admin common.Address) (common.Address, *types.Transaction, *StakingContracts, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingContractsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StakingContractsBin), backend, _candidates, candidatesOwner, _epochPeriod, _maxValidatorSize, _minValidatorStake, _minVoteCap, _admin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingContracts{StakingContractsCaller: StakingContractsCaller{contract: contract}, StakingContractsTransactor: StakingContractsTransactor{contract: contract}, StakingContractsFilterer: StakingContractsFilterer{contract: contract}}, nil
}

// StakingContracts is an auto generated Go binding around an Evrynet contract.
type StakingContracts struct {
	StakingContractsCaller     // Read-only binding to the contract
	StakingContractsTransactor // Write-only binding to the contract
	StakingContractsFilterer   // Log filterer for contract events
}

// StakingContractsCaller is an auto generated read-only Go binding around an Evrynet contract.
type StakingContractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractsTransactor is an auto generated write-only Go binding around an Evrynet contract.
type StakingContractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractsFilterer is an auto generated log filtering Go binding around an Evrynet contract events.
type StakingContractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingContractsSession is an auto generated Go binding around an Evrynet contract,
// with pre-set call and transact options.
type StakingContractsSession struct {
	Contract     *StakingContracts // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingContractsCallerSession is an auto generated read-only Go binding around an Evrynet contract,
// with pre-set call options.
type StakingContractsCallerSession struct {
	Contract *StakingContractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// StakingContractsTransactorSession is an auto generated write-only Go binding around an Evrynet contract,
// with pre-set transact options.
type StakingContractsTransactorSession struct {
	Contract     *StakingContractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// StakingContractsRaw is an auto generated low-level Go binding around an Evrynet contract.
type StakingContractsRaw struct {
	Contract *StakingContracts // Generic contract binding to access the raw methods on
}

// StakingContractsCallerRaw is an auto generated low-level read-only Go binding around an Evrynet contract.
type StakingContractsCallerRaw struct {
	Contract *StakingContractsCaller // Generic read-only contract binding to access the raw methods on
}

// StakingContractsTransactorRaw is an auto generated low-level write-only Go binding around an Evrynet contract.
type StakingContractsTransactorRaw struct {
	Contract *StakingContractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingContracts creates a new instance of StakingContracts, bound to a specific deployed contract.
func NewStakingContracts(address common.Address, backend bind.ContractBackend) (*StakingContracts, error) {
	contract, err := bindStakingContracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingContracts{StakingContractsCaller: StakingContractsCaller{contract: contract}, StakingContractsTransactor: StakingContractsTransactor{contract: contract}, StakingContractsFilterer: StakingContractsFilterer{contract: contract}}, nil
}

// NewStakingContractsCaller creates a new read-only instance of StakingContracts, bound to a specific deployed contract.
func NewStakingContractsCaller(address common.Address, caller bind.ContractCaller) (*StakingContractsCaller, error) {
	contract, err := bindStakingContracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingContractsCaller{contract: contract}, nil
}

// NewStakingContractsTransactor creates a new write-only instance of StakingContracts, bound to a specific deployed contract.
func NewStakingContractsTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingContractsTransactor, error) {
	contract, err := bindStakingContracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingContractsTransactor{contract: contract}, nil
}

// NewStakingContractsFilterer creates a new log filterer instance of StakingContracts, bound to a specific deployed contract.
func NewStakingContractsFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingContractsFilterer, error) {
	contract, err := bindStakingContracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingContractsFilterer{contract: contract}, nil
}

// bindStakingContracts binds a generic wrapper to an already deployed contract.
func bindStakingContracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingContractsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingContracts *StakingContractsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StakingContracts.Contract.StakingContractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingContracts *StakingContractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContracts.Contract.StakingContractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingContracts *StakingContractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingContracts.Contract.StakingContractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingContracts *StakingContractsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StakingContracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingContracts *StakingContractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingContracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingContracts *StakingContractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingContracts.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_StakingContracts *StakingContractsCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_StakingContracts *StakingContractsSession) Admin() (common.Address, error) {
	return _StakingContracts.Contract.Admin(&_StakingContracts.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() constant returns(address)
func (_StakingContracts *StakingContractsCallerSession) Admin() (common.Address, error) {
	return _StakingContracts.Contract.Admin(&_StakingContracts.CallOpts)
}

// CandidateData is a free data retrieval call binding the contract method 0x91a9634f.
//
// Solidity: function candidateData(address ) constant returns(bool isCandidate, uint256 latestTotalStakes, address owner)
func (_StakingContracts *StakingContractsCaller) CandidateData(opts *bind.CallOpts, arg0 common.Address) (struct {
	IsCandidate       bool
	LatestTotalStakes *big.Int
	Owner             common.Address
}, error) {
	ret := new(struct {
		IsCandidate       bool
		LatestTotalStakes *big.Int
		Owner             common.Address
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "candidateData", arg0)
	return *ret, err
}

// CandidateData is a free data retrieval call binding the contract method 0x91a9634f.
//
// Solidity: function candidateData(address ) constant returns(bool isCandidate, uint256 latestTotalStakes, address owner)
func (_StakingContracts *StakingContractsSession) CandidateData(arg0 common.Address) (struct {
	IsCandidate       bool
	LatestTotalStakes *big.Int
	Owner             common.Address
}, error) {
	return _StakingContracts.Contract.CandidateData(&_StakingContracts.CallOpts, arg0)
}

// CandidateData is a free data retrieval call binding the contract method 0x91a9634f.
//
// Solidity: function candidateData(address ) constant returns(bool isCandidate, uint256 latestTotalStakes, address owner)
func (_StakingContracts *StakingContractsCallerSession) CandidateData(arg0 common.Address) (struct {
	IsCandidate       bool
	LatestTotalStakes *big.Int
	Owner             common.Address
}, error) {
	return _StakingContracts.Contract.CandidateData(&_StakingContracts.CallOpts, arg0)
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsCaller) Candidates(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "candidates", arg0)
	return *ret0, err
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsSession) Candidates(arg0 *big.Int) (common.Address, error) {
	return _StakingContracts.Contract.Candidates(&_StakingContracts.CallOpts, arg0)
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsCallerSession) Candidates(arg0 *big.Int) (common.Address, error) {
	return _StakingContracts.Contract.Candidates(&_StakingContracts.CallOpts, arg0)
}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) EpochPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "epochPeriod")
	return *ret0, err
}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() constant returns(uint256)
func (_StakingContracts *StakingContractsSession) EpochPeriod() (*big.Int, error) {
	return _StakingContracts.Contract.EpochPeriod(&_StakingContracts.CallOpts)
}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) EpochPeriod() (*big.Int, error) {
	return _StakingContracts.Contract.EpochPeriod(&_StakingContracts.CallOpts)
}

// EpochSize is a free data retrieval call binding the contract method 0x572d356e.
//
// Solidity: function epochSize() constant returns(uint32)
func (_StakingContracts *StakingContractsCaller) EpochSize(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "epochSize")
	return *ret0, err
}

// EpochSize is a free data retrieval call binding the contract method 0x572d356e.
//
// Solidity: function epochSize() constant returns(uint32)
func (_StakingContracts *StakingContractsSession) EpochSize() (uint32, error) {
	return _StakingContracts.Contract.EpochSize(&_StakingContracts.CallOpts)
}

// EpochSize is a free data retrieval call binding the contract method 0x572d356e.
//
// Solidity: function epochSize() constant returns(uint32)
func (_StakingContracts *StakingContractsCallerSession) EpochSize() (uint32, error) {
	return _StakingContracts.Contract.EpochSize(&_StakingContracts.CallOpts)
}

// GetAllCandidates is a free data retrieval call binding the contract method 0x2e6997fe.
//
// Solidity: function getAllCandidates() constant returns(address[] _candidates)
func (_StakingContracts *StakingContractsCaller) GetAllCandidates(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getAllCandidates")
	return *ret0, err
}

// GetAllCandidates is a free data retrieval call binding the contract method 0x2e6997fe.
//
// Solidity: function getAllCandidates() constant returns(address[] _candidates)
func (_StakingContracts *StakingContractsSession) GetAllCandidates() ([]common.Address, error) {
	return _StakingContracts.Contract.GetAllCandidates(&_StakingContracts.CallOpts)
}

// GetAllCandidates is a free data retrieval call binding the contract method 0x2e6997fe.
//
// Solidity: function getAllCandidates() constant returns(address[] _candidates)
func (_StakingContracts *StakingContractsCallerSession) GetAllCandidates() ([]common.Address, error) {
	return _StakingContracts.Contract.GetAllCandidates(&_StakingContracts.CallOpts)
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isCandidate, address _owner, uint256 _latestTotalStakes)
func (_StakingContracts *StakingContractsCaller) GetCandidateData(opts *bind.CallOpts, _candidate common.Address) (struct {
	IsCandidate       bool
	Owner             common.Address
	LatestTotalStakes *big.Int
}, error) {
	ret := new(struct {
		IsCandidate       bool
		Owner             common.Address
		LatestTotalStakes *big.Int
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "getCandidateData", _candidate)
	return *ret, err
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isCandidate, address _owner, uint256 _latestTotalStakes)
func (_StakingContracts *StakingContractsSession) GetCandidateData(_candidate common.Address) (struct {
	IsCandidate       bool
	Owner             common.Address
	LatestTotalStakes *big.Int
}, error) {
	return _StakingContracts.Contract.GetCandidateData(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isCandidate, address _owner, uint256 _latestTotalStakes)
func (_StakingContracts *StakingContractsCallerSession) GetCandidateData(_candidate common.Address) (struct {
	IsCandidate       bool
	Owner             common.Address
	LatestTotalStakes *big.Int
}, error) {
	return _StakingContracts.Contract.GetCandidateData(&_StakingContracts.CallOpts, _candidate)
}

// GetCurrentEpoch is a free data retrieval call binding the contract method 0xb97dd9e2.
//
// Solidity: function getCurrentEpoch() constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) GetCurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getCurrentEpoch")
	return *ret0, err
}

// GetCurrentEpoch is a free data retrieval call binding the contract method 0xb97dd9e2.
//
// Solidity: function getCurrentEpoch() constant returns(uint256)
func (_StakingContracts *StakingContractsSession) GetCurrentEpoch() (*big.Int, error) {
	return _StakingContracts.Contract.GetCurrentEpoch(&_StakingContracts.CallOpts)
}

// GetCurrentEpoch is a free data retrieval call binding the contract method 0xb97dd9e2.
//
// Solidity: function getCurrentEpoch() constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) GetCurrentEpoch() (*big.Int, error) {
	return _StakingContracts.Contract.GetCurrentEpoch(&_StakingContracts.CallOpts)
}

// GetListCandidates is a free data retrieval call binding the contract method 0x690ff8a1.
//
// Solidity: function getListCandidates() constant returns(address[] _candidates, uint256[] _stakes, uint32 _maxValSize, uint32 _epochSize)
func (_StakingContracts *StakingContractsCaller) GetListCandidates(opts *bind.CallOpts) (struct {
	Candidates []common.Address
	Stakes     []*big.Int
	MaxValSize uint32
	EpochSize  uint32
}, error) {
	ret := new(struct {
		Candidates []common.Address
		Stakes     []*big.Int
		MaxValSize uint32
		EpochSize  uint32
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "getListCandidates")
	return *ret, err
}

// GetListCandidates is a free data retrieval call binding the contract method 0x690ff8a1.
//
// Solidity: function getListCandidates() constant returns(address[] _candidates, uint256[] _stakes, uint32 _maxValSize, uint32 _epochSize)
func (_StakingContracts *StakingContractsSession) GetListCandidates() (struct {
	Candidates []common.Address
	Stakes     []*big.Int
	MaxValSize uint32
	EpochSize  uint32
}, error) {
	return _StakingContracts.Contract.GetListCandidates(&_StakingContracts.CallOpts)
}

// GetListCandidates is a free data retrieval call binding the contract method 0x690ff8a1.
//
// Solidity: function getListCandidates() constant returns(address[] _candidates, uint256[] _stakes, uint32 _maxValSize, uint32 _epochSize)
func (_StakingContracts *StakingContractsCallerSession) GetListCandidates() (struct {
	Candidates []common.Address
	Stakes     []*big.Int
	MaxValSize uint32
	EpochSize  uint32
}, error) {
	return _StakingContracts.Contract.GetListCandidates(&_StakingContracts.CallOpts)
}

// GetTotalStakes is a free data retrieval call binding the contract method 0x118ea1af.
//
// Solidity: function getTotalStakes(address _candidate, uint256 epoch) constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) GetTotalStakes(opts *bind.CallOpts, _candidate common.Address, epoch *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getTotalStakes", _candidate, epoch)
	return *ret0, err
}

// GetTotalStakes is a free data retrieval call binding the contract method 0x118ea1af.
//
// Solidity: function getTotalStakes(address _candidate, uint256 epoch) constant returns(uint256)
func (_StakingContracts *StakingContractsSession) GetTotalStakes(_candidate common.Address, epoch *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.GetTotalStakes(&_StakingContracts.CallOpts, _candidate, epoch)
}

// GetTotalStakes is a free data retrieval call binding the contract method 0x118ea1af.
//
// Solidity: function getTotalStakes(address _candidate, uint256 epoch) constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) GetTotalStakes(_candidate common.Address, epoch *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.GetTotalStakes(&_StakingContracts.CallOpts, _candidate, epoch)
}

// GetVoterLatestStake is a free data retrieval call binding the contract method 0x432fc981.
//
// Solidity: function getVoterLatestStake(address _candidate, address _voter) constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) GetVoterLatestStake(opts *bind.CallOpts, _candidate common.Address, _voter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getVoterLatestStake", _candidate, _voter)
	return *ret0, err
}

// GetVoterLatestStake is a free data retrieval call binding the contract method 0x432fc981.
//
// Solidity: function getVoterLatestStake(address _candidate, address _voter) constant returns(uint256)
func (_StakingContracts *StakingContractsSession) GetVoterLatestStake(_candidate common.Address, _voter common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterLatestStake(&_StakingContracts.CallOpts, _candidate, _voter)
}

// GetVoterLatestStake is a free data retrieval call binding the contract method 0x432fc981.
//
// Solidity: function getVoterLatestStake(address _candidate, address _voter) constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) GetVoterLatestStake(_candidate common.Address, _voter common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterLatestStake(&_StakingContracts.CallOpts, _candidate, _voter)
}

// GetVoterStake is a free data retrieval call binding the contract method 0x1fad0a71.
//
// Solidity: function getVoterStake(address _candidate, address _voter, uint256 _epoch) constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) GetVoterStake(opts *bind.CallOpts, _candidate common.Address, _voter common.Address, _epoch *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getVoterStake", _candidate, _voter, _epoch)
	return *ret0, err
}

// GetVoterStake is a free data retrieval call binding the contract method 0x1fad0a71.
//
// Solidity: function getVoterStake(address _candidate, address _voter, uint256 _epoch) constant returns(uint256)
func (_StakingContracts *StakingContractsSession) GetVoterStake(_candidate common.Address, _voter common.Address, _epoch *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStake(&_StakingContracts.CallOpts, _candidate, _voter, _epoch)
}

// GetVoterStake is a free data retrieval call binding the contract method 0x1fad0a71.
//
// Solidity: function getVoterStake(address _candidate, address _voter, uint256 _epoch) constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) GetVoterStake(_candidate common.Address, _voter common.Address, _epoch *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStake(&_StakingContracts.CallOpts, _candidate, _voter, _epoch)
}

// InitCandidates is a free data retrieval call binding the contract method 0x8106d590.
//
// Solidity: function initCandidates(uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsCaller) InitCandidates(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "initCandidates", arg0)
	return *ret0, err
}

// InitCandidates is a free data retrieval call binding the contract method 0x8106d590.
//
// Solidity: function initCandidates(uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsSession) InitCandidates(arg0 *big.Int) (common.Address, error) {
	return _StakingContracts.Contract.InitCandidates(&_StakingContracts.CallOpts, arg0)
}

// InitCandidates is a free data retrieval call binding the contract method 0x8106d590.
//
// Solidity: function initCandidates(uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsCallerSession) InitCandidates(arg0 *big.Int) (common.Address, error) {
	return _StakingContracts.Contract.InitCandidates(&_StakingContracts.CallOpts, arg0)
}

// MaxValidatorSize is a free data retrieval call binding the contract method 0x2de7dd5f.
//
// Solidity: function maxValidatorSize() constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) MaxValidatorSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "maxValidatorSize")
	return *ret0, err
}

// MaxValidatorSize is a free data retrieval call binding the contract method 0x2de7dd5f.
//
// Solidity: function maxValidatorSize() constant returns(uint256)
func (_StakingContracts *StakingContractsSession) MaxValidatorSize() (*big.Int, error) {
	return _StakingContracts.Contract.MaxValidatorSize(&_StakingContracts.CallOpts)
}

// MaxValidatorSize is a free data retrieval call binding the contract method 0x2de7dd5f.
//
// Solidity: function maxValidatorSize() constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) MaxValidatorSize() (*big.Int, error) {
	return _StakingContracts.Contract.MaxValidatorSize(&_StakingContracts.CallOpts)
}

// MinValidatorStake is a free data retrieval call binding the contract method 0x017ddd35.
//
// Solidity: function minValidatorStake() constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) MinValidatorStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "minValidatorStake")
	return *ret0, err
}

// MinValidatorStake is a free data retrieval call binding the contract method 0x017ddd35.
//
// Solidity: function minValidatorStake() constant returns(uint256)
func (_StakingContracts *StakingContractsSession) MinValidatorStake() (*big.Int, error) {
	return _StakingContracts.Contract.MinValidatorStake(&_StakingContracts.CallOpts)
}

// MinValidatorStake is a free data retrieval call binding the contract method 0x017ddd35.
//
// Solidity: function minValidatorStake() constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) MinValidatorStake() (*big.Int, error) {
	return _StakingContracts.Contract.MinValidatorStake(&_StakingContracts.CallOpts)
}

// MinVoterCap is a free data retrieval call binding the contract method 0xf8ac9dd5.
//
// Solidity: function minVoterCap() constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) MinVoterCap(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "minVoterCap")
	return *ret0, err
}

// MinVoterCap is a free data retrieval call binding the contract method 0xf8ac9dd5.
//
// Solidity: function minVoterCap() constant returns(uint256)
func (_StakingContracts *StakingContractsSession) MinVoterCap() (*big.Int, error) {
	return _StakingContracts.Contract.MinVoterCap(&_StakingContracts.CallOpts)
}

// MinVoterCap is a free data retrieval call binding the contract method 0xf8ac9dd5.
//
// Solidity: function minVoterCap() constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) MinVoterCap() (*big.Int, error) {
	return _StakingContracts.Contract.MinVoterCap(&_StakingContracts.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) StartBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "startBlock")
	return *ret0, err
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() constant returns(uint256)
func (_StakingContracts *StakingContractsSession) StartBlock() (*big.Int, error) {
	return _StakingContracts.Contract.StartBlock(&_StakingContracts.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) StartBlock() (*big.Int, error) {
	return _StakingContracts.Contract.StartBlock(&_StakingContracts.CallOpts)
}

// WithdrawalCap is a free data retrieval call binding the contract method 0x909b4053.
//
// Solidity: function withdrawalCap(address , uint256 ) constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) WithdrawalCap(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "withdrawalCap", arg0, arg1)
	return *ret0, err
}

// WithdrawalCap is a free data retrieval call binding the contract method 0x909b4053.
//
// Solidity: function withdrawalCap(address , uint256 ) constant returns(uint256)
func (_StakingContracts *StakingContractsSession) WithdrawalCap(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.WithdrawalCap(&_StakingContracts.CallOpts, arg0, arg1)
}

// WithdrawalCap is a free data retrieval call binding the contract method 0x909b4053.
//
// Solidity: function withdrawalCap(address , uint256 ) constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) WithdrawalCap(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.WithdrawalCap(&_StakingContracts.CallOpts, arg0, arg1)
}

// Register is a paid mutator transaction binding the contract method 0xaa677354.
//
// Solidity: function register(address _candidate, address _owner) returns()
func (_StakingContracts *StakingContractsTransactor) Register(opts *bind.TransactOpts, _candidate common.Address, _owner common.Address) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "register", _candidate, _owner)
}

// Register is a paid mutator transaction binding the contract method 0xaa677354.
//
// Solidity: function register(address _candidate, address _owner) returns()
func (_StakingContracts *StakingContractsSession) Register(_candidate common.Address, _owner common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Register(&_StakingContracts.TransactOpts, _candidate, _owner)
}

// Register is a paid mutator transaction binding the contract method 0xaa677354.
//
// Solidity: function register(address _candidate, address _owner) returns()
func (_StakingContracts *StakingContractsTransactorSession) Register(_candidate common.Address, _owner common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Register(&_StakingContracts.TransactOpts, _candidate, _owner)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(address _candidate) returns()
func (_StakingContracts *StakingContractsTransactor) Resign(opts *bind.TransactOpts, _candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "resign", _candidate)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(address _candidate) returns()
func (_StakingContracts *StakingContractsSession) Resign(_candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Resign(&_StakingContracts.TransactOpts, _candidate)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(address _candidate) returns()
func (_StakingContracts *StakingContractsTransactorSession) Resign(_candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Resign(&_StakingContracts.TransactOpts, _candidate)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_StakingContracts *StakingContractsTransactor) TransferAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "transferAdmin", newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_StakingContracts *StakingContractsSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.TransferAdmin(&_StakingContracts.TransactOpts, newAdmin)
}

// TransferAdmin is a paid mutator transaction binding the contract method 0x75829def.
//
// Solidity: function transferAdmin(address newAdmin) returns()
func (_StakingContracts *StakingContractsTransactorSession) TransferAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.TransferAdmin(&_StakingContracts.TransactOpts, newAdmin)
}

// Unvote is a paid mutator transaction binding the contract method 0x02aa9be2.
//
// Solidity: function unvote(address candidate, uint256 amount) returns()
func (_StakingContracts *StakingContractsTransactor) Unvote(opts *bind.TransactOpts, candidate common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "unvote", candidate, amount)
}

// Unvote is a paid mutator transaction binding the contract method 0x02aa9be2.
//
// Solidity: function unvote(address candidate, uint256 amount) returns()
func (_StakingContracts *StakingContractsSession) Unvote(candidate common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.Unvote(&_StakingContracts.TransactOpts, candidate, amount)
}

// Unvote is a paid mutator transaction binding the contract method 0x02aa9be2.
//
// Solidity: function unvote(address candidate, uint256 amount) returns()
func (_StakingContracts *StakingContractsTransactorSession) Unvote(candidate common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.Unvote(&_StakingContracts.TransactOpts, candidate, amount)
}

// UpdateMaxValidatorSize is a paid mutator transaction binding the contract method 0x0619624f.
//
// Solidity: function updateMaxValidatorSize(uint256 newMaxValidatorSize) returns()
func (_StakingContracts *StakingContractsTransactor) UpdateMaxValidatorSize(opts *bind.TransactOpts, newMaxValidatorSize *big.Int) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "updateMaxValidatorSize", newMaxValidatorSize)
}

// UpdateMaxValidatorSize is a paid mutator transaction binding the contract method 0x0619624f.
//
// Solidity: function updateMaxValidatorSize(uint256 newMaxValidatorSize) returns()
func (_StakingContracts *StakingContractsSession) UpdateMaxValidatorSize(newMaxValidatorSize *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.UpdateMaxValidatorSize(&_StakingContracts.TransactOpts, newMaxValidatorSize)
}

// UpdateMaxValidatorSize is a paid mutator transaction binding the contract method 0x0619624f.
//
// Solidity: function updateMaxValidatorSize(uint256 newMaxValidatorSize) returns()
func (_StakingContracts *StakingContractsTransactorSession) UpdateMaxValidatorSize(newMaxValidatorSize *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.UpdateMaxValidatorSize(&_StakingContracts.TransactOpts, newMaxValidatorSize)
}

// UpdateMinValidateStake is a paid mutator transaction binding the contract method 0xb2c76f10.
//
// Solidity: function updateMinValidateStake(uint256 _newCap) returns()
func (_StakingContracts *StakingContractsTransactor) UpdateMinValidateStake(opts *bind.TransactOpts, _newCap *big.Int) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "updateMinValidateStake", _newCap)
}

// UpdateMinValidateStake is a paid mutator transaction binding the contract method 0xb2c76f10.
//
// Solidity: function updateMinValidateStake(uint256 _newCap) returns()
func (_StakingContracts *StakingContractsSession) UpdateMinValidateStake(_newCap *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.UpdateMinValidateStake(&_StakingContracts.TransactOpts, _newCap)
}

// UpdateMinValidateStake is a paid mutator transaction binding the contract method 0xb2c76f10.
//
// Solidity: function updateMinValidateStake(uint256 _newCap) returns()
func (_StakingContracts *StakingContractsTransactorSession) UpdateMinValidateStake(_newCap *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.UpdateMinValidateStake(&_StakingContracts.TransactOpts, _newCap)
}

// UpdateMinVoteCap is a paid mutator transaction binding the contract method 0x3a1d8c5a.
//
// Solidity: function updateMinVoteCap(uint256 _newCap) returns()
func (_StakingContracts *StakingContractsTransactor) UpdateMinVoteCap(opts *bind.TransactOpts, _newCap *big.Int) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "updateMinVoteCap", _newCap)
}

// UpdateMinVoteCap is a paid mutator transaction binding the contract method 0x3a1d8c5a.
//
// Solidity: function updateMinVoteCap(uint256 _newCap) returns()
func (_StakingContracts *StakingContractsSession) UpdateMinVoteCap(_newCap *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.UpdateMinVoteCap(&_StakingContracts.TransactOpts, _newCap)
}

// UpdateMinVoteCap is a paid mutator transaction binding the contract method 0x3a1d8c5a.
//
// Solidity: function updateMinVoteCap(uint256 _newCap) returns()
func (_StakingContracts *StakingContractsTransactorSession) UpdateMinVoteCap(_newCap *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.UpdateMinVoteCap(&_StakingContracts.TransactOpts, _newCap)
}

// Vote is a paid mutator transaction binding the contract method 0x6dd7d8ea.
//
// Solidity: function vote(address candidate) returns()
func (_StakingContracts *StakingContractsTransactor) Vote(opts *bind.TransactOpts, candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "vote", candidate)
}

// Vote is a paid mutator transaction binding the contract method 0x6dd7d8ea.
//
// Solidity: function vote(address candidate) returns()
func (_StakingContracts *StakingContractsSession) Vote(candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Vote(&_StakingContracts.TransactOpts, candidate)
}

// Vote is a paid mutator transaction binding the contract method 0x6dd7d8ea.
//
// Solidity: function vote(address candidate) returns()
func (_StakingContracts *StakingContractsTransactorSession) Vote(candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Vote(&_StakingContracts.TransactOpts, candidate)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 epoch) returns(bool)
func (_StakingContracts *StakingContractsTransactor) Withdraw(opts *bind.TransactOpts, epoch *big.Int) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "withdraw", epoch)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 epoch) returns(bool)
func (_StakingContracts *StakingContractsSession) Withdraw(epoch *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.Withdraw(&_StakingContracts.TransactOpts, epoch)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 epoch) returns(bool)
func (_StakingContracts *StakingContractsTransactorSession) Withdraw(epoch *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.Withdraw(&_StakingContracts.TransactOpts, epoch)
}

// StakingContractsRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the StakingContracts contract.
type StakingContractsRegisteredIterator struct {
	Event *StakingContractsRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingContractsRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractsRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingContractsRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingContractsRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractsRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractsRegistered represents a Registered event raised by the StakingContracts contract.
type StakingContractsRegistered struct {
	Candidate common.Address
	Owner     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x0a31ee9d46a828884b81003c8498156ea6aa15b9b54bdd0ef0b533d9eba57e55.
//
// Solidity: event Registered(address candidate, address owner)
func (_StakingContracts *StakingContractsFilterer) FilterRegistered(opts *bind.FilterOpts) (*StakingContractsRegisteredIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &StakingContractsRegisteredIterator{contract: _StakingContracts.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x0a31ee9d46a828884b81003c8498156ea6aa15b9b54bdd0ef0b533d9eba57e55.
//
// Solidity: event Registered(address candidate, address owner)
func (_StakingContracts *StakingContractsFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *StakingContractsRegistered) (event.Subscription, error) {

	logs, sub, err := _StakingContracts.contract.WatchLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractsRegistered)
				if err := _StakingContracts.contract.UnpackLog(event, "Registered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakingContractsResignedIterator is returned from FilterResigned and is used to iterate over the raw logs and unpacked data for Resigned events raised by the StakingContracts contract.
type StakingContractsResignedIterator struct {
	Event *StakingContractsResigned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingContractsResignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractsResigned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingContractsResigned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingContractsResignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractsResignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractsResigned represents a Resigned event raised by the StakingContracts contract.
type StakingContractsResigned struct {
	Candidate common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterResigned is a free log retrieval operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address _candidate)
func (_StakingContracts *StakingContractsFilterer) FilterResigned(opts *bind.FilterOpts) (*StakingContractsResignedIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return &StakingContractsResignedIterator{contract: _StakingContracts.contract, event: "Resigned", logs: logs, sub: sub}, nil
}

// WatchResigned is a free log subscription operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address _candidate)
func (_StakingContracts *StakingContractsFilterer) WatchResigned(opts *bind.WatchOpts, sink chan<- *StakingContractsResigned) (event.Subscription, error) {

	logs, sub, err := _StakingContracts.contract.WatchLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractsResigned)
				if err := _StakingContracts.contract.UnpackLog(event, "Resigned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakingContractsUnvotedIterator is returned from FilterUnvoted and is used to iterate over the raw logs and unpacked data for Unvoted events raised by the StakingContracts contract.
type StakingContractsUnvotedIterator struct {
	Event *StakingContractsUnvoted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingContractsUnvotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractsUnvoted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingContractsUnvoted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingContractsUnvotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractsUnvotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractsUnvoted represents a Unvoted event raised by the StakingContracts contract.
type StakingContractsUnvoted struct {
	Voter     common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnvoted is a free log retrieval operation binding the contract event 0x7958395da8e26969cc7c09ee58e9507a2601574c3bd232617e2d6354224ff836.
//
// Solidity: event Unvoted(address voter, address candidate, uint256 amount)
func (_StakingContracts *StakingContractsFilterer) FilterUnvoted(opts *bind.FilterOpts) (*StakingContractsUnvotedIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Unvoted")
	if err != nil {
		return nil, err
	}
	return &StakingContractsUnvotedIterator{contract: _StakingContracts.contract, event: "Unvoted", logs: logs, sub: sub}, nil
}

// WatchUnvoted is a free log subscription operation binding the contract event 0x7958395da8e26969cc7c09ee58e9507a2601574c3bd232617e2d6354224ff836.
//
// Solidity: event Unvoted(address voter, address candidate, uint256 amount)
func (_StakingContracts *StakingContractsFilterer) WatchUnvoted(opts *bind.WatchOpts, sink chan<- *StakingContractsUnvoted) (event.Subscription, error) {

	logs, sub, err := _StakingContracts.contract.WatchLogs(opts, "Unvoted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractsUnvoted)
				if err := _StakingContracts.contract.UnpackLog(event, "Unvoted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakingContractsVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the StakingContracts contract.
type StakingContractsVotedIterator struct {
	Event *StakingContractsVoted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingContractsVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractsVoted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingContractsVoted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingContractsVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractsVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractsVoted represents a Voted event raised by the StakingContracts contract.
type StakingContractsVoted struct {
	Voter     common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0x174ba19ba3c8bb5c679c87e51db79fff7c3f04bb84c1fd55b7dacb470b674aa6.
//
// Solidity: event Voted(address voter, address candidate, uint256 amount)
func (_StakingContracts *StakingContractsFilterer) FilterVoted(opts *bind.FilterOpts) (*StakingContractsVotedIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Voted")
	if err != nil {
		return nil, err
	}
	return &StakingContractsVotedIterator{contract: _StakingContracts.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0x174ba19ba3c8bb5c679c87e51db79fff7c3f04bb84c1fd55b7dacb470b674aa6.
//
// Solidity: event Voted(address voter, address candidate, uint256 amount)
func (_StakingContracts *StakingContractsFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *StakingContractsVoted) (event.Subscription, error) {

	logs, sub, err := _StakingContracts.contract.WatchLogs(opts, "Voted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractsVoted)
				if err := _StakingContracts.contract.UnpackLog(event, "Voted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakingContractsWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the StakingContracts contract.
type StakingContractsWithdrawIterator struct {
	Event *StakingContractsWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingContractsWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingContractsWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingContractsWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingContractsWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingContractsWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingContractsWithdraw represents a Withdraw event raised by the StakingContracts contract.
type StakingContractsWithdraw struct {
	Staker common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address _staker, uint256 _amount)
func (_StakingContracts *StakingContractsFilterer) FilterWithdraw(opts *bind.FilterOpts) (*StakingContractsWithdrawIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &StakingContractsWithdrawIterator{contract: _StakingContracts.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address _staker, uint256 _amount)
func (_StakingContracts *StakingContractsFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StakingContractsWithdraw) (event.Subscription, error) {

	logs, sub, err := _StakingContracts.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingContractsWithdraw)
				if err := _StakingContracts.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
