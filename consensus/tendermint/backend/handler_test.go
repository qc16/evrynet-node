
package backend

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestTendermintMessage(t *testing.T) {
	privateKey, _ := generatePrivateKey()
	b := New(privateKey).(*backend)
	b.Start()
	// generate one msg
	data := []byte("data1")
	msg := makeMsg(tendermintMsg, data)
	addr := getAddress()

	// 2. this message should be in cache after we handle it
	_, err := b.HandleMsg(addr, msg)
	if err != nil {
		t.Fatalf("handle message failed: %v", err)
	}
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}
