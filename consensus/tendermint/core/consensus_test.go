package core

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests_utils"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/utils"
	"github.com/evrynet-official/evrynet-client/core/types"
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
	core.currentState.commitRound = voteRound

	view := tendermint.View{
		BlockNumber: voteBlockNumber,
		Round:       voteRound,
	}

	testCases := []struct {
		name           string
		validatorVotes map[int]int
		totalReceived  int
		assertFn       func(block *types.Block, err error)
	}{
		{
			name: "Case 1",
			validatorVotes: map[int]int{
				0: 2,
				1: 2,
				2: 2,
				3: 2,
			},
			totalReceived: 4,
			assertFn: func(block *types.Block, err error) {
				assert.NotNil(t, block)
				assert.NoError(t, err)
			},
		},
		{
			name: "Case 2: Validator 0,2,3 vote for block 2. Validator 2 votes for block 1",
			validatorVotes: map[int]int{
				0: 2,
				1: 1,
				2: 2,
				3: 2,
			},
			totalReceived: 3,
			assertFn: func(block *types.Block, err error) {
				assert.NotNil(t, block)
				assert.NoError(t, err)
			},
		},
		{
			name: "Case 3: Validator 0,1 will vote for block 2. Validator 2,3 will vote for block 1",
			validatorVotes: map[int]int{
				0: 2,
				1: 2,
				2: 1,
				3: 1,
			},
			totalReceived: 2,
			assertFn: func(block *types.Block, err error) {
				assert.Nil(t, block)
				assert.Error(t, err) // Get error "not enough precommits received expect at least 3 received 2"
			},
		},
	}

	for _, tc := range testCases {
		validateVote := func(t *testing.T) {
			newMsgSet := newMessageSet(core.valSet, msgPrecommit, &view)

			//Create block 1
			genesisHeader.Number = big.NewInt(1)
			bl1 := tests_utils.MakeBlockWithoutSeal(genesisHeader)
			blHash1 := bl1.Hash()
			committedSeal1, err := core.backend.Sign(utils.PrepareCommittedSeal(blHash1))
			require.NoError(t, err)

			//Create block 2
			genesisHeader.Number = big.NewInt(2)
			bl2 := tests_utils.MakeBlockWithoutSeal(genesisHeader)
			blHash2 := bl2.Hash()
			committedSeal2, err := core.backend.Sign(utils.PrepareCommittedSeal(blHash2))
			require.NoError(t, err)
			require.NotEqual(t, bl1.Hash().Hex(), bl2.Hash().Hex(), "Block hash of 2 blocks must be different")

			//Add vote from node 1,2,3,4
			for index, valAddr := range validators {
				msg := message{
					Code:    msgPrecommit,
					Address: valAddr,
				}
				switch tc.validatorVotes[index] {
				case 1:
					ok, err := newMsgSet.AddVote(msg,
						&tendermint.Vote{
							BlockHash:   &blHash1,
							BlockNumber: core.CurrentState().BlockNumber(),
							Round:       voteRound,
							Seal:        committedSeal1,
						})
					require.NoError(t, err)
					assert.True(t, ok)
				case 2:
					ok, err := newMsgSet.AddVote(msg,
						&tendermint.Vote{
							BlockHash:   &blHash2,
							BlockNumber: core.CurrentState().BlockNumber(),
							Round:       voteRound,
							Seal:        committedSeal2,
						})
					require.NoError(t, err)
					assert.True(t, ok)
				default:
					fmt.Println("Not support this case")
				}
			}

			assert.Equal(t, 4, newMsgSet.totalReceived)
			core.currentState.PrecommitsReceived[voteRound] = newMsgSet
			assert.Equal(t, tc.totalReceived, core.currentState.PrecommitsReceived[voteRound].voteByBlock[blHash2].totalReceived, "Total Precommits Received on block 2 must be same when getting vote by block hash")

			tc.assertFn(core.FinalizeBlock(&tendermint.Proposal{
				Block:    bl2,
				Round:    0,
				POLRound: 0,
			}))
		}

		t.Run(tc.name, validateVote)
	}
}
