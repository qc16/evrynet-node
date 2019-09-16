package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/crypto/secp256k1"
)

func TestBackend_VerifyHeader(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = makeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	chain, engine := mustStartTestChainAndBackend(nodePK, genesisHeader)
	assert.NotNil(t, chain)
	assert.NotNil(t, engine)
	assert.Equal(t, true, engine.coreStarted)

	// without seal
	block := makeBlockWithoutSeal(genesisHeader)
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, engine.VerifyHeader(chain, block.Header(), false))

	// with seal but incorrect coinbase
	block = makeBlockWithSeal(engine, genesisHeader)
	header := block.Header()
	header.Coinbase = common.Address{}
	appendSeal(header, engine)
	assert.Equal(t, errCoinBaseInvalid, engine.VerifyHeader(chain, header, false))

	// without committed seal
	block = makeBlockWithSeal(engine, genesisHeader)
	assert.Equal(t, errEmptyCommittedSeals, engine.VerifyHeader(chain, block.Header(), false))

	// with committed seal but is invalid
	block = mustMakeBlockWithCommittedSealInvalid(engine, genesisHeader)
	assert.Equal(t, errInvalidSignature, engine.VerifyHeader(chain, block.Header(), false))

	// with committed seal
	block = mustMakeBlockWithCommittedSeal(engine, genesisHeader, validators)
	assert.NotNil(t, chain)
	err = engine.VerifyHeader(chain, block.Header(), false)
	assert.NoError(t, err)

}
