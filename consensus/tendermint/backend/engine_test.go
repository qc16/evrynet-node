package backend

import (
	"math/big"
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/common/hexutil"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/crypto/secp256k1"
	"github.com/evrynet-official/evrynet-client/ethdb"
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
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be := mustCreateAndStartNewBackend(t, nodePK, genesisHeader)

	// without seal
	block := tests_utils.MakeBlockWithoutSeal(genesisHeader)
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, be.VerifyHeader(be.chain, block.Header(), false))

	err = be.Seal(be.chain, block, nil, nil)

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
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be := mustCreateAndStartNewBackend(t, nodePK, genesisHeader)

	block := tests_utils.MakeBlockWithSeal(be, genesisHeader)
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
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be := mustCreateAndStartNewBackend(t, nodePK, genesisHeader)

	block := tests_utils.MakeBlockWithoutSeal(genesisHeader)
	header := block.Header()

	err = be.Prepare(be.chain, header)
	assert.NoError(t, err)

	header.ParentHash = common.HexToHash("1234567890")
	err = be.Prepare(be.chain, header)
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
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	be := mustCreateAndStartNewBackend(t, nodePK, genesisHeader)

	// cannot verify genesis
	err = be.VerifySeal(be.chain, genesisHeader)
	assert.Equal(t, errUnknownBlock, err)

	block := tests_utils.MakeBlockWithSeal(be, genesisHeader)
	err = be.VerifySeal(be.chain, block.Header())
	assert.NoError(t, err)
}

func newTestEngine() *Backend {
	nodeKey, _ := crypto.GenerateKey()
	be, _ := New(tendermint.DefaultConfig, nodeKey, nil).(*Backend)
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
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
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

	header.Extra, err = tests_utils.PrepareExtra(header)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, header.Extra)

	// append useless information to extra-data
	header.Extra = append(vanity, make([]byte, 15)...)

	header.Extra, err = tests_utils.PrepareExtra(header)
	assert.Nil(t, err)
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
	err = engine.proposedValidator.setProposedValidator(candidate.address, candidate.vote)
	assert.Nil(t, err)
	header.Extra, err = tests_utils.PrepareExtra(header)
	assert.Nil(t, err)
	candidateAddr, err := getModifiedValidator(*header)
	assert.Nil(t, err)
	assert.Equal(t, candidate.address, candidateAddr)

	// the candidate will be repplaced by new candidate when call setProposedValidator and old candidate have not processed yet
	err = engine.proposedValidator.setProposedValidator(newCandidate.address, newCandidate.vote)
	assert.Nil(t, err)
	header.Extra, err = tests_utils.PrepareExtra(header)
	assert.Nil(t, err)
	newCandidateAddr, err := getModifiedValidator(*header)
	assert.Nil(t, err)
	assert.NotEqual(t, candidate.address, newCandidateAddr)
	assert.Equal(t, newCandidate.address, newCandidateAddr)
}
