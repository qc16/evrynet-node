package backend

import (
	"reflect"
	"testing"

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
	err := snap.store(db)
	if err != nil {
		t.Errorf("store snapshot failed: %v", err)
	}

	snap1, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	if err != nil {
		t.Errorf("load snapshot failed: %v", err)
	}
	if snap1 != nil {
		if snap.Epoch != snap1.Epoch {
			t.Errorf("epoch mismatch: have %v, want %v", snap1.Epoch, snap.Epoch)
		}
		if snap.Hash != snap1.Hash {
			t.Errorf("hash mismatch: have %v, want %v", snap1.Number, snap.Number)
		}
		if !reflect.DeepEqual(snap.ValSet, snap.ValSet) {
			t.Errorf("validator set mismatch: have %v, want %v", snap1.ValSet, snap.ValSet)
		}
	}
}
