package core

import (
	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/rlp"
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
		step        = RoundStepPropose
		round       = int64(0)
		blockNumber = uint64(1)
		proposal    = core.getDefaultProposal(core.getLogger(), round)
	)

	core.StoreSentMsg(step, round, proposal)
	proposalMsg, err := core.LookupSentMsg(step, round)
	require.NoError(t, err)
	require.NotNil(t, proposalMsg)
	require.Equal(t, blockNumber, proposalMsg.BlockNumber)
	require.NotNil(t, proposalMsg.Data)

	expectProposal, _ := rlp.EncodeToBytes(proposal)
	require.Equal(t, expectProposal, proposalMsg.Data)
}
