package core

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
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
	be, _ := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig)
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
	be, _ := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig)
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
	be, _ := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig)
	require.NoError(t, core.Start())
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())

	state.UpdateRoundStep(0, RoundStepPrevoteWait)
	require.NoError(t, core.Stop())
	time.Sleep(2 * time.Second)
	require.NoError(t, core.Start())
	assert.Equal(t, RoundStepPrevoteWait, core.CurrentState().Step())

	//wait for 4 second to core's state jump from RoundStepPrevoteWait to RoundStepPrecommit
	time.Sleep(tendermint.DefaultConfig.TimeoutPrevote + (500 * time.Millisecond))
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
	be, _ := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig)
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

func TestCoreFutureMessage(t *testing.T) {
	var (
		nodePrivateKey = tests_utils.MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
		err           error
	)
	//create New test backend and newMockChain
	be, _ := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig)
	require.NoError(t, core.Start())

	logger := zap.NewExample().Sugar()
	zap.ReplaceGlobals(logger.Desugar())

	// create fake block
	tx := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), 800000, big.NewInt(params.GasPriceConfig), nil)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	header := types.CopyHeader(genesisHeader)
	header.Number = big.NewInt(2)
	block := types.NewBlock(header, []*types.Transaction{tx}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block.Transactions(), 1)
	assert.Equal(t, tx.Hash(), block.Transactions()[0].Hash())
	// create fake proposal
	proposal := Proposal{
		Block: block,
		Round: 0,
	}
	msgData, err := rlp.EncodeToBytes(&proposal)
	require.NoError(t, err)
	msg := message{
		Code:    msgPropose,
		Msg:     msgData,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
	}
	sign(t, &msg, nodePrivateKey)
	require.NoError(t, core.handleMsg(msg))
	// create a fake prevote
	hash := block.Hash()
	vote := Vote{
		BlockNumber: big.NewInt(2),
		Round:       0,
		BlockHash:   &hash,
		Seal:        []byte("abc"),
	}
	msgData, err = rlp.EncodeToBytes(&vote)
	require.NoError(t, err)
	prevoteMsg := message{
		Code:    msgPrevote,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
		Msg:     msgData,
	}
	sign(t, &prevoteMsg, nodePrivateKey)
	require.NoError(t, core.handleMsg(prevoteMsg))

	_, err = core.processFutureMessages(logger)
	require.NoError(t, err)
}

func sign(t *testing.T, msg *message, privateKey *ecdsa.PrivateKey) {
	rawPayLoad, err := msg.PayLoadWithoutSignature()
	require.NoError(t, err)
	hashData := crypto.Keccak256(rawPayLoad)
	msg.Signature, err = crypto.Sign(hashData, privateKey)
	require.NoError(t, err)
}
