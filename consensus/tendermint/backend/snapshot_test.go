package backend

import (
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/stretchr/testify/assert"
)

type testerVote struct {
	validator string
	voted     string
	auth      bool
}

// testerAccountPool is a pool to maintain currently active tester accounts,
// mapped from textual names used in the tests below to actual Ethereum private
// keys capable of signing transactions.
type testerAccountPool struct {
	accounts map[string]*ecdsa.PrivateKey
}

func newTesterAccountPool() *testerAccountPool {
	return &testerAccountPool{
		accounts: make(map[string]*ecdsa.PrivateKey),
	}
}

func (ap *testerAccountPool) sign(header *types.Header, validator string) {
	// Ensure we have a persistent key for the validator
	if ap.accounts[validator] == nil {
		ap.accounts[validator], _ = crypto.GenerateKey()
	}
	// Sign the header and embed the signature in extra data
	hashData := crypto.Keccak256([]byte(utils.SigHash(header).Bytes()))
	sig, _ := crypto.Sign(hashData, ap.accounts[validator])

	utils.WriteSeal(header, sig)
}

func (ap *testerAccountPool) address(account string) common.Address {
	// Ensure we have a persistent key for the account
	if ap.accounts[account] == nil {
		ap.accounts[account], _ = crypto.GenerateKey()
	}
	// Resolve and return the Ethereum address
	return crypto.PubkeyToAddress(ap.accounts[account].PublicKey)
}

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

func TestVoting(t *testing.T) {
	// Define the various voting scenarios to test
	// tests := []struct {
	// 	epoch      uint64
	// 	validators []string
	// 	votes      []testerVote
	// 	results    []string
	// }{{
	// 	// Single validator, no votes cast
	// 	validators: []string{"A"},
	// 	votes:      []testerVote{{validator: "A"}},
	// 	results:    []string{"A"},
	// }, {
	// 	// Single validator, voting to add two others (only accept first, second needs 2 votes)
	// 	validators: []string{"A"},
	// 	votes: []testerVote{
	// 		{validator: "A", voted: "B", auth: true},
	// 		{validator: "B"},
	// 		{validator: "A", voted: "C", auth: true},
	// 	},
	// 	results: []string{"A", "B"},
	// }}

	// for i, tt := range tests {
	// 	// Create the account pool and generate the initial set of validators
	// 	accounts := newTesterAccountPool()

	// 	validators := make([]common.Address, len(tt.validators))
	// 	for j, validator := range tt.validators {
	// 		validators[j] = accounts.address(validator)
	// 	}
	// 	for j := 0; j < len(validators); j++ {
	// 		for k := j + 1; k < len(validators); k++ {
	// 			if bytes.Compare(validators[j][:], validators[k][:]) > 0 {
	// 				validators[j], validators[k] = validators[k], validators[j]
	// 			}
	// 		}
	// 	}
	// }
}
