package backend

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	evrynetCore "github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/params"
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

	if signer != tests_utils.GetAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), tests_utils.GetAddress().Hex())
	}
}

func TestValidators(t *testing.T) {
	var (
		nodePrivateKey = tests_utils.MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
		be            = mustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader)
	)

	valSet0 := be.Validators(big.NewInt(0))

	assert.Equal(t, 1, valSet0.Size())

	list := valSet0.List()
	log.Println("validator set of block 0 is")

	for _, val := range list {
		log.Println(val.String())
	}

	valSet1 := be.Validators(big.NewInt(1))

	assert.Equal(t, 1, valSet1.Size())

	list = valSet1.List()
	log.Println("validator set of block 1 is")

	for _, val := range list {
		log.Println(val.String())
	}

	valSet2 := be.Validators(big.NewInt(2))
	assert.Equal(t, 0, valSet2.Size())
}

func mustCreateAndStartNewBackend(t *testing.T, nodePrivateKey *ecdsa.PrivateKey, genesisHeader *types.Header) *backend {
	var (
		address = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		trigger = false
		statedb = tests_utils.MustCreateStateDB(t)

		testTxPoolConfig evrynetCore.TxPoolConfig
		blockchain       = &tests_utils.MockChainReader{
			GenesisHeader: genesisHeader,
			MockBlockChain: &tests_utils.MockBlockChain{
				Statedb:       statedb,
				GasLimit:      1000000000,
				ChainHeadFeed: new(event.Feed),
			},
			Address: address,
			Trigger: &trigger,
		}
		pool   = evrynetCore.NewTxPool(testTxPoolConfig, params.TendermintTestChainConfig, blockchain)
		memDB  = ethdb.NewMemDatabase()
		config = tendermint.DefaultConfig
		be     = New(config, nodePrivateKey, pool, WithDB(memDB)).(*backend)
	)
	statedb.SetBalance(address, new(big.Int).SetUint64(params.Ether))
	defer pool.Stop()
	be.chain = blockchain
	be.currentBlock = blockchain.CurrentBlock

	return be
}
