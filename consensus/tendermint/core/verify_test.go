package core

import (
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	evrynetCore "github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	var (
		nodePrivateKey     = tests.MakeNodeKey()
		nodeFakePrivateKey = tests.MakeNodeKey()
		nodeAddr           = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators         = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)

	//create New test backend and newMockChain
	chain, be := tests.MustStartTestChainAndBackend(nodePrivateKey, genesisHeader)
	assert.NotNil(t, chain)
	assert.NotNil(t, be)

	// --------CASE 1--------
	// Will get errMismatchTxhashes
	block := tests.MakeBlockWithoutSeal(genesisHeader)
	proposal := tendermint.Proposal{
		Block:    block,
		Round:    0,
		POLRound: 0,
	}
	// Should get error if transactions in block is 0
	assert.EqualError(t, be.Core().Verify(proposal), errMismatchTxhashes.Error())

	// --------CASE 2--------
	// Pass all validation
	tx1 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), 800000, big.NewInt(params.GasPriceConfig), nil)
	tx1, err := types.SignTx(tx1, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block2 := types.NewBlock(genesisHeader, []*types.Transaction{tx1}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block2.Transactions(), 1)
	assert.Equal(t, tx1.Hash(), block2.Transactions()[0].Hash())
	proposal = tendermint.Proposal{
		Block: block2,
	}
	err = be.Core().Verify(proposal)
	// Should get no error if block has transactions
	assert.NoError(t, be.Core().Verify(proposal))

	// --------CASE 3--------
	// Will get ErrInsufficientFunds
	tx2 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), params.GasPriceConfig, big.NewInt(params.GasPriceConfig), nil)
	tx2, err = types.SignTx(tx2, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block3 := types.NewBlock(genesisHeader, []*types.Transaction{tx2}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block3.Transactions(), 1)
	assert.Equal(t, tx2.Hash(), block3.Transactions()[0].Hash())
	proposal = tendermint.Proposal{
		Block: block3,
	}
	// Should get error ErrInsufficientFunds
	assert.EqualError(t, be.Core().Verify(proposal), evrynetCore.ErrInsufficientFunds.Error())

	// --------CASE 4--------
	// Node propose fake block (fake signature)
	// backend.VerifyHeader() will return error
	tx3 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), 800000, big.NewInt(params.GasPriceConfig), nil)
	tx3, err = types.SignTx(tx3, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block4 := types.NewBlock(genesisHeader, []*types.Transaction{tx3}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block4.Transactions(), 1)
	assert.Equal(t, tx3.Hash(), block4.Transactions()[0].Hash())
	proposal = tendermint.Proposal{
		Block: block4,
		Round: 1,
	}

	msgData, err := rlp.EncodeToBytes(&proposal)
	assert.NoError(t, err)

	// Create fake message from another node address
	msg := message{
		Code:    0,
		Msg:     msgData,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
	}

	msgPayLoadWithoutSignature, _ := rlp.EncodeToBytes(&message{
		Code:          msg.Code,
		Address:       msg.Address,
		Msg:           msg.Msg,
		Signature:     []byte{},
		CommittedSeal: msg.CommittedSeal,
	})

	signature, err := crypto.Sign(crypto.Keccak256(msgPayLoadWithoutSignature), nodeFakePrivateKey)
	assert.NoError(t, err)
	msg.Signature = signature

	err = be.Core().VerifyProposal(proposal, msg)
	// Should get error when node send signed msg by fake private key
	assert.EqualError(t, err, ErrInvalidProposalSignature.Error())
}
