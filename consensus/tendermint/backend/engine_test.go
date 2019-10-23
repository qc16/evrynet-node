package backend

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/common/hexutil"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/crypto/secp256k1"
	"github.com/evrynet-official/evrynet-client/ethdb"
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

func newTestEngine() *backend {
	nodeKey, _ := crypto.GenerateKey()
	be, _ := New(tendermint.DefaultConfig, nodeKey).(*backend)
	be.address = crypto.PubkeyToAddress(nodeKey.PublicKey)
	be.db = ethdb.NewMemDatabase()
	return be
}

// TestPrepareExtra
// 0xd8c094000000000000000000000000000000000000000080c0
func TestPrepareExtra(t *testing.T) {
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

	vanity := make([]byte, types.TendermintExtraVanity)
	data := hexutil.MustDecode("0xd8c094000000000000000000000000000000000000000080c0")
	expectedResult := append(vanity, data...)

	header := &types.Header{
		Extra:  vanity,
		Number: big.NewInt(0),
	}

	header.Extra = engine.prepareExtra(header)
	assert.Equal(t, expectedResult, header.Extra)

	// append useless information to extra-data
	header.Extra = append(vanity, make([]byte, 15)...)

	header.Extra = engine.prepareExtra(header)
	assert.Equal(t, expectedResult, header.Extra)

	var (
		candidate = ProposalValidator{
			address: common.HexToAddress("123456"),
			vote:    true,
		}
		newCandidate = ProposalValidator{
			address: common.HexToAddress("654321"),
			vote:    true,
		}
	)

	// will attach a candidate to voting
	engine.proposedValidator.setProposedValidator(candidate.address, candidate.vote)
	header.Extra = engine.prepareExtra(header)
	candidateAddr, _ := getModifiedValidator(*header)
	assert.Equal(t, candidate.address, candidateAddr)

	// the candidate will be repplaced by new candidate when call setProposedValidator and old candidate have not processed yet
	engine.proposedValidator.setProposedValidator(newCandidate.address, newCandidate.vote)
	header.Extra = engine.prepareExtra(header)
	newCandidateAddr, _ := getModifiedValidator(*header)
	assert.NotEqual(t, candidate.address, newCandidateAddr)
	assert.Equal(t, newCandidate.address, newCandidateAddr)
}
