package backend

import (
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/common/hexutil"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
)

// TestSimulateSubscribeAndReceiveToSeal is a simple test to pass a block to backend.Seal()
// on core.handleEvents(), the block is received.
func TestSimulateSubscribeAndReceiveToSeal(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be, ok := mustCreateAndStartNewBackend(nodePK, genesisHeader)
	assert.True(t, ok)
	assert.Equal(t, true, be.IsCoreStarted())

	// without seal
	block := tests.MakeBlockWithoutSeal(genesisHeader)
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, be.VerifyHeader(be.Chain(), block.Header(), false))

	err = be.Seal(be.Chain(), block, nil, nil)

	// Sleep to make sure that the block can be received from
	time.Sleep(2 * time.Second)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

}

// TestAuthor is a simple test to pass a block to backend.Seal()
// on core.handleEvents(), the block is received.
func TestAuthor(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be, ok := mustCreateAndStartNewBackend(nodePK, genesisHeader)
	assert.True(t, ok)
	assert.Equal(t, true, be.IsCoreStarted())

	block := tests.MakeBlockWithSeal(be, genesisHeader)
	header := block.Header()
	signer, err := be.Author(header)
	assert.NoError(t, err)
	assert.Equal(t, be.Address(), signer)
}

// TestPrepare
func TestPrepare(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be, ok := mustCreateAndStartNewBackend(nodePK, genesisHeader)
	assert.True(t, ok)
	assert.Equal(t, true, be.IsCoreStarted())

	block := tests.MakeBlockWithoutSeal(genesisHeader)
	header := block.Header()

	err = be.Prepare(be.Chain(), header)
	assert.NoError(t, err)

	header.ParentHash = common.HexToHash("1234567890")
	err = be.Prepare(be.Chain(), header)
	assert.Equal(t, consensus.ErrUnknownAncestor, err)
}

// TestVerifySeal
func TestVerifySeal(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be, ok := mustCreateAndStartNewBackend(nodePK, genesisHeader)
	assert.True(t, ok)
	assert.Equal(t, true, be.IsCoreStarted())

	// cannot verify genesis
	err = be.VerifySeal(be.Chain(), genesisHeader)
	assert.Equal(t, errUnknownBlock, err)

	block := tests.MakeBlockWithSeal(be, genesisHeader)
	err = be.VerifySeal(be.Chain(), block.Header())
	assert.NoError(t, err)
}

// TestPrepareExtra
// 0xd8c094000000000000000000000000000000000000000080c0
func TestPrepareExtra(t *testing.T) {
	vanity := make([]byte, types.TendermintExtraVanity)
	data := hexutil.MustDecode("0xd8c094000000000000000000000000000000000000000080c0")
	expectedResult := append(vanity, data...)

	header := &types.Header{
		Extra: vanity,
	}

	payload, err := prepareExtra(header)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, payload)

	// append useless information to extra-data
	header.Extra = append(vanity, make([]byte, 15)...)

	payload, err = prepareExtra(header)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, payload)
}
