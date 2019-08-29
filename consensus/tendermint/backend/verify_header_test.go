package backend

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/stretchr/testify/assert"
)

const (
	msgCommit uint64 = iota
)

var (
	emptyNonce = types.BlockNonce{}
	nodeKey    = makeNodeKey()
	validators = []common.Address{
		getAddress(),
	}
	genesisHeader = makeGenesisHeader()
	genesisBlock = types.NewBlockWithHeader(genesisHeader)
)

func TestBackend_VerifyHeader(t *testing.T) {
	//create New test backend and newMockChain
	chain, engine := makeBlockChain()
	assert.NotNil(t, chain)
	assert.NotNil(t, engine)
	assert.Equal(t, true, engine.coreStarted)

	// without seal
	block := makeBlockWithoutSeal(chain, engine)
	err := engine.VerifyHeader(chain, block.Header(), false)
	assert.Error(t, err)
	assert.Equal(t, "invalid signature length", err.Error())

	// with seal but incorrect coinbase
	block = makeBlockWithSeal(chain, engine)
	assert.NotNil(t, chain)
	header := block.Header()
	header.Coinbase = common.Address{}
	appendSeal(header, engine)
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, "invalid coin base address", err.Error())

	// without committed seal
	block = makeBlockWithSeal(chain, engine)
	assert.NotNil(t, chain)
	err = engine.VerifyHeader(chain, block.Header(), false)
	assert.Equal(t, "zero committed seals", err.Error())

	// with committed seal but is invalid
	block = makeBlockWithCommittedSealInvalid(chain, engine)
	assert.NotNil(t, chain)
	err = engine.VerifyHeader(chain, block.Header(), false)
	assert.Equal(t, "invalid signature", err.Error())

	// with committed seal
	block = makeBlockWithCommittedSeal(chain, engine)
	assert.NotNil(t, chain)
	err = engine.VerifyHeader(chain, block.Header(), false)
	assert.NoError(t, err)
}

func makeGenesisHeader() *types.Header {
	var header = &types.Header{
		Number:     big.NewInt(int64(0)),
		ParentHash: common.HexToHash("0x01"),
		UncleHash:  types.CalcUncleHash(nil),
		Root:       common.HexToHash("0x0"),
		Difficulty: defaultDifficulty,
	}
	extra, _ := prepareExtra(header, validators)
	header.Extra = extra
	return header
}

func makeNodeKey() *ecdsa.PrivateKey {
	key, _ := generatePrivateKey()
	return key
}

func makeBlockChain() (*mockChain, *backend) {
	memDB := ethdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	b, _ := New(config, nodeKey).(*backend)
	b.db = memDB
	b.privateKey = nodeKey
	b.address = getAddress()

	chain := mockChain{}
	b.Start(&chain, nil)
	return &chain, b
}

func makeBlockWithoutSeal(chain *mockChain, engine *backend) *types.Block {
	var header = makeHeader(chain.Genesis())
	block := types.NewBlockWithHeader(header)
	return block
}

func makeBlockWithSeal(chain *mockChain, engine *backend) *types.Block {
	var header = makeHeader(chain.Genesis())
	appendSeal(header, engine)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

func makeBlockWithCommittedSealInvalid(chain *mockChain, engine *backend) *types.Block {
	var header = makeHeader(chain.Genesis())
	appendSeal(header, engine)
	appendCommittedSeal(header, engine, []byte{})
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

func makeBlockWithCommittedSeal(chain *mockChain, engine *backend) *types.Block {
	var header = makeHeader(chain.Genesis())
	appendSeal(header, engine)

	commitHash := tendermintCore.PrepareCommittedSeal(header.Hash())
	committedSeal, _ := engine.Sign(commitHash)

	appendCommittedSeal(header, engine, committedSeal)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

func appendSeal(header *types.Header, engine *backend) {
	// sign the hash
	seal, _ := engine.Sign(sigHash(header).Bytes())
	writeSeal(header, seal)
}

func appendCommittedSeal(header *types.Header, engine *backend, committedSeal []byte) {
	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], committedSeal[:])
	writeCommittedSeals(header, committedSeals)
}

func makeHeader(parent *types.Block) *types.Header {
	header := &types.Header{
		Coinbase: getAddress(),
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Difficulty: defaultDifficulty,
	}
	extra, _ := prepareExtra(header, validators)
	header.Extra = extra
	return header
}

type commitMessage struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

type message struct {
	messages map[common.Address]*commitMessage
}
type mockChain struct {
}

//GetHeader implement a mock version of chainReader.GetHeader
//It returns correctParentHeader with Number set to input blockNumber
func (mc *mockChain) GetHeader(hash common.Hash, blockNumber uint64) *types.Header {
	if blockNumber == 0 {
		return mc.Genesis().Header()
	}

	return nil
}

//GetHeaderByNumber implement a mock version of chainReader.GetHeaderByNumber
//It returns genesis Header if blockNumber is 0, else return an empty Header
func (mc *mockChain) GetHeaderByNumber(blockNumber uint64) *types.Header {
	if blockNumber == 0 {
		return mc.Genesis().Header()
	}
	return nil
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

func (mc *mockChain) Genesis() *types.Block {
	return genesisBlock
}
