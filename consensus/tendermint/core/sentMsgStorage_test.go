package core

import (
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCore_LookupSentMsg(t *testing.T) {
	var (
		nodePrivateKey = tests_utils.MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	//create New test backend and newMockChain
	be, txPool := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig, txPool)
	require.NoError(t, core.Start())
	var (
		step     = RoundStepPropose
		round    = int64(0)
		proposal = core.getDefaultProposal(core.getLogger(), round)
	)

	msgData, _ := rlp.EncodeToBytes(proposal)
	core.storeSentMsg(step, round, msgData)
	index := core.lookupSentMsg(step, round)
	assert.Equal(t, int64(0), index)
}
