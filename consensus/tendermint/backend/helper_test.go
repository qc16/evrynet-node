package backend

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"log"
	"math/big"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/eth/transaction"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
)

type testChain struct {
	genesisHeader *types.Header
	*testBlockChain
	address common.Address
	trigger *bool
}

func (c *testChain) Config() *params.ChainConfig {
	panic("implement me")
}

func (c *testChain) CurrentHeader() *types.Header {
	return c.CurrentBlock().Header()
}

//GetHeader implement a mock version of chainReader.GetHeader
//It returns correctParentHeader with Number set to input blockNumber
func (c *testChain) GetHeader(hash common.Hash, blockNumber uint64) *types.Header {
	if c.genesisHeader.Hash() == hash {
		return c.genesisHeader
	}
	return nil
}

//GetHeaderByNumber implement a mock version of chainReader.GetHeaderByNumber
//It returns genesis Header if blockNumber is 0, else return an empty Header
func (c *testChain) GetHeaderByNumber(blockNumber uint64) *types.Header {
	return c.genesisHeader
}

func (c *testChain) GetHeaderByHash(hash common.Hash) *types.Header {
	panic("implement me")
}

// testChain.State() is used multiple times to reset the pending state.
// when simulate is true it will create a state that indicates
// that tx0 and tx1 are included in the chain.
func (c *testChain) State() (*state.StateDB, error) {
	// delay "state change" by one. The tx pool fetches the
	// state multiple times and by delaying it a bit we simulate
	// a state change between those fetches.
	stdb := c.statedb
	if *c.trigger {
		c.statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
		c.statedb.SetNonce(c.address, 0)
		c.statedb.SetBalance(c.address, new(big.Int).SetUint64(params.Ether))
		*c.trigger = false
	}
	return stdb, nil
}

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.gasLimit,
	}, nil, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

// ------------------------------------
func makeNodeKey() *ecdsa.PrivateKey {
	key, _ := generatePrivateKey()
	return key
}

func mustStartTestChainAndBackend(nodePK *ecdsa.PrivateKey, genesisHeader *types.Header) (*testChain, *backend) {
	address := crypto.PubkeyToAddress(nodePK.PublicKey)
	trigger := false
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	statedb.SetBalance(address, new(big.Int).SetUint64(params.Ether))
	var testTxPoolConfig core.TxPoolConfig

	blockchain := &testChain{genesisHeader, &testBlockChain{statedb, 1000000000, new(event.Feed)}, address, &trigger}
	pool := core.NewTxPool(testTxPoolConfig, params.TendermintTestChainConfig, blockchain)
	defer pool.Stop()

	memDB := ethdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	b, ok := New(config, nodePK, WithTxPoolOpts(&transaction.TxPoolOpts{CoreTxPool: pool}), WithDB(memDB)).(*backend)
	if !ok {
		panic("New() cannot be asserted back to backend")
	}

	snap, err := b.snapshot(blockchain, 0, common.Hash{}, nil)
	if err != nil {
		panic(err)
	}
	if snap == nil {
		panic("failed to get snapshot")
	}

	if err := b.Start(blockchain, nil); err != nil {
		log.Panicf("cannot start backend, error:%v", err)
	}
	return blockchain, b
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
	seal, _ := engine.Sign(utils.SigHash(header).Bytes())
	utils.WriteSeal(header, seal)
}

//appendCommittedSeal
func appendCommittedSeal(header *types.Header, committedSeal []byte) {
	//TODO: make this logic as the same as appendSeal, which involve signing commit before writeCommittedSeal
	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], committedSeal[:])
	utils.WriteCommittedSeals(header, committedSeals)
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
