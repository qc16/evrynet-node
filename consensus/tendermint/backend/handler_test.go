
package backend

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/params"
)

func TestTendermintMessage(t *testing.T) {
	_, b := newBlockChain()
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

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func newBlockChain() (*core.BlockChain, *backend) {
	var (
		engine consensus.Engine =  newEngine()
		db                      = rawdb.NewMemoryDatabase()
	)

	// Use the first key as private key
	blockchain, _ := core.NewBlockChain(db, nil, params.TendermintTestChainConfig, engine, vm.Config{}, nil)

	b := engine.(*backend)
	b.Start(blockchain, nil)
	return blockchain, b
}