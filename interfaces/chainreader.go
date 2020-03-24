package interfaces

import (
	"github.com/Evrynetlabs/evrynet-node/core/types"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header
}
