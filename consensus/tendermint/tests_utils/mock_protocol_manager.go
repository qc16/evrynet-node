package tests_utils

import (
	"github.com/Evrynetlabs/evrynet-client/common"
	"github.com/Evrynetlabs/evrynet-client/consensus"
	"github.com/Evrynetlabs/evrynet-client/core/types"
)

type MockProtocolManager struct{}

// FindPeers retrives peers by addresses
func (pm *MockProtocolManager) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	return make(map[common.Address]consensus.Peer)
}

// Enqueue adds a block into fetcher queue
func (pm *MockProtocolManager) Enqueue(id string, block *types.Block) {}
