package backend

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	tendermintCore "github.com/evrynet-official/evrynet-client/consensus/tendermint/core"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	evrynetCode "github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/state"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/eth/transaction"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1")
	b := &backend{
		privateKey: privateKey,
	}
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Check signature recover
	hashData := crypto.Keccak256([]byte(data))
	pubkey, _ := crypto.Ecrecover(hashData, sig)

	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	if signer != tests.GetAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), tests.GetAddress().Hex())
	}
}

func TestValidators(t *testing.T) {
	var (
		nodePrivateKey = tests.MakeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests.MakeGenesisHeader(validators)
	)

	//create New test backend and newMockChain
	be, ok := mustCreateAndStartNewBackend(nodePrivateKey, genesisHeader)
	assert.True(t, ok)
	assert.NotNil(t, be.TxPool())

	valSet0 := be.Validators(big.NewInt(0))
	if valSet0.Size() != 1 {
		t.Errorf("Valset size of zero block should be 1, get: %d", valSet0.Size())
	}
	list := valSet0.List()
	log.Println("validator set of block 0 is")

	for _, val := range list {
		log.Println(val.String())
	}
	valSet1 := be.Validators(big.NewInt(1))
	if valSet1.Size() != 1 {
		t.Errorf("Valset size of block 1st should be 1, get: %d", valSet1.Size())
	}
	list = valSet1.List()
	log.Println("validator set of block 1 is")

	for _, val := range list {
		log.Println(val.String())
	}
	valSet2 := be.Validators(big.NewInt(2))
	if valSet2.Size() != 0 {
		t.Errorf("Valset size of block 2th should be 0, get: %d", valSet2.Size())
	}
}

func TestVerifyProposal(t *testing.T) {
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
	be, ok := mustCreateAndStartNewBackend(nodePrivateKey, genesisHeader)
	assert.True(t, ok)

	core := be.Core().(tests.TestEngine)
	err := core.Start()
	assert.Nil(t, err)

	// --------CASE 1--------
	// Will get errMismatchTxhashes
	block := tests.MakeBlockWithoutSeal(genesisHeader)
	proposal := tendermint.Proposal{
		Block:    block,
		Round:    1,
		POLRound: 0,
	}
	msgData, err := rlp.EncodeToBytes(&proposal)
	assert.NoError(t, err)
	// Create fake message from another node address
	msg := tendermintCore.Message{
		Code:    0,
		Msg:     msgData,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
	}
	msgPayLoadWithoutSignature, _ := rlp.EncodeToBytes(&tendermintCore.Message{
		Code:          msg.Code,
		Address:       msg.Address,
		Msg:           msg.Msg,
		Signature:     []byte{},
		CommittedSeal: msg.CommittedSeal,
	})
	signature, err := crypto.Sign(crypto.Keccak256(msgPayLoadWithoutSignature), nodePrivateKey)
	assert.NoError(t, err)
	msg.Signature = signature
	// Should get error if transactions in block is 0
	assert.EqualError(t, core.VerifyProposal(proposal, msg), tendermint.ErrMismatchTxhashes.Error())

	// --------CASE 2--------
	// Pass all validation
	tx1 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), 800000, big.NewInt(params.GasPriceConfig), nil)
	tx1, err = types.SignTx(tx1, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	block2 := types.NewBlock(genesisHeader, []*types.Transaction{tx1}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block2.Transactions(), 1)
	assert.Equal(t, tx1.Hash(), block2.Transactions()[0].Hash())
	proposal = tendermint.Proposal{
		Block: block2,
		Round: 1,
	}
	msgData, err = rlp.EncodeToBytes(&proposal)
	assert.NoError(t, err)
	// Create fake message from another node address
	msg = tendermintCore.Message{
		Code:    0,
		Msg:     msgData,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
	}
	msgPayLoadWithoutSignature, _ = rlp.EncodeToBytes(&tendermintCore.Message{
		Code:          msg.Code,
		Address:       msg.Address,
		Msg:           msg.Msg,
		Signature:     []byte{},
		CommittedSeal: msg.CommittedSeal,
	})
	signature, err = crypto.Sign(crypto.Keccak256(msgPayLoadWithoutSignature), nodePrivateKey)
	assert.NoError(t, err)
	msg.Signature = signature
	// Should get no error if block has transactions
	assert.NoError(t, core.VerifyProposal(proposal, msg))

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
		Round: 1,
	}
	msgData, err = rlp.EncodeToBytes(&proposal)
	assert.NoError(t, err)
	// Create fake message from another node address
	msg = tendermintCore.Message{
		Code:    0,
		Msg:     msgData,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
	}
	msgPayLoadWithoutSignature, _ = rlp.EncodeToBytes(&tendermintCore.Message{
		Code:          msg.Code,
		Address:       msg.Address,
		Msg:           msg.Msg,
		Signature:     []byte{},
		CommittedSeal: msg.CommittedSeal,
	})
	signature, err = crypto.Sign(crypto.Keccak256(msgPayLoadWithoutSignature), nodePrivateKey)
	assert.NoError(t, err)
	msg.Signature = signature
	// Should get error ErrInsufficientFunds
	assert.EqualError(t, core.VerifyProposal(proposal, msg), evrynetCode.ErrInsufficientFunds.Error())

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

	msgData, err = rlp.EncodeToBytes(&proposal)
	assert.NoError(t, err)

	// Create fake message from another node address
	msg = tendermintCore.Message{
		Code:    0,
		Msg:     msgData,
		Address: crypto.PubkeyToAddress(nodePrivateKey.PublicKey),
	}

	msgPayLoadWithoutSignature, _ = rlp.EncodeToBytes(&tendermintCore.Message{
		Code:          msg.Code,
		Address:       msg.Address,
		Msg:           msg.Msg,
		Signature:     []byte{},
		CommittedSeal: msg.CommittedSeal,
	})

	signature, err = crypto.Sign(crypto.Keccak256(msgPayLoadWithoutSignature), nodeFakePrivateKey)
	assert.NoError(t, err)
	msg.Signature = signature

	err = core.VerifyProposal(proposal, msg)
	// Should get error when node send signed msg by fake private key
	assert.EqualError(t, err, tendermintCore.ErrInvalidProposalSignature.Error())
}

func mustCreateAndStartNewBackend(nodePrivateKey *ecdsa.PrivateKey, genesisHeader *types.Header) (tests.TestBackend, bool) {
	address := crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
	trigger := false
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
	statedb.SetBalance(address, new(big.Int).SetUint64(params.Ether))
	var testTxPoolConfig evrynetCode.TxPoolConfig
	blockchain := &tests.TestChain{
		GenesisHeader: genesisHeader,
		TestBlockChain: &tests.TestBlockChain{
			Statedb:       statedb,
			GasLimit:      1000000000,
			ChainHeadFeed: new(event.Feed),
		},
		Address: address,
		Trigger: &trigger,
	}
	pool := evrynetCode.NewTxPool(testTxPoolConfig, params.TendermintTestChainConfig, blockchain)
	defer pool.Stop()
	memDB := ethdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	be := New(config, nodePrivateKey, WithTxPoolOpts(&transaction.TxPoolOpts{CoreTxPool: pool}), WithDB(memDB)).(tests.TestBackend)
	ok := tests.MustStartTestChainAndBackend(be, blockchain)
	return be, ok
}
