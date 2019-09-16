package backend

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	privateKey, _ := generatePrivateKey()
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

	if signer != getAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
	}
}

func TestValidators(t *testing.T) {
	backend, _, blockchain, err := createBlockchainAndBackendFromGenesis()
	assert.NoError(t, err)

	backend.Start(blockchain, nil)

	valSet0 := backend.Validators(big.NewInt(0))
	if valSet0.Size() != 1 {
		t.Errorf("Valset size of zero block should be 1, get: %d", valSet0.Size())
	}
	list := valSet0.List()
	log.Println("validator set of block 0 is")

	for _, val := range list {
		log.Println(val.String())
	}
	valSet1 := backend.Validators(big.NewInt(1))
	if valSet1.Size() != 1 {
		t.Errorf("Valset size of block 1st should be 1, get: %d", valSet1.Size())
	}
	list = valSet1.List()
	log.Println("validator set of block 1 is")

	for _, val := range list {
		log.Println(val.String())
	}
	valSet2 := backend.Validators(big.NewInt(2))
	if valSet2.Size() != 0 {
		t.Errorf("Valset size of block 2th should be 0, get: %d", valSet2.Size())
	}
}

func TestVerify(t *testing.T) {
	var (
		nodePrivateKey = makeNodeKey()
		nodeAddr       = crypto.PubkeyToAddress(nodePrivateKey.PublicKey)
		validators     = []common.Address{
			nodeAddr,
		}
		genesisHeader = makeGenesisHeader(validators)
	)

	//create New test backend and newMockChain
	chain, engine := mustStartTestChainAndBackend(nodePrivateKey, genesisHeader)
	assert.NotNil(t, chain)
	assert.NotNil(t, engine)
	assert.Equal(t, true, engine.coreStarted)

	// --------CASE 1--------
	// Will get errMismatchTxhashes
	block := makeBlockWithoutSeal(genesisHeader)
	proposal := tendermint.Proposal{
		Block:    block,
		Round:    0,
		POLRound: 0,
	}
	// Should get error if transactions in block is 0
	assert.Error(t, engine.Verify(proposal), errMismatchTxhashes)

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
	err = engine.Verify(proposal)
	// Should get no error if block has transactions
	assert.NoError(t, engine.Verify(proposal))

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
	assert.Error(t, engine.Verify(proposal), core.ErrInsufficientFunds)

	// --------CASE 4--------
	// Header Difficulty is nil
	// backend.VerifyHeader() will return error
	tx3 := types.NewTransaction(0, common.HexToAddress("A8A620a156121f6Ef0Bb0bF0FFe1B6A0e02834a1"), big.NewInt(10), params.GasPriceConfig, big.NewInt(params.GasPriceConfig), nil)
	tx3, err = types.SignTx(tx3, types.HomesteadSigner{}, nodePrivateKey)
	assert.NoError(t, err)

	editedHeader := *genesisHeader
	editedHeader.Difficulty = nil
	block4 := types.NewBlock(&editedHeader, []*types.Transaction{tx3}, []*types.Header{}, []*types.Receipt{})
	assert.Len(t, block4.Transactions(), 1)
	assert.Equal(t, tx3.Hash(), block4.Transactions()[0].Hash())
	proposal = tendermint.Proposal{
		Block: block4,
	}
	// Should get error when header Header Difficulty is nil
	assert.Error(t, engine.Verify(proposal), errInvalidDifficulty)
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestValidatorSet(n int) tendermint.ValidatorSet {
	// generate validators
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, tendermint.RoundRobin, int64(0))
	return vset
}
