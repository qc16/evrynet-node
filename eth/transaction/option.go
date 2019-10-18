package transaction

import (
	"github.com/evrynet-official/evrynet-client/core"
)

// TxPoolOpts is txpool options which are used when creating consensus engine
type TxPoolOpts struct {
	CoreTxPool *core.TxPool
}
