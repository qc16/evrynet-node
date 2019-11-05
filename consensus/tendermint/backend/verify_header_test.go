package backend

import (
	"testing"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
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
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be := mustCreateAndStartNewBackend(t, nodePK, genesisHeader)

	// without seal
	block := tests_utils.MakeBlockWithoutSeal(genesisHeader)
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, be.VerifyHeader(be.chain, block.Header(), false))

	// with seal but incorrect coinbase
	block = tests_utils.MakeBlockWithSeal(be, genesisHeader)
	header := block.Header()
	header.Coinbase = common.Address{}
	tests_utils.AppendSeal(header, be)
	assert.Equal(t, errCoinBaseInvalid, be.VerifyHeader(be.chain, header, false))

	// without committed seal
	block = tests_utils.MakeBlockWithSeal(be, genesisHeader)
	assert.Equal(t, tendermint.ErrEmptyCommittedSeals, be.VerifyHeader(be.chain, block.Header(), false))

	// with committed seal but is invalid
	block = tests_utils.MustMakeBlockWithCommittedSealInvalid(be, genesisHeader)
	assert.Equal(t, errInvalidSignature, be.VerifyHeader(be.chain, block.Header(), false))

	// with committed seal
	block = tests_utils.MustMakeBlockWithCommittedSeal(be, genesisHeader, validators)
	assert.NotNil(t, be.chain)
	err = be.VerifyHeader(be.chain, block.Header(), false)
	assert.NoError(t, err)
}
