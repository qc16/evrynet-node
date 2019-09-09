package backend

import (
	"bytes"
	"crypto/ecdsa"
	"log"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/crypto/secp256k1"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
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

func makeGenesisHeader(validators []common.Address) *types.Header {

	var header = &types.Header{
		Number:     big.NewInt(int64(0)),
		ParentHash: common.HexToHash("0x01"),
		UncleHash:  types.CalcUncleHash(nil),
		Root:       common.HexToHash("0x0"),
		Difficulty: defaultDifficulty,
		MixDigest:  types.TendermintDigest,
	}
	extra, _ := prepareExtra(header)

	var buf bytes.Buffer
	buf.Write(extra[:types.TendermintExtraVanity])
	tdm := &types.TendermintExtra{
		Validators:    validators,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}
	payload, _ := rlp.EncodeToBytes(&tdm)

	header.Extra = append(buf.Bytes(), payload...)
	return header
}

func makeNodeKey() *ecdsa.PrivateKey {
	key, _ := generatePrivateKey()
	return key
}

func mustStartTestChainAndBackend(nodePK *ecdsa.PrivateKey, genesisHeader *types.Header) (*mockChain, *backend) {
	memDB := ethdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	b, ok := New(config, nodePK, WithDB(memDB)).(*backend)
	if !ok {
		panic("New() cannot be asserted back to backend")
	}
	chain := mockChain{
		genesisHeader: genesisHeader,
	}

	snap, err := b.snapshot(&chain, 0, common.Hash{}, nil)
	if err != nil {
		panic(err)
	}
	if snap == nil {
		panic("failed to get snapshot")
	}

	if err := b.Start(&chain, nil); err != nil {
		log.Panicf("cannot start backend, error:%v", err)
	}
	return &chain, b
}

func makeBlockWithoutSeal(pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	return types.NewBlockWithHeader(header)
}

func makeBlockWithSeal(engine *backend, pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	appendSeal(header, engine)
	return types.NewBlockWithHeader(header)
}

func mustMakeBlockWithCommittedSealInvalid(engine *backend, pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	appendSeal(header, engine)
	invalidCommitSeal := make([]byte, types.TendermintExtraSeal)
	_, err := rand.Read(invalidCommitSeal)
	if err != nil {
		panic(err)
	}
	appendCommittedSeal(header, invalidCommitSeal)
	return types.NewBlockWithHeader(header)
}

func mustMakeBlockWithCommittedSeal(engine *backend, pHeader *types.Header, validators []common.Address) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	appendSeal(header, engine)
	commitHash := tendermintCore.PrepareCommittedSeal(header.Hash())
	committedSeal, err := engine.Sign(commitHash)
	if err != nil {
		panic(err)
	}
	appendCommittedSeal(header, committedSeal)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

//appendSeal sign the header with the engine's key and write the seal to the input header's extra data
func appendSeal(header *types.Header, engine *backend) {
	// sign the hash
	seal, _ := engine.Sign(sigHash(header).Bytes())
	writeSeal(header, seal)
}

//appendCommittedSeal
func appendCommittedSeal(header *types.Header, committedSeal []byte) {
	//TODO: make this logic as the same as appendSeal, which involve signing commit before writeCommittedSeal
	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], committedSeal[:])
	writeCommittedSeals(header, committedSeals)
}

//makeHeaderFromParent return a new block With valid information from its parents.
func makeHeaderFromParent(parent *types.Block) *types.Header {
	header := &types.Header{
		Coinbase:   getAddress(),
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Difficulty: defaultDifficulty,
		MixDigest:  types.TendermintDigest,
	}
	extra, _ := prepareExtra(header)
	header.Extra = extra
	return header
}

//mockChain implement consensus.ChainReader interface. It will return pseudo data for testing purposes.
type mockChain struct {
	genesisHeader *types.Header
}

//GetHeader implement a mock version of chainReader.GetHeader
//It returns correctParentHeader with Number set to input blockNumber
func (mc *mockChain) GetHeader(hash common.Hash, blockNumber uint64) *types.Header {
	if mc.genesisHeader.Hash() == hash {
		return mc.genesisHeader
	}

	return nil
}

//GetHeaderByNumber implement a mock version of chainReader.GetHeaderByNumber
//It returns genesis Header if blockNumber is 0, else return an empty Header
func (mc *mockChain) GetHeaderByNumber(blockNumber uint64) *types.Header {
	return mc.genesisHeader
}

func (mc *mockChain) Config() *params.ChainConfig {
	panic("implement me")
}

func (mc *mockChain) CurrentHeader() *types.Header {
	panic("implement me")
}

func (mc *mockChain) GetHeaderByHash(hash common.Hash) *types.Header {
	panic("implement me")
}

func (mc *mockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	panic("implement me")
}
