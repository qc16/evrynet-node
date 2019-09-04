package tendermint

import (
	"github.com/evrynet-official/evrynet-client/core/types"
)

// NewBlockEvent is the event sent from Backend to Core after engine.Seal() is called.
// It included the latest eligible block from tx_pool
type NewBlockEvent struct {
	Block *types.Block
}

// MessageEvent is posted for Tendermint engine communication
type MessageEvent struct {
	Payload []byte
}
