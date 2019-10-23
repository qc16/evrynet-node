package backend

import (
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint/tests"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/p2p"
	"github.com/evrynet-official/evrynet-client/rlp"
	"github.com/stretchr/testify/assert"
)

func TestHandleMsg(t *testing.T) {
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

	// generate one msg
	data := []byte("data1")
	msg := makeMsg(consensus.TendermintMsg, data)
	addr := tests.GetAddress()

	// 2. this message should be in cache after we handle it
	handled, err := be.HandleMsg(addr, msg)
	if err != nil {
		t.Errorf("expected message being handled successfully but got %s", err)
	}
	if !handled {
		t.Errorf("expected message not being handled")
	}
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}
