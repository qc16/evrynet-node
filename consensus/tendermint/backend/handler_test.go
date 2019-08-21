package backend

import (
	"testing"

	"github.com/evrynet-official/evrynet-client/consensus"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/vm"
	"github.com/evrynet-official/evrynet-client/p2p"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
)

func TestHandleMsg(t *testing.T) {
	b := newTestBackend()
	// generate one msg
	data := []byte("data1")
	msg := makeMsg(tendermintMsg, data)
	addr := getAddress()

	// 2. this message should be in cache after we handle it
	handled, err := b.HandleMsg(addr, msg)
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

func newTestBackend() *backend {
	var (
		engine consensus.Engine = newEngine()
		db                      = rawdb.NewMemoryDatabase()
	)

	blockchain, _ := core.NewBlockChain(db, nil, params.TendermintTestChainConfig, engine, vm.Config{}, nil)

	b := engine.(*backend)
	b.Start(blockchain, nil)
	return b
}
