package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	tendermintBackend "github.com/ethereum/go-ethereum/consensus/tendermint/backend"
	"github.com/ethereum/go-ethereum/consensus/tendermint/validator"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

func TestTendermintBroadcaster(t *testing.T) {
	var (
		thisNodePk  = mustGeneratePrivateKey(t)
		otherNodePk = mustGeneratePrivateKey(t)
		tbe         = tendermintBackend.New(tendermint.DefaultConfig, thisNodePk)
		addrs       = []common.Address{
			crypto.PubkeyToAddress(thisNodePk.PublicKey),
			crypto.PubkeyToAddress(otherNodePk.PublicKey),
		}
		validatorSet = validator.NewSet(addrs, tendermint.RoundRobin)
	)
	pm, err := newTestProtocolManagerWithConsensus(tbe)
	assert.NoError(t, err)
	defer pm.Stop()

	assert.NoError(t, tbe.Start(nil, nil))
	// create nodes for test peer
	n1 := enode.MustParseV4("enode://" + hex.EncodeToString(crypto.FromECDSAPub(&thisNodePk.PublicKey)[1:]) + "@33.4.2.1:30303")
	n2 := enode.MustParseV4("enode://" + hex.EncodeToString(crypto.FromECDSAPub(&otherNodePk.PublicKey)[1:]) + "@33.4.2.1:30304")

	newTestPeerFromNode(fmt.Sprintf("peer %d", 0), eth63, pm, true, n1)
	newTestPeerFromNode(fmt.Sprintf("peer %d", 1), eth63, pm, true, n2)

	//We don't close the peers since peer.send is asynchronous. When test terminated, the peers will be terminated as well.

	bc, ok := tbe.(tendermint.Backend)
	assert.Equal(t, true, ok)
	payload := []byte("vote message")

	assert.NoError(t, bc.Broadcast(validatorSet, payload))

	//DOTO: set up a test that can receive real message. This might require 2 protocolManagers.
}

func mustGeneratePrivateKey(t *testing.T) *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fail()
	}
	return privateKey
}
