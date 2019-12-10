package core

import (
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/stretchr/testify/require"
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
	time.Sleep(tendermint.DefaultConfig.TimeoutPropose + (1 * time.Second))
	proposalMsg, err := core.LookupSentMsg(RoundStepPropose, 0)
	require.NoError(t, err)
	require.NotNil(t, proposalMsg)
	require.NotNil(t, proposalMsg.Data)

	time.Sleep(tendermint.DefaultConfig.TimeoutPrevote + (1 * time.Second))
	prevoteMsg, err := core.LookupSentMsg(RoundStepPrevote, 0)
	require.NoError(t, err)
	require.NotNil(t, prevoteMsg)
	require.NotNil(t, proposalMsg.Data)
}
