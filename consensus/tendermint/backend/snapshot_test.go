package backend

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/ethdb"
)

func TestSaveAndLoad(t *testing.T) {
	var (
		hash   = common.HexToHash("1234567890")
		valSet = validator.NewSet([]common.Address{
			common.HexToAddress("1234567894"),
			common.HexToAddress("1234567895"),
		}, tendermint.RoundRobin)
		snap = newSnapshot(5, 10, hash, valSet)
	)

	db := ethdb.NewMemDatabase()
	assert.NoError(t, snap.store(db))

	snap1, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	assert.NoError(t, err)
	assert.NotNil(t, snap1)
	assert.Equal(t,snap.Epoch, snap1.Epoch)
	assert.Equal(t, snap.Hash, snap1.Hash)
		if !reflect.DeepEqual(snap.ValSet, snap.ValSet) {
			t.Errorf("validator set mismatch: have %v, want %v", snap1.ValSet, snap.ValSet)
		}
}
