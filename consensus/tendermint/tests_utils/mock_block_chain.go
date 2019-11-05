package tests_utils

import (
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/event"
)

//MockBlockChain is mock struct for block chain
type MockBlockChain struct {
	Statedb       *state.StateDB
	GasLimit      uint64
	ChainHeadFeed *event.Feed
}

func (bc *MockBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.GasLimit,
	}, nil, nil, nil)
}

func (bc *MockBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *MockBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.Statedb, nil
}

func (bc *MockBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.ChainHeadFeed.Subscribe(ch)
}

