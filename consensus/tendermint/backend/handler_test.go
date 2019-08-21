package backend

import (
	"testing"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
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
