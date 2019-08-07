package backend

import (
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/core/types"
)

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {
	// post block into tendermint engine
	go sb.EventMux().Post(tendermint.RequestEvent{
		Proposal: block,
	})
	return nil
}
