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
const StakingContractsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"minValidatorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unvote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMaxValidatorSize\",\"type\":\"uint256\"}],\"name\":\"updateMaxValidatorSize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getWithdrawEpochs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"epochs\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getVoterStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stake\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"getWithdrawCap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"cap\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"getCandidateData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_isActiveCandidate\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_totalStake\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"voters\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxValidatorSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"candidates\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newCap\",\"type\":\"uint256\"}],\"name\":\"updateMinVoteCap\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"getCandidateStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"startBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getListCandidates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minValidatorCap\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"vote\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"candidateVoters\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"candidateData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isCandidate\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"totalStake\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"withdrawWithIndex\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"resign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newCap\",\"type\":\"uint256\"}],\"name\":\"updateMinValidateStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"getCandidateOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"}],\"name\":\"isCandidate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getWithdrawEpochsAndCaps\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"epochs\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"caps\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"voters\",\"type\":\"address[]\"}],\"name\":\"getVoterStakes\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"stakes\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minVoterCap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_candidates\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"candidateOwners\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxValidatorSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minValidatorStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minVoteCap\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unvoted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_candidate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"}]"

// StakingContractsBin is the compiled bytecode used for deploying new contracts.
const StakingContractsBin = `608060405260016000553480156200001657600080fd5b50604051620023463803806200234683398181016040526101008110156200003d57600080fd5b81019080805160405193929190846401000000008211156200005e57600080fd5b9083019060208201858111156200007457600080fd5b82518660208202830111640100000000821117156200009257600080fd5b82525081516020918201928201910280838360005b83811015620000c1578181015183820152602001620000a7565b5050505090500160405260200180516040519392919084640100000000821115620000eb57600080fd5b9083019060208201858111156200010157600080fd5b82518660208202830111640100000000821117156200011f57600080fd5b82525081516020918201928201910280838360005b838110156200014e57818101518382015260200162000134565b505050509190910160409081526020830151908301516060840151608085015160a086015160c090960151939750919550939092509085620001f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f65706f6368206d75737420626520706f73697469766500000000000000000000604482015290519081900360640190fd5b86518851146200026257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6c656e677468206e6f74206d6174636800000000000000000000000000000000604482015290519081900360640190fd5b600686905560078490556008839055600982905587518410156200028557600080fd5b87516200029a9060049060208b019062000489565b5060005b885181101562000458576040518060600160405280600115158152602001858152602001898381518110620002cf57fe5b60200260200101516001600160a01b0316815250600360008b8481518110620002f457fe5b6020908102919091018101516001600160a01b03908116835282820193909352604091820160009081208551815460ff19169015151781559185015160018301559390910151600290910180546001600160a01b03191691909216179055895185916003918c90859081106200036657fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060030160008a8481518110620003a057fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002081905550600260008a8381518110620003dd57fe5b60200260200101516001600160a01b03166001600160a01b031681526020019081526020016000208882815181106200041257fe5b60209081029190910181015182546001808201855560009485529290932090920180546001600160a01b0319166001600160a01b0390931692909217909155016200029e565b50600a80546001600160a01b0319166001600160a01b0392909216919091179055505050600555506200051d915050565b828054828255906000526020600020908101928215620004e1579160200282015b82811115620004e157825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620004aa565b50620004ef929150620004f3565b5090565b6200051a91905b80821115620004ef5780546001600160a01b0319168155600101620004fa565b90565b611e19806200052d6000396000f3fe6080604052600436106101d85760003560e01c80636dd7d8ea11610102578063b5b7a18411610095578063d5816bfa11610064578063d5816bfa146107d0578063e2db89b51461087e578063f851a4401461093e578063f8ac9dd514610953576101d8565b8063b5b7a18414610740578063b642facd14610755578063b97dd9e214610788578063d51b9e931461079d576101d8565b806396c23442116100d157806396c2344214610678578063aa677354146106a8578063ae6e43f5146106e3578063b2c76f1014610716576101d8565b80636dd7d8ea1461058a57806375829def146105b05780637b728966146105e357806391a9634f1461061c576101d8565b80632d15cc041161017a5780633a1d8c5a116101495780633a1d8c5a14610455578063484da9611461047f57806348cd4cb1146104b2578063690ff8a1146104c7576101d8565b80632d15cc04146103895780632de7dd5f146103bc5780632e1a7d4d146103d15780633477ee2e1461040f576101d8565b80630e0516aa116101b65780630e0516aa14610264578063158a65f6146102c957806315febd68146103045780632a466ac71461032e576101d8565b8063017ddd35146101da57806302aa9be2146102015780630619624f1461023a575b005b3480156101e657600080fd5b506101ef610968565b60408051918252519081900360200190f35b34801561020d57600080fd5b506101d86004803603604081101561022457600080fd5b506001600160a01b03813516906020013561096e565b34801561024657600080fd5b506101d86004803603602081101561025d57600080fd5b5035610cc4565b34801561027057600080fd5b50610279610ce0565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156102b557818101518382015260200161029d565b505050509050019250505060405180910390f35b3480156102d557600080fd5b506101ef600480360360408110156102ec57600080fd5b506001600160a01b0381358116916020013516610d44565b34801561031057600080fd5b506101ef6004803603602081101561032757600080fd5b5035610d70565b34801561033a57600080fd5b506103616004803603602081101561035157600080fd5b50356001600160a01b0316610d8d565b6040805193151584526001600160a01b03909216602084015282820152519081900360600190f35b34801561039557600080fd5b50610279600480360360208110156103ac57600080fd5b50356001600160a01b0316610dc1565b3480156103c857600080fd5b506101ef610e37565b3480156103dd57600080fd5b506103fb600480360360208110156103f457600080fd5b5035610e3d565b604080519115158252519081900360200190f35b34801561041b57600080fd5b506104396004803603602081101561043257600080fd5b5035610f48565b604080516001600160a01b039092168252519081900360200190f35b34801561046157600080fd5b506101d86004803603602081101561047857600080fd5b5035610f6f565b34801561048b57600080fd5b506101ef600480360360208110156104a257600080fd5b50356001600160a01b0316610f8b565b3480156104be57600080fd5b506101ef610fa9565b3480156104d357600080fd5b506104dc610faf565b604051808060200180602001868152602001858152602001848152602001838103835288818151815260200191508051906020019060200280838360005b8381101561053257818101518382015260200161051a565b50505050905001838103825287818151815260200191508051906020019060200280838360005b83811015610571578181015183820152602001610559565b5050505090500197505050505050505060405180910390f35b6101d8600480360360208110156105a057600080fd5b50356001600160a01b03166110c5565b3480156105bc57600080fd5b506101d8600480360360208110156105d357600080fd5b50356001600160a01b031661128d565b3480156105ef57600080fd5b506104396004803603604081101561060657600080fd5b506001600160a01b0381351690602001356112d9565b34801561062857600080fd5b5061064f6004803603602081101561063f57600080fd5b50356001600160a01b031661130e565b60408051931515845260208401929092526001600160a01b031682820152519081900360600190f35b34801561068457600080fd5b506103fb6004803603604081101561069b57600080fd5b508035906020013561133c565b3480156106b457600080fd5b506101d8600480360360408110156106cb57600080fd5b506001600160a01b03813581169160200135166115b2565b3480156106ef57600080fd5b506101d86004803603602081101561070657600080fd5b50356001600160a01b0316611809565b34801561072257600080fd5b506101d86004803603602081101561073957600080fd5b5035611ac9565b34801561074c57600080fd5b506101ef611ae5565b34801561076157600080fd5b506104396004803603602081101561077857600080fd5b50356001600160a01b0316611aeb565b34801561079457600080fd5b506101ef611b0c565b3480156107a957600080fd5b506103fb600480360360208110156107c057600080fd5b50356001600160a01b0316611b3b565b3480156107dc57600080fd5b506107e5611b59565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b83811015610829578181015183820152602001610811565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610868578181015183820152602001610850565b5050505090500194505050505060405180910390f35b34801561088a57600080fd5b50610279600480360360408110156108a157600080fd5b6001600160a01b0382351691908101906040810160208201356401000000008111156108cc57600080fd5b8201836020820111156108de57600080fd5b8035906020019184602083028401116401000000008311171561090057600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550611c54945050505050565b34801561094a57600080fd5b50610439611d0e565b34801561095f57600080fd5b506101ef611d1d565b60085481565b60008054600101908190558282806109cd576040805162461bcd60e51b815260206004820152601760248201527f5f6361702073686f756c6420626520706f736974697665000000000000000000604482015290519081900360640190fd5b6001600160a01b03821660009081526003602081815260408084203380865293019091529091205482811015610a41576040805162461bcd60e51b81526020600482015260146024820152736e6f7420656e6f75676820746f20756e766f746560601b604482015290519081900360640190fd5b6001600160a01b0384811660009081526003602052604090206002015481169083161415610acb57600854610a7c828563ffffffff611d2316565b1015610ac6576040805162461bcd60e51b81526020600482015260146024820152736e6f7420656e6f75676820746f20756e766f746560601b604482015290519081900360640190fd5b610b36565b6000610add828563ffffffff611d2316565b9050801580610aee57506009548110155b610b34576040805162461bcd60e51b81526020600482015260126024820152711a5b9d985b1a59081d5b9d9bdd1948185b5d60721b604482015290519081900360640190fd5b505b6000610b40611b0c565b6001600160a01b03891660009081526003602081815260408084203380865293019091529091205491925090610b7c908963ffffffff611d2316565b6001600160a01b03808b1660008181526003602081815260408084209588168452858301825283209590955591905290915260010154610bc2908963ffffffff611d2316565b6001600160a01b038a16600090815260036020526040812060010191909155610bf283600263ffffffff611d3816565b6001600160a01b0383166000908152600160209081526040808320848452909152902054909150610c29908a63ffffffff611d3816565b6001600160a01b0380841660008181526001602081815260408084208885528083528185209790975582825295820180549283018155835291829020018590558351918252918d16918101919091528082018b905290517f7958395da8e26969cc7c09ee58e9507a2601574c3bd232617e2d6354224ff8369181900360600190a1505050505050506000548114610cbf57600080fd5b505050565b600a546001600160a01b03163314610cdb57600080fd5b600755565b33600090815260016020818152604092839020909101805483518184028101840190945280845260609392830182828015610d3a57602002820191906000526020600020905b815481526020019060010190808311610d26575b5050505050905090565b6001600160a01b0391821660009081526003602081815260408084209490951683529201909152205490565b336000908152600160209081526040808320938352929052205490565b6001600160a01b0390811660009081526003602052604090208054600282015460019092015460ff90911693919092169190565b6001600160a01b038116600090815260026020908152604091829020805483518184028101840190945280845260609392830182828015610e2b57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610e0d575b50505050509050919050565b60075481565b6000805460010180825581610e50611b0c565b905083811015610e915760405162461bcd60e51b8152600401808060200182810382526021815260200180611dc46021913960400191505060405180910390fd5b3360008181526001602090815260408083208884529091528120805491905580610ef6576040805162461bcd60e51b81526020600482015260116024820152700776974686472617720636170206973203607c1b604482015290519081900360640190fd5b6040516001600160a01b0383169082156108fc029083906000818181858888f19350505050158015610f2c573d6000803e3d6000fd5b50600194505050506000548114610f4257600080fd5b50919050565b60048181548110610f5557fe5b6000918252602090912001546001600160a01b0316905081565b600a546001600160a01b03163314610f8657600080fd5b600955565b6001600160a01b031660009081526003602052604090206001015490565b60055481565b6060806000806000610fbf611b0c565b925060075491506008549050600480548060200260200160405190810160405280929190818152602001828054801561102157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611003575b505050505094508451604051908082528060200260200182016040528015611053578160200160208202803883390190505b50935060005b85518110156110bd576003600087838151811061107257fe5b60200260200101516001600160a01b03166001600160a01b03168152602001908152602001600020600101548582815181106110aa57fe5b6020908102919091010152600101611059565b509091929394565b60095434101561110d576040805162461bcd60e51b815260206004820152600e60248201526d1b1bddc81d9bdd1948185b5bdd5d60921b604482015290519081900360640190fd5b6001600160a01b038116600090815260036020526040902054819060ff16151560011461113957600080fd5b6001600160a01b0382166000908152600360208181526040808420338086529301909152909120543491906111a7576001600160a01b0384811660009081526002602090815260408220805460018101825590835291200180546001600160a01b0319169183169190911790555b6001600160a01b038085166000908152600360208181526040808420948616845293909101905220546111e0908363ffffffff611d3816565b6001600160a01b0380861660008181526003602081815260408084209588168452858301825283209590955591905290915260010154611226908363ffffffff611d3816565b6001600160a01b0380861660008181526003602090815260409182902060010194909455805192851683529282015280820184905290517f174ba19ba3c8bb5c679c87e51db79fff7c3f04bb84c1fd55b7dacb470b674aa69181900360600190a150505050565b600a546001600160a01b031633146112a457600080fd5b6001600160a01b0381166112b757600080fd5b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b600260205281600052604060002081815481106112f257fe5b6000918252602090912001546001600160a01b03169150829050565b60036020526000908152604090208054600182015460029092015460ff90911691906001600160a01b031683565b600080546001018082558161134f611b0c565b9050848110156113905760405162461bcd60e51b8152600401808060200182810382526021815260200180611dc46021913960400191505060405180910390fd5b336000818152600160208190526040909120018054879190879081106113b257fe5b906000526020600020015414611403576040805162461bcd60e51b81526020600482015260116024820152700dcdee840c6dee4e4cac6e840d2dcc8caf607b1b604482015290519081900360640190fd5b6001600160a01b03811660009081526001602090815260408083208984529091529020548061146d576040805162461bcd60e51b81526020600482015260116024820152700776974686472617720636170206973203607c1b604482015290519081900360640190fd5b6001600160a01b03821660008181526001602081815260408084208c855280835290842084905593909252908190520180549060001982018281106114ae57fe5b906000526020600020015460016000856001600160a01b03166001600160a01b0316815260200190815260200160002060010188815481106114ec57fe5b60009182526020808320909101929092556001600160a01b038516815260019182905260409020018054600019830190811061152457fe5b600091825260208083209091018290556001600160a01b038516825260019081905260409091200180549061155d906000198301611d86565b506040516001600160a01b0384169083156108fc029084906000818181858888f19350505050158015611594573d6000803e3d6000fd5b50600195505050505060005481146115ab57600080fd5b5092915050565b600a546001600160a01b031633146115c957600080fd5b6001600160a01b038216600090815260036020526040902054829060ff16156115f157600080fd5b6001600160a01b03831661164c576040805162461bcd60e51b815260206004820152601d60248201527f5f63616e6469646174652061646472657373206973206d697373696e67000000604482015290519081900360640190fd5b6001600160a01b0382166116a7576040805162461bcd60e51b815260206004820152601960248201527f5f6f776e65722061646472657373206973206d697373696e6700000000000000604482015290519081900360640190fd5b6001600160a01b038316600090815260036020526040902060010154600454608011611710576040805162461bcd60e51b8152602060048201526013602482015272746f6f206d616e792063616e6469646174657360681b604482015290519081900360640190fd5b60408051606081018252600180825260208083018581526001600160a01b038881168587018181528b83166000818152600387528981209851895460ff19169015151789559451888801559051600297880180546001600160a01b031990811692909516919091179055600480548088019091557f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b0180548416821790559584528683208054958601815583529183902090930180549093168117909255835192835282015281517f0a31ee9d46a828884b81003c8498156ea6aa15b9b54bdd0ef0b533d9eba57e55929181900390910190a150505050565b6001600160a01b038116600090815260036020526040902054819060ff16151560011461183557600080fd5b6001600160a01b038083166000908152600360205260409020600201548391163314611894576040805162461bcd60e51b81526020600482015260096024820152683737ba1037bbb732b960b91b604482015290519081900360640190fd5b33600061189f611b0c565b905060005b60045481101561198d57856001600160a01b0316600482815481106118c557fe5b6000918252602090912001546001600160a01b03161415611985576004805460001981019081106118f257fe5b600091825260209091200154600480546001600160a01b03909216918390811061191857fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560048054600019810190811061195357fe5b600091825260209091200180546001600160a01b0319169055600480549061197f906000198301611d86565b5061198d565b6001016118a4565b506001600160a01b038086166000818152600360208181526040808420805460ff191681559588168452858301825283208054908490559390925290526001909101546119e0908263ffffffff611d2316565b6001600160a01b038716600090815260036020526040812060010191909155611a1083600263ffffffff611d3816565b6001600160a01b0385166000908152600160209081526040808320848452909152902054909150611a47908363ffffffff611d3816565b6001600160a01b0380861660009081526001602081815260408084208785528083528185209690965582825294820180549283018155835291829020018490558251918a168252810185905281517f886e0db046874dde595498040d176412e81183750ceb33fc46f0450362bbc241929181900390910190a150505050505050565b600a546001600160a01b03163314611ae057600080fd5b600855565b60065481565b6001600160a01b039081166000908152600360205260409020600201541690565b6000611b35600654611b2960055443611d2390919063ffffffff16565b9063ffffffff611d5116565b90505b90565b6001600160a01b031660009081526003602052604090205460ff1690565b336000908152600160208181526040928390209091018054835181840281018401909452808452606093849390929190830182828015611bb857602002820191906000526020600020905b815481526020019060010190808311611ba4575b505050505091508151604051908082528060200260200182016040528015611bea578160200160208202803883390190505b50905060005b8251811015611c4f573360009081526001602052604081208451909190859084908110611c1957fe5b6020026020010151815260200190815260200160002054828281518110611c3c57fe5b6020908102919091010152600101611bf0565b509091565b60608151604051908082528060200260200182016040528015611c81578160200160208202803883390190505b50905060005b82518110156115ab5760036000856001600160a01b03166001600160a01b031681526020019081526020016000206003016000848381518110611cc657fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002054828281518110611cfb57fe5b6020908102919091010152600101611c87565b600a546001600160a01b031681565b60095481565b600082821115611d3257600080fd5b50900390565b600082820183811015611d4a57600080fd5b9392505050565b6000808211611d5f57600080fd5b6000828481611d6a57fe5b049050828481611d7657fe5b06818402018414611d4a57600080fd5b815481835581811115610cbf57600083815260209020610cbf918101908301611b3891905b80821115611dbf5760008155600101611dab565b509056fe63616e206e6f7420776974686472617720666f72206675747572652065706f6368a265627a7a723158206dfa26126052a83dc61aaa1e216b16211108477ad238da8a78cf339cc89e8bf764736f6c634300050b0032`

// DeployStakingContracts deploys a new Evrynet contract, binding an instance of StakingContracts to it.
func DeployStakingContracts(auth *bind.TransactOpts, backend bind.ContractBackend, _candidates []common.Address, candidateOwners []common.Address, _epochPeriod *big.Int, _startBlock *big.Int, _maxValidatorSize *big.Int, _minValidatorStake *big.Int, _minVoteCap *big.Int, _admin common.Address) (common.Address, *types.Transaction, *StakingContracts, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingContractsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StakingContractsBin), backend, _candidates, candidateOwners, _epochPeriod, _startBlock, _maxValidatorSize, _minValidatorStake, _minVoteCap, _admin)
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

// CandidateVoters is a free data retrieval call binding the contract method 0x7b728966.
//
// Solidity: function candidateVoters(address , uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsCaller) CandidateVoters(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "candidateVoters", arg0, arg1)
	return *ret0, err
}

// CandidateVoters is a free data retrieval call binding the contract method 0x7b728966.
//
// Solidity: function candidateVoters(address , uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsSession) CandidateVoters(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _StakingContracts.Contract.CandidateVoters(&_StakingContracts.CallOpts, arg0, arg1)
}

// CandidateVoters is a free data retrieval call binding the contract method 0x7b728966.
//
// Solidity: function candidateVoters(address , uint256 ) constant returns(address)
func (_StakingContracts *StakingContractsCallerSession) CandidateVoters(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _StakingContracts.Contract.CandidateVoters(&_StakingContracts.CallOpts, arg0, arg1)
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

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isActiveCandidate, address _owner, uint256 _totalStake)
func (_StakingContracts *StakingContractsCaller) GetCandidateData(opts *bind.CallOpts, _candidate common.Address) (struct {
	IsActiveCandidate bool
	Owner             common.Address
	TotalStake        *big.Int
}, error) {
	ret := new(struct {
		IsActiveCandidate bool
		Owner             common.Address
		TotalStake        *big.Int
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "getCandidateData", _candidate)
	return *ret, err
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isActiveCandidate, address _owner, uint256 _totalStake)
func (_StakingContracts *StakingContractsSession) GetCandidateData(_candidate common.Address) (struct {
	IsActiveCandidate bool
	Owner             common.Address
	TotalStake        *big.Int
}, error) {
	return _StakingContracts.Contract.GetCandidateData(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateData is a free data retrieval call binding the contract method 0x2a466ac7.
//
// Solidity: function getCandidateData(address _candidate) constant returns(bool _isActiveCandidate, address _owner, uint256 _totalStake)
func (_StakingContracts *StakingContractsCallerSession) GetCandidateData(_candidate common.Address) (struct {
	IsActiveCandidate bool
	Owner             common.Address
	TotalStake        *big.Int
}, error) {
	return _StakingContracts.Contract.GetCandidateData(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateOwner is a free data retrieval call binding the contract method 0xb642facd.
//
// Solidity: function getCandidateOwner(address _candidate) constant returns(address)
func (_StakingContracts *StakingContractsCaller) GetCandidateOwner(opts *bind.CallOpts, _candidate common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getCandidateOwner", _candidate)
	return *ret0, err
}

// GetCandidateOwner is a free data retrieval call binding the contract method 0xb642facd.
//
// Solidity: function getCandidateOwner(address _candidate) constant returns(address)
func (_StakingContracts *StakingContractsSession) GetCandidateOwner(_candidate common.Address) (common.Address, error) {
	return _StakingContracts.Contract.GetCandidateOwner(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateOwner is a free data retrieval call binding the contract method 0xb642facd.
//
// Solidity: function getCandidateOwner(address _candidate) constant returns(address)
func (_StakingContracts *StakingContractsCallerSession) GetCandidateOwner(_candidate common.Address) (common.Address, error) {
	return _StakingContracts.Contract.GetCandidateOwner(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateStake is a free data retrieval call binding the contract method 0x484da961.
//
// Solidity: function getCandidateStake(address _candidate) constant returns(uint256)
func (_StakingContracts *StakingContractsCaller) GetCandidateStake(opts *bind.CallOpts, _candidate common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getCandidateStake", _candidate)
	return *ret0, err
}

// GetCandidateStake is a free data retrieval call binding the contract method 0x484da961.
//
// Solidity: function getCandidateStake(address _candidate) constant returns(uint256)
func (_StakingContracts *StakingContractsSession) GetCandidateStake(_candidate common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetCandidateStake(&_StakingContracts.CallOpts, _candidate)
}

// GetCandidateStake is a free data retrieval call binding the contract method 0x484da961.
//
// Solidity: function getCandidateStake(address _candidate) constant returns(uint256)
func (_StakingContracts *StakingContractsCallerSession) GetCandidateStake(_candidate common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetCandidateStake(&_StakingContracts.CallOpts, _candidate)
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
// Solidity: function getListCandidates() constant returns(address[] _candidates, uint256[] stakes, uint256 epoch, uint256 validatorSize, uint256 minValidatorCap)
func (_StakingContracts *StakingContractsCaller) GetListCandidates(opts *bind.CallOpts) (struct {
	Candidates      []common.Address
	Stakes          []*big.Int
	Epoch           *big.Int
	ValidatorSize   *big.Int
	MinValidatorCap *big.Int
}, error) {
	ret := new(struct {
		Candidates      []common.Address
		Stakes          []*big.Int
		Epoch           *big.Int
		ValidatorSize   *big.Int
		MinValidatorCap *big.Int
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "getListCandidates")
	return *ret, err
}

// GetListCandidates is a free data retrieval call binding the contract method 0x690ff8a1.
//
// Solidity: function getListCandidates() constant returns(address[] _candidates, uint256[] stakes, uint256 epoch, uint256 validatorSize, uint256 minValidatorCap)
func (_StakingContracts *StakingContractsSession) GetListCandidates() (struct {
	Candidates      []common.Address
	Stakes          []*big.Int
	Epoch           *big.Int
	ValidatorSize   *big.Int
	MinValidatorCap *big.Int
}, error) {
	return _StakingContracts.Contract.GetListCandidates(&_StakingContracts.CallOpts)
}

// GetListCandidates is a free data retrieval call binding the contract method 0x690ff8a1.
//
// Solidity: function getListCandidates() constant returns(address[] _candidates, uint256[] stakes, uint256 epoch, uint256 validatorSize, uint256 minValidatorCap)
func (_StakingContracts *StakingContractsCallerSession) GetListCandidates() (struct {
	Candidates      []common.Address
	Stakes          []*big.Int
	Epoch           *big.Int
	ValidatorSize   *big.Int
	MinValidatorCap *big.Int
}, error) {
	return _StakingContracts.Contract.GetListCandidates(&_StakingContracts.CallOpts)
}

// GetVoterStake is a free data retrieval call binding the contract method 0x158a65f6.
//
// Solidity: function getVoterStake(address _candidate, address _voter) constant returns(uint256 stake)
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
// Solidity: function getVoterStake(address _candidate, address _voter) constant returns(uint256 stake)
func (_StakingContracts *StakingContractsSession) GetVoterStake(_candidate common.Address, _voter common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStake(&_StakingContracts.CallOpts, _candidate, _voter)
}

// GetVoterStake is a free data retrieval call binding the contract method 0x158a65f6.
//
// Solidity: function getVoterStake(address _candidate, address _voter) constant returns(uint256 stake)
func (_StakingContracts *StakingContractsCallerSession) GetVoterStake(_candidate common.Address, _voter common.Address) (*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStake(&_StakingContracts.CallOpts, _candidate, _voter)
}

// GetVoterStakes is a free data retrieval call binding the contract method 0xe2db89b5.
//
// Solidity: function getVoterStakes(address _candidate, address[] voters) constant returns(uint256[] stakes)
func (_StakingContracts *StakingContractsCaller) GetVoterStakes(opts *bind.CallOpts, _candidate common.Address, voters []common.Address) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getVoterStakes", _candidate, voters)
	return *ret0, err
}

// GetVoterStakes is a free data retrieval call binding the contract method 0xe2db89b5.
//
// Solidity: function getVoterStakes(address _candidate, address[] voters) constant returns(uint256[] stakes)
func (_StakingContracts *StakingContractsSession) GetVoterStakes(_candidate common.Address, voters []common.Address) ([]*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStakes(&_StakingContracts.CallOpts, _candidate, voters)
}

// GetVoterStakes is a free data retrieval call binding the contract method 0xe2db89b5.
//
// Solidity: function getVoterStakes(address _candidate, address[] voters) constant returns(uint256[] stakes)
func (_StakingContracts *StakingContractsCallerSession) GetVoterStakes(_candidate common.Address, voters []common.Address) ([]*big.Int, error) {
	return _StakingContracts.Contract.GetVoterStakes(&_StakingContracts.CallOpts, _candidate, voters)
}

// GetVoters is a free data retrieval call binding the contract method 0x2d15cc04.
//
// Solidity: function getVoters(address _candidate) constant returns(address[] voters)
func (_StakingContracts *StakingContractsCaller) GetVoters(opts *bind.CallOpts, _candidate common.Address) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getVoters", _candidate)
	return *ret0, err
}

// GetVoters is a free data retrieval call binding the contract method 0x2d15cc04.
//
// Solidity: function getVoters(address _candidate) constant returns(address[] voters)
func (_StakingContracts *StakingContractsSession) GetVoters(_candidate common.Address) ([]common.Address, error) {
	return _StakingContracts.Contract.GetVoters(&_StakingContracts.CallOpts, _candidate)
}

// GetVoters is a free data retrieval call binding the contract method 0x2d15cc04.
//
// Solidity: function getVoters(address _candidate) constant returns(address[] voters)
func (_StakingContracts *StakingContractsCallerSession) GetVoters(_candidate common.Address) ([]common.Address, error) {
	return _StakingContracts.Contract.GetVoters(&_StakingContracts.CallOpts, _candidate)
}

// GetWithdrawCap is a free data retrieval call binding the contract method 0x15febd68.
//
// Solidity: function getWithdrawCap(uint256 epoch) constant returns(uint256 cap)
func (_StakingContracts *StakingContractsCaller) GetWithdrawCap(opts *bind.CallOpts, epoch *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getWithdrawCap", epoch)
	return *ret0, err
}

// GetWithdrawCap is a free data retrieval call binding the contract method 0x15febd68.
//
// Solidity: function getWithdrawCap(uint256 epoch) constant returns(uint256 cap)
func (_StakingContracts *StakingContractsSession) GetWithdrawCap(epoch *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.GetWithdrawCap(&_StakingContracts.CallOpts, epoch)
}

// GetWithdrawCap is a free data retrieval call binding the contract method 0x15febd68.
//
// Solidity: function getWithdrawCap(uint256 epoch) constant returns(uint256 cap)
func (_StakingContracts *StakingContractsCallerSession) GetWithdrawCap(epoch *big.Int) (*big.Int, error) {
	return _StakingContracts.Contract.GetWithdrawCap(&_StakingContracts.CallOpts, epoch)
}

// GetWithdrawEpochs is a free data retrieval call binding the contract method 0x0e0516aa.
//
// Solidity: function getWithdrawEpochs() constant returns(uint256[] epochs)
func (_StakingContracts *StakingContractsCaller) GetWithdrawEpochs(opts *bind.CallOpts) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "getWithdrawEpochs")
	return *ret0, err
}

// GetWithdrawEpochs is a free data retrieval call binding the contract method 0x0e0516aa.
//
// Solidity: function getWithdrawEpochs() constant returns(uint256[] epochs)
func (_StakingContracts *StakingContractsSession) GetWithdrawEpochs() ([]*big.Int, error) {
	return _StakingContracts.Contract.GetWithdrawEpochs(&_StakingContracts.CallOpts)
}

// GetWithdrawEpochs is a free data retrieval call binding the contract method 0x0e0516aa.
//
// Solidity: function getWithdrawEpochs() constant returns(uint256[] epochs)
func (_StakingContracts *StakingContractsCallerSession) GetWithdrawEpochs() ([]*big.Int, error) {
	return _StakingContracts.Contract.GetWithdrawEpochs(&_StakingContracts.CallOpts)
}

// GetWithdrawEpochsAndCaps is a free data retrieval call binding the contract method 0xd5816bfa.
//
// Solidity: function getWithdrawEpochsAndCaps() constant returns(uint256[] epochs, uint256[] caps)
func (_StakingContracts *StakingContractsCaller) GetWithdrawEpochsAndCaps(opts *bind.CallOpts) (struct {
	Epochs []*big.Int
	Caps   []*big.Int
}, error) {
	ret := new(struct {
		Epochs []*big.Int
		Caps   []*big.Int
	})
	out := ret
	err := _StakingContracts.contract.Call(opts, out, "getWithdrawEpochsAndCaps")
	return *ret, err
}

// GetWithdrawEpochsAndCaps is a free data retrieval call binding the contract method 0xd5816bfa.
//
// Solidity: function getWithdrawEpochsAndCaps() constant returns(uint256[] epochs, uint256[] caps)
func (_StakingContracts *StakingContractsSession) GetWithdrawEpochsAndCaps() (struct {
	Epochs []*big.Int
	Caps   []*big.Int
}, error) {
	return _StakingContracts.Contract.GetWithdrawEpochsAndCaps(&_StakingContracts.CallOpts)
}

// GetWithdrawEpochsAndCaps is a free data retrieval call binding the contract method 0xd5816bfa.
//
// Solidity: function getWithdrawEpochsAndCaps() constant returns(uint256[] epochs, uint256[] caps)
func (_StakingContracts *StakingContractsCallerSession) GetWithdrawEpochsAndCaps() (struct {
	Epochs []*big.Int
	Caps   []*big.Int
}, error) {
	return _StakingContracts.Contract.GetWithdrawEpochsAndCaps(&_StakingContracts.CallOpts)
}

// IsCandidate is a free data retrieval call binding the contract method 0xd51b9e93.
//
// Solidity: function isCandidate(address _candidate) constant returns(bool)
func (_StakingContracts *StakingContractsCaller) IsCandidate(opts *bind.CallOpts, _candidate common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _StakingContracts.contract.Call(opts, out, "isCandidate", _candidate)
	return *ret0, err
}

// IsCandidate is a free data retrieval call binding the contract method 0xd51b9e93.
//
// Solidity: function isCandidate(address _candidate) constant returns(bool)
func (_StakingContracts *StakingContractsSession) IsCandidate(_candidate common.Address) (bool, error) {
	return _StakingContracts.Contract.IsCandidate(&_StakingContracts.CallOpts, _candidate)
}

// IsCandidate is a free data retrieval call binding the contract method 0xd51b9e93.
//
// Solidity: function isCandidate(address _candidate) constant returns(bool)
func (_StakingContracts *StakingContractsCallerSession) IsCandidate(_candidate common.Address) (bool, error) {
	return _StakingContracts.Contract.IsCandidate(&_StakingContracts.CallOpts, _candidate)
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

// WithdrawWithIndex is a paid mutator transaction binding the contract method 0x96c23442.
//
// Solidity: function withdrawWithIndex(uint256 epoch, uint256 index) returns(bool)
func (_StakingContracts *StakingContractsTransactor) WithdrawWithIndex(opts *bind.TransactOpts, epoch *big.Int, index *big.Int) (*types.Transaction, error) {
	return _StakingContracts.contract.Transact(opts, "withdrawWithIndex", epoch, index)
}

// WithdrawWithIndex is a paid mutator transaction binding the contract method 0x96c23442.
//
// Solidity: function withdrawWithIndex(uint256 epoch, uint256 index) returns(bool)
func (_StakingContracts *StakingContractsSession) WithdrawWithIndex(epoch *big.Int, index *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.WithdrawWithIndex(&_StakingContracts.TransactOpts, epoch, index)
}

// WithdrawWithIndex is a paid mutator transaction binding the contract method 0x96c23442.
//
// Solidity: function withdrawWithIndex(uint256 epoch, uint256 index) returns(bool)
func (_StakingContracts *StakingContractsTransactorSession) WithdrawWithIndex(epoch *big.Int, index *big.Int) (*types.Transaction, error) {
	return _StakingContracts.Contract.WithdrawWithIndex(&_StakingContracts.TransactOpts, epoch, index)
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
	Epoch     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterResigned is a free log retrieval operation binding the contract event 0x886e0db046874dde595498040d176412e81183750ceb33fc46f0450362bbc241.
//
// Solidity: event Resigned(address _candidate, uint256 _epoch)
func (_StakingContracts *StakingContractsFilterer) FilterResigned(opts *bind.FilterOpts) (*StakingContractsResignedIterator, error) {

	logs, sub, err := _StakingContracts.contract.FilterLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return &StakingContractsResignedIterator{contract: _StakingContracts.contract, event: "Resigned", logs: logs, sub: sub}, nil
}

// WatchResigned is a free log subscription operation binding the contract event 0x886e0db046874dde595498040d176412e81183750ceb33fc46f0450362bbc241.
//
// Solidity: event Resigned(address _candidate, uint256 _epoch)
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
