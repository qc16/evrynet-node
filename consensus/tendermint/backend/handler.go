package backend

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/sha3"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/p2p"
	"github.com/evrynet-official/evrynet-client/rlp"
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode istanbul message")
)

func rLPHash(v interface{}) (h common.Hash) {
	hw := sha3.New256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

func (sb *backend) decode(msg p2p.Msg) ([]byte, common.Hash, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, common.Hash{}, errDecodeFailed
	}
	return data, rLPHash(data), nil
}

// HandleMsg implements consensus.Handler.HandleMsg
// return false if the message cannot be handle by Tendermint Backend
func (sb *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	switch msg.Code {
	case tendermintMsg:
		if !sb.coreStarted {
			return true, tendermint.ErrStoppedEngine
		}

		data, hash, err := sb.decode(msg)
		if err != nil {
			return true, errDecodeFailed
		}
		//log is used at local package level for testing now
		log.Printf("--- Message from peer %s, msgCode is %d msg Hash is %s \n", addr.Hex(), msg.Code, hash.String())
		//TODO: mark peer's message and self known message with the hash get from message

		go func() {
			if err := sb.tendermintEventMux.Post(data); err != nil {
				log.Printf("error in Post event %v", err)
			}
		}()

		return true, nil
	default:
		return false, fmt.Errorf("unknown message code %d for Tendermint's protocol", msg.Code)
		//TODO:Handler other cases
		//Case 1: NewBlock when this node is the propose.
		//More cases to be added...
	}
}
