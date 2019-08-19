
package backend

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode istanbul message")
)

func (sb *backend) decode(msg p2p.Msg) ([]byte, common.Hash, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, common.Hash{}, errDecodeFailed
	}

	hash := common.Hash{}//TODO get from tendermint.RLPHash(data)
	return data, hash, nil
}

// HandleMsg implements consensus.Handler.HandleMsg
// return false if cannot handle message else return false if the message is not handled
func (sb *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if msg.Code == tendermintMsg {
		if !sb.coreStarted {
			return true, tendermint.ErrStoppedEngine
		}

		data, hash, err := sb.decode(msg)
		if err != nil {
			return true, errDecodeFailed
		}

		log.Debug("got the message's hash from peers", "hash", hash.String())
		//TODO: mark peer's message and self known message with the hash get from message
		
		go func() {
			if err := sb.tendermintEventMux.Post(data); err != nil {
				fmt.Printf("error in Post event %v", err)
			}
		}()
		
		return true, nil
	}
	if msg.Code == NewBlockMsg && sb.core.IsProposer() { // eth.NewBlockMsg: import cycle
		// this case is to safeguard the race of similar block which gets propagated from other node while this node is proposing
		// as p2p.Msg can only be decoded once (get EOF for any subsequence read), we need to make sure the payload is restored after we decode it
		//TODO: implement with code is new blockMsg and check is Proposer
	}
	return false, nil
}
