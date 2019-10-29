package backend

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	evrynetCore "github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1")
	b := &backend{
		privateKey: privateKey,
	}
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Check signature recover
	hashData := crypto.Keccak256([]byte(data))
	pubkey, _ := crypto.Ecrecover(hashData, sig)

	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	if signer != tests.GetAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), tests.GetAddress().Hex())
	}
}

func TestValidators(t *testing.T) {
	var (
		nodePrivateKey = tests.MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)

	//create New test backend and newMockChain
	be, ok := mustCreateAndStartNewBackend(nodePrivateKey, genesisHeader)
	assert.True(t, ok)

	valSet0 := be.Validators(big.NewInt(0))
	if valSet0.Size() != 1 {
		t.Errorf("Valset size of zero block should be 1, get: %d", valSet0.Size())
	}
	list := valSet0.List()
	log.Println("validator set of block 0 is")

	for _, val := range list {
		log.Println(val.String())
	}
	valSet1 := be.Validators(big.NewInt(1))
	if valSet1.Size() != 1 {
		t.Errorf("Valset size of block 1st should be 1, get: %d", valSet1.Size())
	}
	list = valSet1.List()
	log.Println("validator set of block 1 is")

	for _, val := range list {
		log.Println(val.String())
	}
	valSet2 := be.Validators(big.NewInt(2))
	if valSet2.Size() != 0 {
		t.Errorf("Valset size of block 2th should be 0, get: %d", valSet2.Size())
	}
}

func mustCreateAndStartNewBackend(nodePrivateKey *ecdsa.PrivateKey, genesisHeader *types.Header) (tests.TestBackend, bool) {
	address := crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
	trigger := false
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	statedb.SetBalance(address, new(big.Int).SetUint64(params.Ether))
	var testTxPoolConfig evrynetCore.TxPoolConfig
	blockchain := &tests.TestChain{
		GenesisHeader: genesisHeader,
		TestBlockChain: &tests.TestBlockChain{
			Statedb:       statedb,
			GasLimit:      1000000000,
			ChainHeadFeed: new(event.Feed),
		},
		Address: address,
		Trigger: &trigger,
	}
	pool := evrynetCore.NewTxPool(testTxPoolConfig, params.TendermintTestChainConfig, blockchain)
	defer pool.Stop()
	memDB := ethdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	be := New(config, nodePrivateKey, WithTxPool(pool), WithDB(memDB)).(tests.TestBackend)
	ok := tests.MustStartTestChainAndBackend(be, blockchain)
	return be, ok
}
