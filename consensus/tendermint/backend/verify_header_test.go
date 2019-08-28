package backend

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
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
	genesisHeader = types.Header{
		Number:     big.NewInt(int64(0)),
		ParentHash: common.HexToHash("0x01"),
		UncleHash:  types.CalcUncleHash(nil),
		Root:       common.HexToHash("0x0"),
		Extra:      makeExtra(true),
		Difficulty: defaultDifficulty,
	}
	genesisBlock = types.NewBlock(&genesisHeader, nil, nil, nil)
)

func makeNodeKey() *ecdsa.PrivateKey {
	key, _ := generatePrivateKey()
	return key
}

func TestBackend_VerifyHeader(t *testing.T) {
	//create New test backend and newMockChain
	chain, engine := makeBlockChain()
	assert.NotNil(t, chain)
	assert.NotNil(t, engine)
	assert.Equal(t, true, engine.coreStarted)

	// without seal
	block := makeBlockWithoutSeal(chain)
	err := engine.VerifyHeader(chain, block.Header(), false)
	assert.Error(t, err)
	assert.Equal(t, "invalid signature length", err.Error())

	// with seal but incorrect coinbase
	block = makeBlockWithSeal(chain)
	assert.NotNil(t, chain)
	header := block.Header()
	header.Coinbase = common.Address{}
	appendSeal(header)
	err = engine.VerifyHeader(chain, header, false)
	assert.Equal(t, "invalid coin base address", err.Error())

	// without committed seal
	block = makeBlockWithSeal(chain)
	assert.NotNil(t, chain)
	err = engine.VerifyHeader(chain, block.Header(), false)
	assert.Equal(t, "zero committed seals", err.Error())

	// with committed seal but is invalid
	block = makeBlockWithCommittedSeal(chain)
	assert.NotNil(t, chain)
	err = engine.VerifyHeader(chain, block.Header(), false)
	assert.Equal(t, "invalid signature", err.Error())
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

func makeBlockWithoutSeal(chain *mockChain) *types.Block {
	var header = makeHeader(chain.Genesis())
	prepare(chain, header)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

func makeBlockWithSeal(chain *mockChain) *types.Block {
	var header = makeHeader(chain.Genesis())
	prepare(chain, header)
	appendSeal(header)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

func makeBlockWithCommittedSeal(chain *mockChain) *types.Block {
	var header = makeHeader(chain.Genesis())
	prepare(chain, header)
	appendCommittedSeal(header)
	appendSeal(header)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

func appendSeal(header *types.Header) {
	// sign the hash
	hash := sigHash(header)
	seal := sign(hash.Bytes())

	// add seal to extradata of the header
	extra, _ := types.ExtractTendermintExtra(header)
	extra.Seal = seal
	payload, _ := rlp.EncodeToBytes(&extra)
	header.Extra = append(header.Extra[:types.TendermintExtraVanity], payload...)
}

func appendCommittedSeal(header *types.Header) {
	var (
		//messages = map[common.Address]*commitMessage
		addr = getAddress()
	)

	// msg := &commitMessage{
	// 	Code: msgCommit,
	// 	Msg: []byte{}
	// 	Address: addr,
	// 	Signature: []byte{},
	// 	CommittedSeal: addr.Bytes(),
	// }
	// messages[addr] = msg

	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], addr.Bytes()[:])
	// add seal to extradata of the header
	extra, _ := types.ExtractTendermintExtra(header)
	extra.CommittedSeal = committedSeals
	payload, _ := rlp.EncodeToBytes(&extra)
	header.Extra = append(header.Extra[:types.TendermintExtraVanity], payload...)
}

func prepare(chain consensus.ChainReader, header *types.Header) {
	header.Coinbase = getAddress()
	header.Nonce = emptyNonce
	header.MixDigest = types.TendermintDigest

	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// add validators in snapshot to extraData's validators section
	extra := makeExtra(true)
	header.Extra = extra
}

func makeHeader(parent *types.Block) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Extra:      makeExtra(false),
		Difficulty: defaultDifficulty,
	}
	return header
}

func makeExtra(importValidators bool) []byte {
	var extra []byte
	extra = append(extra, bytes.Repeat([]byte{0x00}, types.TendermintExtraVanity)...)
	extra = extra[:types.TendermintExtraVanity]

	data := &types.TendermintExtra{
		Validators:    []common.Address{},
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}
	if importValidators {
		data.Validators = validators
	}

	payload, _ := rlp.EncodeToBytes(&data)
	extra = append(extra, payload...)
	return extra
}

func sign(data []byte) []byte {
	hashData := crypto.Keccak256(data)
	seal, _ := crypto.Sign(hashData, nodeKey)
	return seal
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
