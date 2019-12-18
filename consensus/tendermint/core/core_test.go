package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Evrynetlabs/evrynet-client/common"
	"github.com/Evrynetlabs/evrynet-client/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/Evrynetlabs/evrynet-client/crypto"
)

func TestRecoverCoreTimeoutWithNewHeight(t *testing.T) {
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
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())

	state.UpdateRoundStep(0, RoundStepNewHeight)
	core.currentState = state
	require.NoError(t, core.Stop())
	time.Sleep(2 * time.Second)
	require.NoError(t, core.Start())

	//wait for 4 second to core's state jump from RoundStepPropose to RoundStepPrevote
	time.Sleep(1 * time.Second)
	assert.Equal(t, RoundStepPropose, core.CurrentState().Step())
}

func TestRecoverCoreTimeoutWithPropose(t *testing.T) {
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
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())

	state.UpdateRoundStep(0, RoundStepPropose)
	core.currentState = state
	require.NoError(t, core.Stop())
	time.Sleep(2 * time.Second)
	require.NoError(t, core.Start())

	//wait for 4 second to core's state jump from RoundStepPropose to RoundStepPrevote
	time.Sleep(tendermint.DefaultConfig.TimeoutPropose + (1 * time.Second))
	assert.Equal(t, RoundStepPrevote, core.CurrentState().Step())
}

func TestRecoverCoreTimeoutWithPrevoteWait(t *testing.T) {
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
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())

	state.UpdateRoundStep(0, RoundStepPrevoteWait)
	require.NoError(t, core.Stop())
	time.Sleep(2 * time.Second)
	require.NoError(t, core.Start())
	assert.Equal(t, RoundStepPrevoteWait, core.CurrentState().Step())

	//wait for 4 second to core's state jump from RoundStepPrevoteWait to RoundStepPrecommit
	time.Sleep(tendermint.DefaultConfig.TimeoutPrevote +  (500 * time.Millisecond))
	assert.Equal(t, RoundStepPrecommit, core.CurrentState().Step())
}

func TestRecoverCoreTimeoutWithPreCommit(t *testing.T) {
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
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())

	state.UpdateRoundStep(0, RoundStepPrecommitWait)
	require.NoError(t, core.Stop())
	time.Sleep(2 * time.Second)
	require.NoError(t, core.Start())
	assert.Equal(t, RoundStepPrecommitWait, core.CurrentState().Step())

	//wait for core's state jump from RoundStepPrecommit to RoundStepPropose
	time.Sleep(tendermint.DefaultConfig.TimeoutPrecommit + (1 * time.Second))
	assert.Equal(t, RoundStepPropose, core.CurrentState().Step())
}
