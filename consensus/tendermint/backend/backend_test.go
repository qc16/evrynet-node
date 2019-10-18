package backend

import (
	"log"
	"math/big"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	privateKey, _ := tests.GeneratePrivateKey()
	b := &Backend{
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
