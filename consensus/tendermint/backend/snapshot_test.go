package backend

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethdb"
)

func TestSaveAndLoad(t *testing.T) {
	var (
		hash   = common.HexToHash("1234567890")
		valSet = validator.NewSet([]common.Address{
			common.HexToAddress("1234567894"),
			common.HexToAddress("1234567895"),
		}, tendermint.RoundRobin, int64(0))
		snap = newSnapshot(5, 10, hash, valSet)
	)

	db := ethdb.NewMemDatabase()
	assert.NoError(t, snap.store(db))

	snap1, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	assert.NoError(t, err)
	assert.NotNil(t, snap1)
	assert.Equal(t, snap.Epoch, snap1.Epoch)
	assert.Equal(t, snap.Hash, snap1.Hash)
	if !reflect.DeepEqual(snap.ValSet, snap.ValSet) {
		t.Errorf("validator set mismatch: have %v, want %v", snap1.ValSet, snap.ValSet)
	}
}

func TestApplyHeaders(t *testing.T) {
	var (
		hash   = common.HexToHash("1234567890")
		addr1  = common.HexToAddress("1")
		addr2  = common.HexToAddress("2")
		addr3  = common.HexToAddress("3")
		addr4  = common.HexToAddress("4")
		valSet = validator.NewSet([]common.Address{
			addr1, addr2, addr3, addr4,
		}, tendermint.RoundRobin, int64(1))

		//init snapshot with epoch is 5, number is 0
		snap    = newSnapshot(5, 1, hash, valSet)
		headers []*types.Header
	)

	if val := snap.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr1) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr1)
	}

	//Apply new 2 headers, new proposer should be addr2
	headers = createHeaderArr(2, 2)
	snap1, err := snap.apply(headers)
	assert.NoError(t, err)

	if val := snap1.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr3) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr3)
	}

	//Apply new 5 headers, new proposer should be addr2
	headers = createHeaderArr(2, 5)
	snap2, err := snap.apply(headers)
	assert.NoError(t, err)

	if val := snap2.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr2) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr2)
	}

	// test save and load after that apply headers
	db := ethdb.NewMemDatabase()
	assert.NoError(t, snap.store(db))
	snap3, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	if val := snap3.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr1) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr1)
	}

	//Apply new 9 headers, new proposer should be addr2
	headers = createHeaderArr(2, 9)
	snap4, err := snap3.apply(headers)
	assert.NoError(t, err)

	if val := snap4.ValSet.GetProposer(); !reflect.DeepEqual(val.Address(), addr2) {
		t.Errorf("validator mismatch: have %v, want %v", val.Address(), addr2)
	}

}

func createHeaderArr(startNumber int, countHeader int) []*types.Header {
	headers := make([]*types.Header, countHeader)
	for i := startNumber; i < startNumber+countHeader; i++ {
		headers[i-startNumber] = &types.Header{
			Number: big.NewInt(int64(i)),
		}
	}
	return headers
}
