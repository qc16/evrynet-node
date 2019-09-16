package backend

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/validator"
	"github.com/evrynet-official/evrynet-client/crypto"
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

	valSet1 := backend.Validators(big.NewInt(1))
	if valSet1.Size() != 1 {
		t.Errorf("Valset size of block 1st should be 1, get: %d", valSet1.Size())
	}

	valSet2 := backend.Validators(big.NewInt(2))
	if valSet2.Size() != 0 {
		t.Errorf("Valset size of block 2th should be 0, get: %d", valSet2.Size())
	}
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
	vset := validator.NewSet(addrs, tendermint.RoundRobin)
	return vset
}
