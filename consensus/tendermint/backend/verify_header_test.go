package backend

import (
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"

)

var (
	correctParentHeader = types.Header{}
	genesisHeader = types.Header{
		ParentHash : common.HexToHash("0x01"),
		UncleHash:types.CalcUncleHash(nil),
		Root: common.HexToHash("0x0"),
	}
)

func TestBackend_VerifyHeader(t *testing.T) {
	//create New test backend and newMockChain
	//then create a valid block and verify it
}

type mockChain struct{}


//GetHeader implement a mock version of chainReader.GetHeader
//It returns correctParentHeader with Number set to input blockNumber
func (mc *mockChain) GetHeader(hash common.Hash, blockNumber uint64) *types.Header{
	var h  = types.Header{}
	h = correctParentHeader
	h.Number=  big.NewInt(int64(blockNumber))
	return &h

}

//GetHeaderByNumber implement a mock version of chainReader.GetHeaderByNumber
//It returns genesis Header if blockNumber is 0, else return an empty Header
func (mc *mockChain) GetHeaderByNumber(blockNumber uint64) *types.Header{
	if blockNumber==0 {
		return &genesisHeader
	}
	return nil
}

