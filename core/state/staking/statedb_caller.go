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

// GetCandidates returns list candidate's address
func (c *stateDBStakingCaller) GetCandidates(stakingContractAddr common.Address) ([]common.Address, error) {
	arrLength := c.getStorageLocation(stakingContractAddr, c.config.CandidatesLayout, common.Hash{}, nil)
	if arrLength.Big().Cmp(big.NewInt(0)) == 0 {
		return nil, ErrEmptyValidatorSet
	}
	var candidates []common.Address
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		ret := c.getStorageLocation(stakingContractAddr, c.config.CandidatesLayout, common.Hash{}, new(big.Int).SetUint64(i))
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
	ret := c.getStorageLocation(stakingContractAddr, c.config.CandidateDataLayout, candidate.Hash(), new(big.Int).SetUint64(1))
	return ret.Big()
}

// GetStartBlock returns the startblock
func (c *stateDBStakingCaller) GetStartBlock(stakingContractAddr common.Address) int {
	ret := c.getStorageLocation(stakingContractAddr, c.config.StartBlockLayout, common.Hash{}, nil)
	return int(ret.Big().Int64())
}

// GetEpochPeriod returns the epochperiod
func (c *stateDBStakingCaller) GetEpochPeriod(stakingContractAddr common.Address) int {
	ret := c.getStorageLocation(stakingContractAddr, c.config.EpochPeriodLayout, common.Hash{}, nil)
	return int(ret.Big().Int64())
}

// GetMaxValidatorSize returns maximum validators allowed
func (c *stateDBStakingCaller) GetMaxValidatorSize(stakingContractAddr common.Address) int {
	ret := c.getStorageLocation(stakingContractAddr, c.config.MaxValidatorSizeLayout, common.Hash{}, nil)
	return int(ret.Big().Int64())
}

// GetMinValidatorStake returns the min stake of a validator
func (c *stateDBStakingCaller) GetMinValidatorStake(stakingContractAddr common.Address) *big.Int {
	ret := c.getStorageLocation(stakingContractAddr, c.config.MinValidatorStakeLayout, common.Hash{}, nil)
	return ret.Big()
}

// GetMinVoterCap returns the MinVoterCap
func (c *stateDBStakingCaller) GetMinVoterCap(stakingContractAddr common.Address) *big.Int {
	ret := c.getStorageLocation(stakingContractAddr, c.config.MinVoterCapLayout, common.Hash{}, nil)
	return ret.Big()
}

// GetAdmin returns admin's address
func (c *stateDBStakingCaller) GetAdmin(stakingContractAddr common.Address) common.Address {
	ret := c.getStorageLocation(stakingContractAddr, c.config.AdminLayout, common.Hash{}, nil)
	return common.HexToAddress(ret.Hex())
}

// GetCandidateOwner returns current owner of a candidate
func (c *stateDBStakingCaller) GetCandidateOwner(stakingContractAddr common.Address, candidate common.Address) common.Address {
	ret := c.getStorageLocation(stakingContractAddr, c.config.CandidateDataLayout, candidate.Hash(), new(big.Int).SetUint64(2))
	return common.HexToAddress(ret.Hex())
}

// if its a primitive type let keyHash = zeroHash and index = nil
// elsif its an array type let keyHash = zeroHash and index = index of element
// else its a map type let keyHash = hash of key and index = index of element
func (c *stateDBStakingCaller) getStorageLocation(stakingContractAddr common.Address, layout LayOut, keyHash common.Hash, index *big.Int) common.Hash {
	var (
		emptyHash = common.Hash{}
		key       common.Hash
	)

	if index == nil {
		// it's a primitive type
		return c.stateDB.GetState(stakingContractAddr, layout.slotHash())
	}

	if keyHash == emptyHash {
		// its an array
		key = getLocDynamicArrAtElement(layout.slotHash(), index.Uint64(), 1)
		return c.stateDB.GetState(stakingContractAddr, key)
	}

	// its a map
	locState := getLocMappingAtKey(keyHash, layout.Slot)
	key = common.BigToHash(locState.Add(locState, index))
	return c.stateDB.GetState(stakingContractAddr, key)
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
