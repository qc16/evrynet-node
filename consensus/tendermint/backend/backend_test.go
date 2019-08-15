package backend

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/consensus/tendermint/validator"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	addressSet = []common.Address{
		common.HexToAddress("0x3Cf628d49Ae46b49b210F0521Fbd9F82B461A9E1"),
		common.HexToAddress("0x723f12209b9C71f17A7b27FCDF16CA5883b7BBB0"),
	}
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

// Test broadcast in consensus
func TestBroadcast(t *testing.T) {
	backend := &backend{}
	payload := []byte("vote message")
	validatorSet := newTestValidatorSet(2)
	err := backend.Broadcast(validatorSet, payload)
	if err != nil {
		t.Fatalf("can't broadcast to validators: %v", err)
	}
}

// Test Gossip between validators in consensus
func TestGossip(t *testing.T) {
	backend := &backend{}
	payload := []byte("vote message")
	validatorSet := newTestValidatorSet(2)
	err := backend.Gossip(validatorSet, payload)
	if err != nil {
		t.Fatalf("can't gossip to validators: %v", err)
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
