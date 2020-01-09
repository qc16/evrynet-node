package backend

import (
	"crypto/ecdsa"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/crypto/secp256k1"
	"github.com/evrynet-official/evrynet-client/evrdb"
)

func TestBackend_VerifyHeader(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	chain, engine := mustStartTestChainAndBackend(nodePK, genesisHeader)
	assert.NotNil(t, chain)
	assert.NotNil(t, engine)

	// without seal
	block := tests_utils.MakeBlockWithoutSeal(genesisHeader)
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, engine.VerifyHeader(engine.chain, block.Header(), false))

	// with seal but incorrect coinbase
	block = tests_utils.MakeBlockWithSeal(engine, genesisHeader)
	header := block.Header()
	header.Coinbase = common.Address{}
	tests_utils.AppendSeal(header, engine)
	assert.Equal(t, errCoinBaseInvalid, engine.VerifyHeader(engine.chain, header, false))

	// without committed seal
	block = tests_utils.MakeBlockWithSeal(engine, genesisHeader)
	assert.Equal(t, tendermint.ErrEmptyCommittedSeals, engine.VerifyHeader(engine.chain, block.Header(), false))

	// with committed seal but is invalid
	block = tests_utils.MustMakeBlockWithCommittedSealInvalid(engine, genesisHeader)
	assert.Equal(t, errInvalidCommittedSeals, engine.VerifyHeader(engine.chain, block.Header(), false))

	// with committed seal
	block = tests_utils.MustMakeBlockWithCommittedSeal(engine, genesisHeader, validators)
	assert.NotNil(t, engine.chain)
	err = engine.VerifyHeader(engine.chain, block.Header(), false)
	assert.NoError(t, err)
}

func mustStartTestChainAndBackend(nodePK *ecdsa.PrivateKey, genesisHeader *types.Header) (*tests_utils.MockChainReader, *Backend) {
	memDB := evrdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	b, ok := New(config, nodePK, WithDB(memDB)).(*Backend)
	b.SetBroadcaster(&tests_utils.MockProtocolManager{})
	if !ok {
		panic("New() cannot be asserted back to backend")
	}
	chain := tests_utils.MockChainReader{
		GenesisHeader: genesisHeader,
	}

	snap, err := b.snapshot(&chain, 0, common.Hash{}, nil)
	if err != nil {
		panic(err)
	}
	if snap == nil {
		panic("failed to get snapshot")
	}
	currentBlock := func() *types.Block {
		tests_utils.AppendSeal(genesisHeader, b)
		return types.NewBlockWithHeader(genesisHeader)
	}
	if err := b.Start(&chain, currentBlock); err != nil {
		log.Panicf("cannot start backend, error:%v", err)
	}
	return &chain, b
}
