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
const StakingContractsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"minValidatorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unvote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMaxValidatorSize\",\"type\":\"uint256\"}],\"name\":\"updateMaxValidatorSize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getVoterStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"getCandidateData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_isCandidate\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_totalStake\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxValidatorSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAllCandidates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"candidates\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newCap\",\"type\":\"uint256\"}],\"name\":\"updateMinVoteCap\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"startBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochSize\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getListCandidates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_stakes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint32\",\"name\":\"_maxValSize\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"_epochSize\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"vote\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"candidateData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isCandidate\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"resign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newCap\",\"type\":\"uint256\"}],\"name\":\"updateMinValidateStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minVoterCap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"candidatesOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxValidatorSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minValidatorStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minVoteCap\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unvoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"}]"

// StakingContractsBin is the compiled bytecode used for deploying new contracts.
const StakingContractsBin = `608060405260016000553480156200001657600080fd5b50604051620022d3380380620022d3833981810160405260c08110156200003c57600080fd5b81019080805160405193929190846401000000008211156200005d57600080fd5b838201915060208201858111156200007457600080fd5b82518660208202830111640100000000821117156200009257600080fd5b8083526020830192505050908051906020019060200280838360005b83811015620000cb578082015181840152602081019050620000ae565b505050509050016040526020018051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050600084116200011b57600080fd5b8360048190555082600581905550816006819055508060078190555085600290805190602001906200014f92919062000356565b5060008090505b8651811015620003015760405180606001604052806001151581526020018481526020018773ffffffffffffffffffffffffffffffffffffffff1681525060016000898481518110620001a557fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548160ff0219169083151502179055506020820151816001015560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505082600160008984815181106200026d57fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550808060010191505062000156565b5033600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550436003819055505050505050506200042b565b828054828255906000526020600020908101928215620003d2579160200282015b82811115620003d15782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509160200191906001019062000377565b5b509050620003e19190620003e5565b5090565b6200042891905b808211156200042457600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905550600101620003ec565b5090565b90565b611e98806200043b6000396000f3fe6080604052600436106101355760003560e01c8063572d356e116100ab578063ae6e43f51161006f578063ae6e43f514610725578063b2c76f1014610776578063b5b7a184146107b1578063b97dd9e2146107dc578063f851a44014610807578063f8ac9dd51461085e57610135565b8063572d356e146104dc578063690ff8a1146105135780636dd7d8ea146105ed57806375829def1461063157806391a9634f1461068257610135565b80632de7dd5f116100fd5780632de7dd5f146103205780632e6997fe1461034b5780633477ee2e146103b75780633a1d8c5a146104325780634420e4861461046d57806348cd4cb1146104b157610135565b8063017ddd351461013757806302aa9be2146101625780630619624f146101bd578063158a65f6146101f85780632a466ac71461027d575b005b34801561014357600080fd5b5061014c610889565b6040518082815260200191505060405180910390f35b34801561016e57600080fd5b506101bb6004803603604081101561018557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061088f565b005b3480156101c957600080fd5b506101f6600480360360208110156101e057600080fd5b8101908080359060200190929190505050610b07565b005b34801561020457600080fd5b506102676004803603604081101561021b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b6b565b6040518082815260200191505060405180910390f35b34801561028957600080fd5b506102cc600480360360208110156102a057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610bf5565b60405180841515151581526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390f35b34801561032c57600080fd5b50610335610cfd565b6040518082815260200191505060405180910390f35b34801561035757600080fd5b50610360610d03565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156103a3578082015181840152602081019050610388565b505050509050019250505060405180910390f35b3480156103c357600080fd5b506103f0600480360360208110156103da57600080fd5b8101908080359060200190929190505050610d91565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561043e57600080fd5b5061046b6004803603602081101561045557600080fd5b8101908080359060200190929190505050610dcd565b005b6104af6004803603602081101561048357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610e31565b005b3480156104bd57600080fd5b506104c66110d0565b6040518082815260200191505060405180910390f35b3480156104e857600080fd5b506104f16110d6565b604051808263ffffffff1663ffffffff16815260200191505060405180910390f35b34801561051f57600080fd5b506105286110e5565b6040518080602001806020018563ffffffff1663ffffffff1681526020018463ffffffff1663ffffffff168152602001838103835287818151815260200191508051906020019060200280838360005b83811015610593578082015181840152602081019050610578565b50505050905001838103825286818151815260200191508051906020019060200280838360005b838110156105d55780820151818401526020810190506105ba565b50505050905001965050505050505060405180910390f35b61062f6004803603602081101561060357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061158d565b005b34801561063d57600080fd5b506106806004803603602081101561065457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061177d565b005b34801561068e57600080fd5b506106d1600480360360208110156106a557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611855565b60405180841515151581526020018381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001935050505060405180910390f35b34801561073157600080fd5b506107746004803603602081101561074857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506118ac565b005b34801561078257600080fd5b506107af6004803603602081101561079957600080fd5b8101908080359060200190929190505050611d65565b005b3480156107bd57600080fd5b506107c6611dc9565b6040518082815260200191505060405180910390f35b3480156107e857600080fd5b506107f1611dcf565b6040518082815260200191505060405180910390f35b34801561081357600080fd5b5061081c611de6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561086a57600080fd5b50610873611e0c565b6040518082815260200191505060405180910390f35b60065481565b6001600080828254019250508190555060008054905081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561093157600080fd5b81600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101600082825403925050819055503373ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015610a54573d6000803e3d6000fd5b507f7958395da8e26969cc7c09ee58e9507a2601574c3bd232617e2d6354224ff836338484604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a16000548114610b0257600080fd5b505050565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610b6157600080fd5b8060058190555050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000806000600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff169250600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169150600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001015490509193909250565b60055481565b60606002805480602002602001604051908101604052809291908181526020018280548015610d8757602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610d3d575b5050505050905090565b60028181548110610d9e57fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610e2757600080fd5b8060078190555050565b60001515600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff16151514610e9157600080fd5b60405180606001604052806001151581526020013481526020013373ffffffffffffffffffffffffffffffffffffffff16815250600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548160ff0219169083151502179055506020820151816001015560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505034600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060028190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550507f6f3bf3fa84e4763a43b3d23f9d79be242d6d5c834941ff4c1111b67469e1150c8134604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a150565b60035481565b60006110e0611dcf565b905090565b60608060008060055491506110f8611dcf565b9050600080600090505b600280549050811015611283576000600160006002848154811061112257fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060065460016000600285815481106111c157fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054106112755782806001019350505b508080600101915050611102565b50806040519080825280602002602001820160405280156112b35781602001602082028038833980820191505090505b509450806040519080825280602002602001820160405280156112e55781602001602082028038833980820191505090505b509350600080905060008090505b600280549050811015611584576000600160006002848154811061131357fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060065460016000600285815481106113b257fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410611576576002828154811061146a57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168884815181106114a157fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505060016000600284815481106114ec57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001015487848151811061156157fe5b60200260200101818152505082806001019350505b5080806001019150506112f3565b50505090919293565b60075434101561159c57600080fd5b8060011515600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff161515146115fd57600080fd5b34600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001016000828254019250508190555034600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055507f174ba19ba3c8bb5c679c87e51db79fff7c3f04bb84c1fd55b7dacb470b674aa6338334604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15050565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146117d757600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561181157600080fd5b80600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60016020528060005260406000206000915090508060000160009054906101000a900460ff16908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b6001600080828254019250508190555060008054905060011515600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900460ff1615151461192257600080fd5b3373ffffffffffffffffffffffffffffffffffffffff16600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146119bc57600080fd5b60008090505b600280549050811015611b39578273ffffffffffffffffffffffffffffffffffffffff16600282815481106119f357fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415611b2c57600260016002805490500381548110611a4f57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660028281548110611a8757fe5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600260016002805490500381548110611ae457fe5b9060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556002805480919060019003611b269190611e12565b50611b39565b80806001019150506119c2565b506000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a81548160ff0219169083151502179055506000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490506000600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000811115611cef573373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015611ced573d6000803e3d6000fd5b505b7fa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d8683604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1506000548114611d6157600080fd5b5050565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611dbf57600080fd5b8060068190555050565b60045481565b6000600454600354430381611de057fe5b04905090565b600860009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60075481565b815481835581811115611e3957818360005260206000209182019101611e389190611e3e565b5b505050565b611e6091905b80821115611e5c576000816000905550600101611e44565b5090565b9056fea265627a7a72315820e3f5a7a6b2d8df9a945afcca908abf5f9ba334eb7d59397ff8ae9200767ad4ba64736f6c634300050b0032`

// DeployStakingContracts deploys a new Evrynet contract, binding an instance of StakingContracts to it.
func DeployStakingContracts(auth *bind.TransactOpts, backend bind.ContractBackend, _candidates []common.Address, candidatesOwner common.Address, _epochPeriod *big.Int, _maxValidatorSize *big.Int, _minValidatorStake *big.Int, _minVoteCap *big.Int) (common.Address, *types.Transaction, *StakingContracts, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingContractsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StakingContractsBin), backend, _candidates, candidatesOwner, _epochPeriod, _maxValidatorSize, _minValidatorStake, _minVoteCap)
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
// Solidity: function candidateData(address ) constant returns(bool isCandidate, uint256 totalStake, address owner)
func (_StakingContracts *StakingContractsCaller) CandidateData(opts *bind.CallOpts, arg0 common.Address) (struct {
	IsCandidate bool
	TotalStake  *big.Int
	Owner       common.Address
}, error) {
	ret := new(struct {
		IsCandidate bool
		TotalStake  *big.Int
		Owner       common.Address
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "candidateData", arg0)
	return *ret, err
}

// CandidateData is a free data retrieval call binding the contract method 0x91a9634f.
//
// Solidity: function candidateData(address ) constant returns(bool isCandidate, uint256 totalStake, address owner)
func (_StakingContracts *StakingContractsSession) CandidateData(arg0 common.Address) (struct {
	IsCandidate bool
	TotalStake  *big.Int
	Owner       common.Address
}, error) {
	return _StakingContracts.Contract.CandidateData(&_StakingContracts.CallOpts, arg0)
}

// CandidateData is a free data retrieval call binding the contract method 0x91a9634f.
//
// Solidity: function candidateData(address ) constant returns(bool isCandidate, uint256 totalStake, address owner)
func (_StakingContracts *StakingContractsCallerSession) CandidateData(arg0 common.Address) (struct {
	IsCandidate bool
	TotalStake  *big.Int
	Owner       common.Address
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
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isCandidate, address _owner, uint256 _totalStake)
func (_StakingContracts *StakingContractsCaller) GetCandidateData(opts *bind.CallOpts, _candidate common.Address) (struct {
	IsCandidate bool
	Owner       common.Address
	TotalStake  *big.Int
}, error) {
	ret := new(struct {
		IsCandidate bool
		Owner       common.Address
		TotalStake  *big.Int
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "getCandidateData", _candidate)
	return *ret, err
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isCandidate, address _owner, uint256 _totalStake)
func (_StakingContracts *StakingContractsSession) GetCandidateData(_candidate common.Address) (struct {
	IsCandidate bool
	Owner       common.Address
	TotalStake  *big.Int
}, error) {
	return _StakingContracts.Contract.GetCandidateData(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isCandidate, address _owner, uint256 _totalStake)
func (_StakingContracts *StakingContractsCallerSession) GetCandidateData(_candidate common.Address) (struct {
	IsCandidate bool
	Owner       common.Address
	TotalStake  *big.Int
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

// GetVoterStake is a free data retrieval call binding the contract method 0x158a65f6.
//
// Solidity: function getVoterStake(address _candidate, address _voter) constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) GetVoterStake(opts *bind.CallOpts, _candidate common.Address, _voter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getVoterStake", _candidate, _voter)
	return *ret0, err
}

// GetVoterStake is a free data retrieval call binding the contract method 0x158a65f6.
//
// Solidity: function getVoterStake(address _candidate, address _voter) constant returns(uint256)
func (_StakingContracts *StakingContractsSession) GetVoterStake(_candidate common.Address, _voter common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStake(&_StakingContracts.CallOpts, _candidate, _voter)
}

// GetVoterStake is a free data retrieval call binding the contract method 0x158a65f6.
//
// Solidity: function getVoterStake(address _candidate, address _voter) constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) GetVoterStake(_candidate common.Address, _voter common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStake(&_StakingContracts.CallOpts, _candidate, _voter)
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

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address _candidate) returns()
func (_StakingContracts *StakingContractsTransactor) Register(opts *bind.TransactOpts, _candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "register", _candidate)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address _candidate) returns()
func (_StakingContracts *StakingContractsSession) Register(_candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Register(&_StakingContracts.TransactOpts, _candidate)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address _candidate) returns()
func (_StakingContracts *StakingContractsTransactorSession) Register(_candidate common.Address) (*types.Transaction, error) {
	return _StakingContracts.Contract.Register(&_StakingContracts.TransactOpts, _candidate)
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
	Stake     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x6f3bf3fa84e4763a43b3d23f9d79be242d6d5c834941ff4c1111b67469e1150c.
//
// Solidity: event Registered(address candidate, uint256 stake)
func (_StakingContracts *StakingContractsFilterer) FilterRegistered(opts *bind.FilterOpts) (*StakingContractsRegisteredIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &StakingContractsRegisteredIterator{contract: _StakingContracts.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x6f3bf3fa84e4763a43b3d23f9d79be242d6d5c834941ff4c1111b67469e1150c.
//
// Solidity: event Registered(address candidate, uint256 stake)
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
