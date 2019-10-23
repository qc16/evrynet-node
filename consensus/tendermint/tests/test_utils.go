package tests

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"log"
	"math/big"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/eth/transaction"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/p2p"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
)

func mustGeneratePrivateKey(key string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		panic(err)
	}
	return privateKey
}

//TestChain is mock struct for chain
type TestChain struct {
	GenesisHeader *types.Header
	*TestBlockChain
	Address common.Address
	Trigger *bool
}

func (c *TestChain) Config() *params.ChainConfig {
	panic("implement me")
}

func (c *TestChain) CurrentHeader() *types.Header {
	return c.CurrentBlock().Header()
}

//GetHeader implement a mock version of chainReader.GetHeader
//It returns correctParentHeader with Number set to input blockNumber
func (c *TestChain) GetHeader(hash common.Hash, blockNumber uint64) *types.Header {
	if c.GenesisHeader.Hash() == hash {
		return c.GenesisHeader
	}
	return nil
}

//GetHeaderByNumber implement a mock version of chainReader.GetHeaderByNumber
//It returns genesis Header if blockNumber is 0, else return an empty Header
func (c *TestChain) GetHeaderByNumber(blockNumber uint64) *types.Header {
	return c.GenesisHeader
}

func (c *TestChain) GetHeaderByHash(hash common.Hash) *types.Header {
	panic("implement me")
}

//State is used multiple times to reset the pending state.
// when simulate is true it will create a state that indicates
// that tx0 and tx1 are included in the chain.
func (c *TestChain) State() (*state.StateDB, error) {
	// delay "state change" by one. The tx pool fetches the
	// state multiple times and by delaying it a bit we simulate
	// a state change between those fetches.
	stdb := c.Statedb
	if *c.Trigger {
		c.Statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
		c.Statedb.SetNonce(c.Address, 0)
		c.Statedb.SetBalance(c.Address, new(big.Int).SetUint64(params.Ether))
		*c.Trigger = false
	}
	return stdb, nil
}

//TestBlockChain is mock struct for block chain
type TestBlockChain struct {
	Statedb       *state.StateDB
	GasLimit      uint64
	ChainHeadFeed *event.Feed
}

func (bc *TestBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.GasLimit,
	}, nil, nil, nil)
}

func (bc *TestBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *TestBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.Statedb, nil
}

func (bc *TestBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.ChainHeadFeed.Subscribe(ch)
}

// ------------------------------------
type testProtocolManager struct{}

// FindPeers retrives peers by addresses
func (pm *testProtocolManager) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	return make(map[common.Address]consensus.Peer)
}

// Enqueue adds a block into fetcher queue
func (pm *testProtocolManager) Enqueue(id string, block *types.Block) {
	return
}

// ------------------------------------

//TestBackend is mock interface of tendermint.Backend
type TestBackend interface {
	Address() common.Address
	EventMux() *event.TypeMux
	Sign([]byte) ([]byte, error)
	Gossip(valSet tendermint.ValidatorSet, payload []byte) error
	Broadcast(valSet tendermint.ValidatorSet, payload []byte) error
	Validators(blockNumber *big.Int) tendermint.ValidatorSet
	CurrentHeadBlock() *types.Block
	FindExistingPeers(targets tendermint.ValidatorSet) map[common.Address]consensus.Peer
	Commit(block *types.Block)
	EnqueueBlock(block *types.Block)
	ClearStoringMsg()
	Verify(tendermint.Proposal) error
	VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error
	TxPool() *transaction.TxPoolOpts
	Chain() consensus.ChainReader
	Start(chain consensus.ChainReader, currentBlock func() *types.Block) error
	SetBroadcaster(broadcaster consensus.Broadcaster)
	IsCoreStarted() bool
	Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error)
	Author(header *types.Header) (common.Address, error)
	Prepare(chain consensus.ChainReader, header *types.Header) error
	VerifySeal(chain consensus.ChainReader, header *types.Header) error
	HandleMsg(addr common.Address, msg p2p.Msg) (bool, error)
	Core() tendermintCore.Engine
}

// ------------------------------------

//TestEngine is mock interface of tendermint.core.Engine
type TestEngine interface {
	Start() error
	Stop() error
	SetBlockForProposal(block *types.Block)
	VerifyProposal(proposal tendermint.Proposal, msg tendermintCore.Message) error
}

// ------------------------------------

func MakeNodeKey() *ecdsa.PrivateKey {
	key, _ := GeneratePrivateKey()
	return key
}

func MustStartTestChainAndBackend(be TestBackend, blockchain *TestChain) bool {
	be.SetBroadcaster(&testProtocolManager{})
	if err := be.Start(blockchain, blockchain.CurrentBlock); err != nil {
		log.Panicf("cannot start backend, error:%v", err)
		return false
	}
	return true
}

func MakeGenesisHeader(validators []common.Address) *types.Header {
	var header = &types.Header{
		Number:     big.NewInt(int64(0)),
		ParentHash: common.HexToHash("0x01"),
		UncleHash:  types.CalcUncleHash(nil),
		Root:       common.HexToHash("0x0"),
		Difficulty: big.NewInt(1),
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

func MakeBlockWithoutSeal(pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	return types.NewBlockWithHeader(header)
}

func MakeBlockWithSeal(be TestBackend, pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	return types.NewBlockWithHeader(header)
}

func MustMakeBlockWithCommittedSealInvalid(be TestBackend, pHeader *types.Header) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	invalidCommitSeal := make([]byte, types.TendermintExtraSeal)
	_, err := rand.Read(invalidCommitSeal)
	if err != nil {
		panic(err)
	}
	appendCommittedSeal(header, invalidCommitSeal)
	return types.NewBlockWithHeader(header)
}

func MustMakeBlockWithCommittedSeal(be TestBackend, pHeader *types.Header, validators []common.Address) *types.Block {
	header := makeHeaderFromParent(types.NewBlockWithHeader(pHeader))
	AppendSeal(header, be)
	commitHash := tendermintCore.PrepareCommittedSeal(header.Hash())
	committedSeal, err := be.Sign(commitHash)
	if err != nil {
		panic(err)
	}
	appendCommittedSeal(header, committedSeal)
	block := types.NewBlockWithHeader(header)
	return block.WithSeal(header)
}

//AppendSeal sign the header with the engine's key and write the seal to the input header's extra data
func AppendSeal(header *types.Header, be tendermint.Backend) {
	// sign the hash
	seal, _ := be.Sign(utils.SigHash(header).Bytes())
	utils.WriteSeal(header, seal)
}

//appendCommittedSeal
func appendCommittedSeal(header *types.Header, committedSeal []byte) {
	//TODO: make this logic as the same as AppendSeal, which involve signing commit before writeCommittedSeal
	committedSeals := make([][]byte, 1)
	committedSeals[0] = make([]byte, types.TendermintExtraSeal)
	copy(committedSeals[0][:], committedSeal[:])
	utils.WriteCommittedSeals(header, committedSeals)
}

//makeHeaderFromParent return a new block With valid information from its parents.
func makeHeaderFromParent(parent *types.Block) *types.Header {
	header := &types.Header{
		Coinbase:   GetAddress(),
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Difficulty: big.NewInt(1),
		MixDigest:  types.TendermintDigest,
	}
	extra, _ := prepareExtra(header)
	header.Extra = extra
	return header
}

func GetAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough TendermintExtraVanity bytes.
	if len(header.Extra) < types.TendermintExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.TendermintExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.TendermintExtraVanity])

	tdm := &types.TendermintExtra{}
	payload, err := rlp.EncodeToBytes(&tdm)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}
