package tests_utils

import (
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/core/types"
)

type MockProtocolManager struct{}

// FindPeers retrives peers by addresses
func (pm *MockProtocolManager) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	return make(map[common.Address]consensus.Peer)
}

// Enqueue adds a block into fetcher queue
func (pm *MockProtocolManager) Enqueue(id string, block *types.Block) {
	return
}
