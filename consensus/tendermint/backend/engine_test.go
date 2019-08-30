package backend

import (
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/rlp"
	"github.com/stretchr/testify/assert"
)

// TestSimulateSubscribeAndReceiveToSeal is a simple test to pass a block to backend.Seal()
// on core.handleEvents(), the block is received.
func TestSimulateSubscribeAndReceiveToSeal(t *testing.T) {
	var block types.Block
	// Create a backend and use it as Engine
	engine := newEngine()
	assert.NoError(t, engine.Start(nil, nil))

	blockEnc := common.FromHex("f90260f901f9a083cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f13b19fa1c8dd55a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a05fe50b260da6308036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67a0bc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa7edae9283e52b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000001832fefd8825208845506eb0780a0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")

	assert.NoError(t, rlp.DecodeBytes(blockEnc, &block))

	err := engine.Seal(nil, &block, nil, nil)

	// Sleep to make sure that the block can be received from
	time.Sleep(2 * time.Second)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

}

func newEngine() *backend {
	nodeKey, _ := crypto.GenerateKey()
	be, _ := New(tendermint.DefaultConfig, nodeKey).(*backend)
	be.address = crypto.PubkeyToAddress(nodeKey.PublicKey)
	be.db = ethdb.NewMemDatabase()
	return be
}
