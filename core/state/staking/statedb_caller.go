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

// GetCandidates returns list candidate's address
func (c *stateDBStakingCaller) GetCandidates(stakingContractAddr common.Address) ([]common.Address, error) {
	slotHash := c.config.CandidatesLayout.slotHash()
	arrLength := c.stateDB.GetState(stakingContractAddr, slotHash)
	if arrLength.Big().Cmp(big.NewInt(0)) == 0 {
		return nil, ErrEmptyValidatorSet
	}
	var candidates []common.Address
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		key := getLocDynamicArrAtElement(slotHash, i, 1)
		ret := c.stateDB.GetState(stakingContractAddr, key)
		candidates = append(candidates, common.HexToAddress(ret.Hex()))
	}
	return candidates, nil
}

// GetCandidateOwner returns current owner of a candidate
func (c *stateDBStakingCaller) GetCandidateOwner(stakingContractAddr common.Address, candidate common.Address) common.Address {
	locCandidateOwner := getStorageLocation(c.config.CandidateDataLayout, candidate.Hash(), 2)
	ret := c.stateDB.GetState(stakingContractAddr, locCandidateOwner)
	return common.HexToAddress(ret.Hex())
}

// GetCandidateStake returns current stake of a candidate
func (c *stateDBStakingCaller) GetCandidateStake(stakingContractAddr common.Address, candidate common.Address) *big.Int {
	locStake := getStorageLocation(c.config.CandidateDataLayout, candidate.Hash(), 1)
	stake := c.stateDB.GetState(stakingContractAddr, locStake)
	return stake.Big()
}

// GetStartBlock returns the startblock
func (c *stateDBStakingCaller) GetStartBlock(stakingContractAddr common.Address) int {
	ret := c.stateDB.GetState(stakingContractAddr, c.config.StartBlockLayout.slotHash())
	return int(ret.Big().Int64())
}

// GetEpochPeriod returns the epochperiod
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

func getStorageLocation(LayOut LayOut, keyHash common.Hash, index uint) common.Hash {
	locState := getLocMappingAtKey(keyHash, LayOut.Slot)
	return common.BigToHash(locState.Add(locState, new(big.Int).SetUint64(uint64(index))))
}

func getLocDynamicArrAtElement(slotHash common.Hash, index uint64, elementSize uint64) common.Hash {
	slotKecBig := crypto.Keccak256Hash(slotHash.Bytes()).Big()
	arrBig := slotKecBig.Add(slotKecBig, new(big.Int).SetUint64(index*elementSize))
	return common.BigToHash(arrBig)
}

func getLocMappingAtKey(key common.Hash, slot uint64) *big.Int {
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	retByte := crypto.Keccak256(key.Bytes(), slotHash.Bytes())
	ret := new(big.Int)
	ret.SetBytes(retByte)
	return ret
}
