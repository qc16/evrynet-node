package staking

import (
	"math/big"
	"sort"
	"strings"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/core/state"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/log"
)

// StateDbCaller creates a wrapper with statedb to implements ContractCaller
type StateDbCaller struct {
	stateDB *state.StateDB
}

// NewStateDbCaller return instance of StakingCaller
func NewStateDbCaller(state *state.StateDB) StakingCaller {
	return &StateDbCaller{
		stateDB: state,
	}
}

var (
	slotStakingMapping = map[string]uint64{
		"admin":             1,
		"candidateVoters":   2,
		"candidateData":     3,
		"candidates":        4,
		"startBlock":        5,
		"epochPeriod":       6,
		"maxValidatorSize":  7,
		"minValidatorStake": 8,
		"minVoterCap":       9,
		"withdrawsState":    10,
	}
)

// GetCandidateStake returns current stake of a candidate
func (c *StateDbCaller) GetCandidateStake(scAddress common.Address, candidate common.Address) *big.Int {
	slot := slotStakingMapping["candidateData"]
	locValidatorsState := getLocMappingAtKey(candidate.Hash(), slot)
	locStake := locValidatorsState.Add(locValidatorsState, new(big.Int).SetUint64(uint64(1)))
	stake := c.stateDB.GetState(scAddress, common.BigToHash(locStake))
	return stake.Big()
}

// GetValidators returns validators from stateDB and block number of the caller by smart-contract's address
func (c *StateDbCaller) GetValidators(scAddress common.Address) ([]common.Address, error) {
	candidates := c.getCandidates(scAddress)
	if candidates == nil || len(candidates) == 0 {
		log.Warn("statedb.GetState returns empty hash")
		return nil, ErrEmptyValidatorSet
	}

	var (
		validators []common.Address
		stakes     = make(map[common.Address]*big.Int)
	)
	minValStake := c.getMinValidatorStake(scAddress)
	for _, candidate := range candidates {
		stake := c.GetCandidateStake(scAddress, candidate)
		if stake.Cmp(minValStake) < 0 {
			continue
		}

		validators = append(validators, candidate)
		stakes[candidate] = stake
	}

	maxValSize := c.getMaxValidatorSize(scAddress)
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

func (c *StateDbCaller) getMaxValidatorSize(scAddress common.Address) int {
	slot := slotStakingMapping["maxValidatorSize"]
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	ret := c.stateDB.GetState(scAddress, slotHash)
	return int(ret.Big().Int64())
}

func (c *StateDbCaller) getMinValidatorStake(scAddress common.Address) *big.Int {
	slot := slotStakingMapping["minValidatorStake"]
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	ret := c.stateDB.GetState(scAddress, slotHash)
	return ret.Big()
}

func (c *StateDbCaller) getCandidates(scAddress common.Address) []common.Address {
	slot := slotStakingMapping["candidates"]
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	arrLength := c.stateDB.GetState(scAddress, slotHash)
	keys := []common.Hash{}
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		key := getLocDynamicArrAtElement(slotHash, i, 1)
		keys = append(keys, key)
	}
	rets := []common.Address{}
	for _, key := range keys {
		ret := c.stateDB.GetState(scAddress, key)
		rets = append(rets, common.HexToAddress(ret.Hex()))
	}
	return rets
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
