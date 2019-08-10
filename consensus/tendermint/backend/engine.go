package backend

import (
	"time"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	now = time.Now
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

// Start implements consensus.Istanbul.Start
func (sb *backend) Start() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return tendermint.ErrStartedEngine
	}

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}
