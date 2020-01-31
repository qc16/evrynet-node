package state

import (
	"math/big"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/crypto"
)

var (
	indexValidatorMapping = map[string]uint64{
		"validators": 0,
	}
)

// GetValidators returns validators from stateDB by smart-contract's address
func GetValidators(stateDB *StateDB, scAddress common.Address) []common.Address {
	index := indexValidatorMapping["validators"]
	indexHash := common.BigToHash(new(big.Int).SetUint64(index))
	arrLength := stateDB.GetState(scAddress, indexHash)
	var keys []common.Hash
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		key := getLocDynamicArrAtElement(indexHash, i, 1)
		keys = append(keys, key)
	}
	var result []common.Address
	for _, key := range keys {
		ret := stateDB.GetState(scAddress, key)
		result = append(result, common.HexToAddress(ret.Hex()))
	}
	return result
}

func getLocDynamicArrAtElement(indexHash common.Hash, index uint64, elementSize uint64) common.Hash {
	indexKecBig := crypto.Keccak256Hash(indexHash.Bytes()).Big()
	//arrBig = slotKecBig + index * elementSize
	arrBig := indexKecBig.Add(indexKecBig, new(big.Int).SetUint64(index*elementSize))
	return common.BigToHash(arrBig)
}
