package staking

import (
	"math/big"
	"sort"
	"strings"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/state"
	"github.com/Evrynetlabs/evrynet-node/crypto"
)

// stateDBStakingCaller creates a wrapper with statedb to implements ContractCaller
// based on https://solidity.readthedocs.io/en/develop/miscellaneous.html
type stateDBStakingCaller struct {
	stateDB *state.StateDB
	config  *IndexConfigs
}

// NewStateDbStakingCaller return instance of StakingCaller which reads data directly from state DB
func NewStateDbStakingCaller(state *state.StateDB, cfg *IndexConfigs) StakingCaller {
	return &stateDBStakingCaller{
		stateDB: state,
		config:  cfg,
	}
}

func (layOut *LayOut) slotHash() common.Hash {
	return common.BigToHash(new(big.Int).SetUint64(layOut.Slot))
}

// GetCandidates returns list candidate's address
func (c *stateDBStakingCaller) GetCandidates(stakingContractAddr common.Address) ([]common.Address, error) {
	//arrLength := c.getStorageLocation(stakingContractAddr, c.config.CandidatesLayout, common.Hash{}, nil)
	arrLength := c.stateDB.GetState(stakingContractAddr, c.config.CandidatesLayout.slotHash())
	if arrLength.Big().Cmp(big.NewInt(0)) == 0 {
		return nil, ErrEmptyValidatorSet
	}
	var candidates []common.Address
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		ret := c.stateDB.GetState(stakingContractAddr, getLocDynamicArrAtElement(c.config.CandidatesLayout.slotHash(), i, 1))
		candidates = append(candidates, common.HexToAddress(ret.Hex()))
	}
	return candidates, nil
}

// GetValidators returns validators from stateDB and block number of the caller by smart-contract's address
func (c *stateDBStakingCaller) GetValidators(stakingContractAddr common.Address) ([]common.Address, error) {
	// check if this address is a valid contract, this will help us return better error
	if codes := c.stateDB.GetCode(stakingContractAddr); len(codes) == 0 {
		return nil, bind.ErrNoCode
	}

	candidates, err := c.GetCandidates(stakingContractAddr)
	if err != nil {
		return nil, err
	}

	var (
		validators []common.Address
		stakes     = make(map[common.Address]*big.Int)
	)
	minValStake := c.GetMinValidatorStake(stakingContractAddr)
	for _, candidate := range candidates {
		stake := c.GetCandidateStake(stakingContractAddr, candidate)
		if stake.Cmp(minValStake) < 0 {
			continue
		}

		validators = append(validators, candidate)
		stakes[candidate] = stake
	}

	maxValSize := c.GetMaxValidatorSize(stakingContractAddr)
	if len(validators) <= maxValSize {
		return validators, nil
	}

	sort.Slice(validators, func(i, j int) bool {
		if stakes[validators[i]].Cmp(stakes[validators[j]]) == 0 {
			return strings.Compare(validators[i].String(), validators[j].String()) > 0
		}
		return stakes[validators[i]].Cmp(stakes[validators[j]]) > 0
	})

	return candidates[:maxValSize], nil
}

func (c *stateDBStakingCaller) GetValidatorsData(common.Address, []common.Address) (map[common.Address]CandidateData, error) {
	panic("implement me")
}

// GetCandidateStake returns current stake of a candidate
func (c *stateDBStakingCaller) GetCandidateStake(stakingContractAddr common.Address, candidate common.Address) *big.Int {
	loc := getLocMapping(c.config.CandidateDataLayout.slotHash(), candidate.Hash())
	loc = getSlot(loc, big.NewInt(1))
	ret := c.stateDB.GetState(stakingContractAddr, loc)
	return ret.Big()
}

// GetStartBlock returns the startblock
func (c *stateDBStakingCaller) GetStartBlock(stakingContractAddr common.Address) int {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.StartBlockLayout.slotHash())
	return int(ret.Big().Int64())
}

// GetEpochPeriod returns the epochPeriod
func (c *stateDBStakingCaller) GetEpochPeriod(stakingContractAddr common.Address) int {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.EpochPeriodLayout.slotHash())
	return int(ret.Big().Int64())
}

// GetMaxValidatorSize returns maximum validators allowed
func (c *stateDBStakingCaller) GetMaxValidatorSize(stakingContractAddr common.Address) int {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.MaxValidatorSizeLayout.slotHash())
	return int(ret.Big().Int64())
}

// GetMinValidatorStake returns the min stake of a validator
func (c *stateDBStakingCaller) GetMinValidatorStake(stakingContractAddr common.Address) *big.Int {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.MinValidatorStakeLayout.slotHash())
	return ret.Big()
}

// GetMinVoterCap returns the MinVoterCap
func (c *stateDBStakingCaller) GetMinVoterCap(stakingContractAddr common.Address) *big.Int {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.MinVoterCapLayout.slotHash())
	return ret.Big()
}

// GetAdmin returns admin's address
func (c *stateDBStakingCaller) GetAdmin(stakingContractAddr common.Address) common.Address {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.AdminLayout.slotHash())
	return common.HexToAddress(ret.Hex())
}

// GetCandidateOwner returns current owner of a candidate
func (c *stateDBStakingCaller) GetCandidateOwner(stakingContractAddr common.Address, candidate common.Address) common.Address {
	loc := getLocMapping(c.config.CandidateDataLayout.slotHash(), candidate.Hash())
	loc = getSlot(loc, big.NewInt(2))
	ret := c.stateDB.GetState(stakingContractAddr, loc)
	return common.HexToAddress(ret.Hex())
}

/**
 * Array data is located at keccak256(p)
 */
func getLocDynamicArrAtElement(slotHash common.Hash, index uint64, elementSize uint64) common.Hash {
	slotKecBig := crypto.Keccak256Hash(slotHash.Bytes()).Big()
	arrBig := slotKecBig.Add(slotKecBig, new(big.Int).SetUint64(index*elementSize))
	return common.BigToHash(arrBig)
}

/**
 * The value to a mapping key k at position p is located at keccak256(k . p) where . is concatenation.
 */
func getLocMapping(root common.Hash, key common.Hash) common.Hash {
	return common.BytesToHash(crypto.Keccak256(key.Bytes(), root.Bytes()))
}

/**
 * Get the position for a field inside a struct
 */
func getSlot(root common.Hash, slot *big.Int) common.Hash {
	rootBig := root.Big()
	arrBig := new(big.Int).Add(rootBig, slot)
	return common.BigToHash(arrBig)
}
