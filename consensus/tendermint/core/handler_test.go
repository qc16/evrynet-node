package core_test

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/backend"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	evrynetCore "github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
)

func TestVerifyProposal(t *testing.T) {
	var (
		nodePrivateKey     = tests.MakeNodeKey()
		nodeFakePrivateKey = tests.MakeNodeKey()
		nodeAddr           = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators         = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
		err           error
	)

	//create New test backend and newMockChain
	be, txPool := mustCreateAndStartNewBackend(t, nodePrivateKey, genesisHeader)

	core := tendermintCore.New(be, tendermint.DefaultConfig, txPool).(tests.TestEngine)
	require.NoError(t, core.Start())
	// --------CASE 1--------
	// Will get errMismatchTxhashes
	block1 := tests.MakeBlockWithoutSeal(genesisHeader)

	// --------CASE 2--------
	// Pass all validation
	tx1 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), 800000, big.NewInt(params.GasPriceConfig), nil)
	tx1, err = types.SignTx(tx1, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block2 := types.NewBlock(genesisHeader, []*types.Transaction{tx1}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block2.Transactions(), 1)
	assert.Equal(t, tx1.Hash(), block2.Transactions()[0].Hash())

	// --------CASE 3--------
	// Will get ErrInsufficientFunds
	tx2 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), params.GasPriceConfig, big.NewInt(params.GasPriceConfig), nil)
	tx2, err = types.SignTx(tx2, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block3 := types.NewBlock(genesisHeader, []*types.Transaction{tx2}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block3.Transactions(), 1)
	assert.Equal(t, tx2.Hash(), block3.Transactions()[0].Hash())

	// --------CASE 4--------
	// Node propose fake block1 (fake signature)
	// backend.VerifyHeader() will return error
	tx3 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), 800000, big.NewInt(params.GasPriceConfig), nil)
	tx3, err = types.SignTx(tx3, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block4 := types.NewBlock(genesisHeader, []*types.Transaction{tx3}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block4.Transactions(), 1)
	assert.Equal(t, tx3.Hash(), block4.Transactions()[0].Hash())

	for _, testCase := range []struct {
		block      *types.Block
		privateKey *ecdsa.PrivateKey
		assertFn   func(err error)
	}{
		{
			block:      block1,
			privateKey: nodePrivateKey,
			assertFn: func(err error) {
				assert.EqualError(t, err, tendermint.ErrMismatchTxhashes.Error())
			},
		},
		{
			block:      block2,
			privateKey: nodePrivateKey,
			assertFn: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			block:      block3,
			privateKey: nodePrivateKey,
			assertFn: func(err error) {
				assert.EqualError(t, err, evrynetCore.ErrInsufficientFunds.Error())
			},
		},
		{
			block:      block4,
			privateKey: nodeFakePrivateKey,
			assertFn: func(err error) {
				assert.EqualError(t, err, tendermintCore.ErrInvalidProposalSignature.Error())
			},
		},
	} {
		proposal := tendermint.Proposal{
			Block: testCase.block,
			Round: 1,
		}

		msgData, err := rlp.EncodeToBytes(&proposal)
		require.NoError(t, err)

		// Create fake message from another node address
		msg := tendermintCore.Message{
			Code:    0,
			Msg:     msgData,
			Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
		}

		msgPayLoadWithoutSignature, err := rlp.EncodeToBytes(&tendermintCore.Message{
			Code:          msg.Code,
			Address:       msg.Address,
			Msg:           msg.Msg,
			Signature:     []byte{},
			CommittedSeal: msg.CommittedSeal,
		})
		require.NoError(t, err)

		signature, err := crypto.Sign(crypto.Keccak256(msgPayLoadWithoutSignature), testCase.privateKey)
		require.NoError(t, err)
		msg.Signature = signature
		testCase.assertFn(core.VerifyProposal(proposal, msg))
	}
}

func mustCreateAndStartNewBackend(t *testing.T, nodePrivateKey *ecdsa.PrivateKey, genesisHeader *types.Header) (tbe tests.TestBackend, txPool *evrynetCore.TxPool) {
	var (
		address = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		trigger = false
		statedb = tests.MustCreateStateDB(t)

		testTxPoolConfig evrynetCore.TxPoolConfig
		blockchain       = &tests.TestChain{
			GenesisHeader: genesisHeader,
			TestBlockChain: &tests.TestBlockChain{
				Statedb:       statedb,
				GasLimit:      1000000000,
				ChainHeadFeed: new(event.Feed),
			},
			Address: address,
			Trigger: &trigger,
		}
		pool   = evrynetCore.NewTxPool(testTxPoolConfig, params.TendermintTestChainConfig, blockchain)
		memDB  = ethdb.NewMemDatabase()
		config = tendermint.DefaultConfig
		be     = backend.New(config, nodePrivateKey, pool, backend.WithDB(memDB)).(tests.TestBackend)
	)
	statedb.SetBalance(address, new(big.Int).SetUint64(params.Ether))
	tests.MustStartTestChainAndBackend(be, blockchain)
	return be, pool
}
