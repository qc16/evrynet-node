package tests

import (
	"crypto/ecdsa"

	"github.com/evrynet-official/evrynet-client/crypto"
)

func mustGeneratePrivateKey(key string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		panic(err)
	}
	return privateKey
}
