package core

import (
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/core/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/crypto"
)

func TestFinalizeBlock(t *testing.T) {
	var (
		nodePrivateKey = tests_utils.MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
			common.HexToAddress("0x954e4BF2C68F13D97C45db0e02645D145dB6911f"),
			common.HexToAddress("0x45F8B547A7f16730c0C8961A21b56c31d84DdB49"),
			common.HexToAddress("0x5be60024b3b7EF2f6e4db97641e8942b85a5124e"),
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	//create New test backend and newMockChain
	be, txPool := tests_utils.MustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader, validators)

	core := newTestCore(be, tendermint.DefaultConfig, txPool)
	require.NoError(t, core.Start())
	state := core.CurrentState()
	assert.Equal(t, RoundStepNewHeight, state.Step())

	//fake core.CurrentState
	var (
		voteRound       = int64(0)
		voteBlockNumber = big.NewInt(1)
	)
	core.currentState = core.getInitializedState()

	view := tendermint.View{
		BlockNumber: voteBlockNumber,
		Round:       voteRound,
	}

	blockHasSeal := tests_utils.MakeBlockWithSeal(be, genesisHeader)
	blockHashHasSeal := blockHasSeal.Hash()
	seal, err := core.backend.Sign(utils.PrepareCommittedSeal(blockHasSeal.Hash()))
	require.NoError(t, err)

	for _, testCase := range []struct {
		addWrongBlockHash map[int]bool
		totalReceived     int
		assertFn          func(block *types.Block, err error)
	}{
		{
			addWrongBlockHash: map[int]bool{},
			totalReceived:     4,
			assertFn: func(block *types.Block, err error) {
				assert.NotNil(t, block)
				assert.NoError(t, err)
			},
		},
		{
			addWrongBlockHash: map[int]bool{
				3: true,
			},
			totalReceived: 3,
			assertFn: func(block *types.Block, err error) {
				assert.NotNil(t, block)
				assert.NoError(t, err)
			},
		},
		{
			addWrongBlockHash: map[int]bool{
				2: true,
				3: true,
			},
			totalReceived: 2,
			assertFn: func(block *types.Block, err error) {
				assert.Nil(t, block)
				assert.Error(t, err) // Get error "not enough precommits received expect at least 3 received 2"
			},
		},
	} {
		newMsgSet := newMessageSet(core.valSet, msgPrecommit, &view)
		newBlock := tests_utils.MakeBlockWithoutSeal(genesisHeader)
		newBlockHash := newBlock.Hash()
		newSeal, err := core.backend.Sign(utils.PrepareCommittedSeal(newBlockHash))
		require.NoError(t, err)

		//Add vote from node 1,2,3,4
		for index, valAddr := range validators {
			if testCase.addWrongBlockHash[index] {
				ok, err := newMsgSet.AddVote(
					message{
						Code:    msgPrecommit,
						Address: valAddr,
					},
					&tendermint.Vote{
						BlockHash:   &blockHashHasSeal,
						BlockNumber: core.CurrentState().BlockNumber(),
						Round:       voteRound,
						Seal:        seal,
					})
				require.NoError(t, err)
				assert.True(t, ok)
			} else {
				ok, err := newMsgSet.AddVote(
					message{
						Code:    msgPrecommit,
						Address: valAddr,
					},
					&tendermint.Vote{
						BlockHash:   &newBlockHash,
						BlockNumber: core.CurrentState().BlockNumber(),
						Round:       voteRound,
						Seal:        newSeal,
					})
				require.NoError(t, err)
				assert.True(t, ok)
			}
		}

		core.currentState.PrecommitsReceived[voteRound] = newMsgSet
		assert.Equal(t, testCase.totalReceived, core.currentState.PrecommitsReceived[voteRound].voteByBlock[newBlockHash].totalReceived)

		testCase.assertFn(core.FinalizeBlock(&tendermint.Proposal{
			Block:    newBlock,
			Round:    1,
			POLRound: 0,
		}))
	}
}
