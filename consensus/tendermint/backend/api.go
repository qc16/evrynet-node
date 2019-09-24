package backend


import (
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
)

// API is a user facing RPC API to dump Istanbul state
type API struct {
	chain    consensus.ChainReader
	tendermint *backend
}

// Propose injects a new authorization candidate that the validator will attempt to
// push through.
func (api *API) Propose(address common.Address, auth bool) {
	api.tendermint.candidatesLock.Lock()
	defer api.tendermint.candidatesLock.Unlock()

	api.tendermint.candidates[address] = auth
}

// Discard drops a currently running candidate, stopping the validator from casting
// further votes (either for or against).
func (api *API) Discard(address common.Address) {
	api.tendermint.candidatesLock.Lock()
	defer api.tendermint.candidatesLock.Unlock()

	delete(api.tendermint.candidates, address)
}
