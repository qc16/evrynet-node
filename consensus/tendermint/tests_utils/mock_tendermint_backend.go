package tests_utils

import (
	"errors"
	"testing"

	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/params"

	"crypto/ecdsa"
	"math/big"
	"sync"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/log"
)

type MockBackend struct {
	config             *tendermint.Config
	tendermintEventMux *event.TypeMux
	privateKey         *ecdsa.PrivateKey
	address            common.Address
	validators         []common.Address

	//once voting finish, the block will be send for commit here
	//it is a map of blocknumber- channels with mutex
	mutex *sync.RWMutex
	chain consensus.ChainReader

	//storingMsgs is used to store msg to handler when core stopped

	currentBlock func() *types.Block
}

func NewMockBackend(privateKey *ecdsa.PrivateKey, blockchain *MockChainReader, validators []common.Address) tendermint.Backend {
	return &MockBackend{
		config:             tendermint.DefaultConfig,
		tendermintEventMux: new(event.TypeMux),
		privateKey:         privateKey,
		address:            crypto.PubkeyToAddress(privateKey.PublicKey),
		mutex:              &sync.RWMutex{},
		chain:              blockchain,
		currentBlock:       blockchain.CurrentBlock,
		validators:         validators,
	}
}

func (mb *MockBackend) VerifyProposalHeader(header *types.Header) error {
	log.Warn("mocked backend always return nil when verifyProposalHeader")
	return nil
}

// EventMux implements tendermint.Backend.EventMux
func (mb *MockBackend) EventMux() *event.TypeMux {
	return mb.tendermintEventMux
}

// Sign implements tendermint.Backend.Sign
func (mb *MockBackend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, mb.privateKey)
}

// Address implements tendermint.Backend.Address
func (mb *MockBackend) Address() common.Address {
	return mb.address
}

// Broadcast implements tendermint.Backend.Broadcast
// It sends message to its validator by calling gossiping, and send message to itself by eventMux
func (mb *MockBackend) Broadcast(valSet tendermint.ValidatorSet, blockNumber *big.Int, payload []byte) error {
	// send to others
	if err := mb.Gossip(valSet, blockNumber, payload); err != nil {
		return err
	}
	// send to self
	go func() {
		if err := mb.EventMux().Post(tendermint.MessageEvent{
			Payload: payload,
		}); err != nil {
			log.Error("failed to post event to self", "error", err)
		}
	}()
	return nil
}

// Gossip implements tendermint.Backend.Gossip
// It sends message to its validators only, not itself.
// The validators must be able to connected through Peer.
// It will return MockBackend.ErrNoBroadcaster if no broadcaster is set for MockBackend
func (mb *MockBackend) Gossip(valSet tendermint.ValidatorSet, _ *big.Int, payload []byte) error {
	return errors.New("not implemented")
}

// Multicast implements tendermint.Backend.Multicast
func (mb *MockBackend) Multicast(targets map[common.Address]bool, payload []byte) error {
	panic("implement me")
}

// Validators return validator set for a block number
// TODO: revise this function once auth vote is implemented
func (mb *MockBackend) Validators(blockNumber *big.Int) tendermint.ValidatorSet {
	log.Error("not implemented")
	return validator.NewSet(mb.validators, mb.config.ProposerPolicy, int64(0))
}

// FindExistingPeers check validator peers exist or not by address
func (mb *MockBackend) FindExistingPeers(valSet tendermint.ValidatorSet) map[common.Address]consensus.Peer {
	log.Error("not implemented")
	return make(map[common.Address]consensus.Peer)
}

//Commit implement tendermint.Backend.Commit()
func (mb *MockBackend) Commit(block *types.Block) {
	log.Error("not implemented")
}

// EnqueueBlock adds a block returned from consensus into fetcher queue
func (mb *MockBackend) EnqueueBlock(block *types.Block) {
	log.Error("not implemented")
}

func (mb *MockBackend) CurrentHeadBlock() *types.Block {
	return mb.currentBlock()
}

func (mb *MockBackend) Cancel(block *types.Block) {
	log.Error("not implemented")
}

// ValidatorsByChainReader returns val-set from snapshot
func (mb *MockBackend) ValidatorsByChainReader(blockNumber *big.Int, chain consensus.ChainReader) tendermint.ValidatorSet {
	log.Error("not implemented")
	return validator.NewSet(nil, 0, int64(0))
}

func MustCreateAndStartNewBackend(t *testing.T, nodePrivateKey *ecdsa.PrivateKey, genesisHeader *types.Header, validators []common.Address) (tendermint.Backend, *core.TxPool) {
	var (
		address = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		trigger = false
		statedb = MustCreateStateDB(t)

		testTxPoolConfig core.TxPoolConfig
		blockchain       = &MockChainReader{
			GenesisHeader: genesisHeader,
			MockBlockChain: &MockBlockChain{
				Statedb:       statedb,
				GasLimit:      1000000000,
				ChainHeadFeed: new(event.Feed),
			},
			Address: address,
			Trigger: &trigger,
		}
		pool = core.NewTxPool(testTxPoolConfig, params.TendermintTestChainConfig, blockchain)
		be   = NewMockBackend(nodePrivateKey, blockchain, validators)
	)
	statedb.SetBalance(address, new(big.Int).SetUint64(params.Ether))
	return be, pool
}
