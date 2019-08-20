
package backend

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/consensus/tendermint"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"

	"golang.org/x/crypto/sha3"
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

		log.Debug("handler.HandleMsg implement me mark peer's message amd self known message", "hash", hash.String())
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
		if reader, ok := msg.Payload.(*bytes.Reader); ok {
			payload, err := ioutil.ReadAll(reader)
			if err != nil {
				return true, err
			}
			reader.Reset(payload)       // ready to be decoded
			defer reader.Reset(payload) // restore so main eth/handler can decode
			var request struct {        // this has to be same as eth/protocol.go#newBlockData as we are reading NewBlockMsg
				Block *types.Block
				TD    *big.Int
			}
			if err := msg.Decode(&request); err != nil {
				log.Debug("Proposer was unable to decode the NewBlockMsg", "error", err)
				return false, nil
			}
			newRequestedBlock := request.Block
			if newRequestedBlock.Header().MixDigest == types.TendermintDigest && sb.core.IsCurrentProposal(newRequestedBlock.Hash()) {
				log.Debug("Proposer already proposed this block", "hash", newRequestedBlock.Hash(), "sender", addr)
				return true, nil
			}
		}
	}
	return false, nil
}
