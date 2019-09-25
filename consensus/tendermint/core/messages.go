package core

import (
	"io"
	"sync"

	"github.com/pkg/errors"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/rlp"
)

var (
	ErrConflictingVotes = errors.New("vote received from the same validator for different block in the same round")
	ErrDifferentMsgType = errors.New("message set is not of the same type of the received message")
)

//Engine abstract the core's functionalities
//Note that backend and other packages doesn't care about core's internal logic.
//It only requires core to start receiving/handling messages
//The sending of events/message from core to backend will be done by calling accessing Backend.EventMux()
type Engine interface {
	Start() error
	Stop() error
	//SetBlockForProposal define a method to allow Injecting a Block for testing purpose
	SetBlockForProposal(block *types.Block)
	EventMux() *event.TypeMux
}

// TODO: More msg codes here if needed
const (
	msgPropose uint64 = iota
	msgPrevote
	msgPrecommit
)

type message struct {
	Code      uint64
	Msg       []byte
	Address   common.Address
	Signature []byte
	//TODO: Is CommitedSeal needed in message of Tendermint?
	CommittedSeal []byte
}

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint64
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

func (m *message) PayLoadWithoutSignature() ([]byte, error) {
	return rlp.EncodeToBytes(&message{
		Code:          m.Code,
		Address:       m.Address,
		Msg:           m.Msg,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

// GetAddressFromSignature gets the signer address from the signature
func (m *message) GetAddressFromSignature() (common.Address, error) {
	payLoad, err := m.PayLoadWithoutSignature()
	if err != nil {
		return common.Address{}, err
	}
	hashData := crypto.Keccak256(payLoad)

	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, m.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

//blockVotes store the voting received for a particular block
type blockVotes struct {
	votes         []*tendermint.Vote // validatorIndex -> *Vote
	totalReceived int
}

type messageSet struct {
	view          *tendermint.View
	valSet        tendermint.ValidatorSet
	msgCode       uint64
	messagesMu    *sync.Mutex
	messages      map[common.Address]*message
	voteByAddress map[common.Address]*tendermint.Vote
	voteByBlock   map[common.Hash]*blockVotes
	maj23         *common.Hash
	totalReceived int
	//TODO: Do we have to keep track of which peer has 2/3Majority?
}

// Construct a new message set to accumulate messages for given height/view number.
func newMessageSet(valSet tendermint.ValidatorSet, code uint64, view *tendermint.View) *messageSet {
	return &messageSet{
		view:          view,
		msgCode:       code,
		messagesMu:    new(sync.Mutex),
		messages:      make(map[common.Address]*message),
		voteByBlock:   make(map[common.Hash]*blockVotes),
		voteByAddress: make(map[common.Address]*tendermint.Vote),
		valSet:        valSet,
	}
}

func (ms *messageSet) VotesByAddress() map[common.Address]*tendermint.Vote {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	var (
		ret = make(map[common.Address]*tendermint.Vote)
	)
	for addr, vote := range ms.voteByAddress {
		ret[addr] = vote
	}
	return ret
}

func (ms *messageSet) AddVote(msg message, vote *tendermint.Vote) (bool, error) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	copyHash := common.HexToHash(vote.BlockHash.Hex())
	if ms.msgCode != msg.Code {
		return false, ErrDifferentMsgType
	}
	index, _ := ms.valSet.GetByAddress(msg.Address)
	if index == -1 {
		return false, errors.Wrapf(ErrVoteInvalidValidatorAddress, "address in vote message:%s ", msg.Address.String())
	}
	if ms.view.BlockNumber.Cmp(vote.BlockNumber) != 0 {
		return false, ErrVoteHeightMismatch
	}
	if ms.view.Round != vote.Round {
		log.Error("message set round is not the same as vote round", "msg_set_round", ms.view.Round, "vote_round", vote.Round)
		return false, errors.New("invalid vote for the message set")
	}
	//Signer is supposed to be checked at previous steps so it doesn't need to be check again.

	// if this message set already got this msg, check if the vote is duplicate or double voting
	current, existed := ms.messages[msg.Address]
	if existed {
		var currentVote tendermint.Vote
		if err := rlp.DecodeBytes(current.Msg, &currentVote); err != nil {
			return false, err
		}
		if currentVote.BlockHash.Hex() != vote.BlockHash.Hex() {
			return false, ErrConflictingVotes
		}
		//log.Info("already got vote, skipping", "from", msg.Address, "round", vote.Round)
		return false, nil
	}

	ms.messages[msg.Address] = &msg
	ms.voteByAddress[msg.Address] = vote
	ms.totalReceived++
	ms.addVoteToBlockVote(vote, index)

	if ms.voteByBlock[copyHash].totalReceived > 2*ms.valSet.F() {
		if ms.maj23 == nil {
			ms.maj23 = &copyHash
		}
	}

	return true, nil
}

func (ms *messageSet) addVoteToBlockVote(vote *tendermint.Vote, index int) error {
	bvotes, exist := ms.voteByBlock[*(vote.BlockHash)]
	if !exist {
		bvotes = &blockVotes{
			votes:         make([]*tendermint.Vote, ms.valSet.Size()),
			totalReceived: 0,
		}
	}
	//shouldn't happen but just making sure
	if bvotes.votes[index] != nil && bvotes.votes[index].BlockHash.Hex() != vote.BlockHash.Hex() {
		return ErrConflictingVotes
	}
	bvotes.votes[index] = vote
	bvotes.totalReceived++
	ms.voteByBlock[*(vote.BlockHash)] = bvotes
	return nil
}

func (ms *messageSet) HasMajority() bool {
	if ms == nil {
		return false
	}
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return ms.maj23 != nil
}

func (ms *messageSet) HasTwoThirdAny() bool {
	if ms == nil {
		return false
	}
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return ms.totalReceived > ms.valSet.F()*2
}

//TwoThirdMajority return a blockHash and a bool inidicate if this messageSet hash got a
//TwoThirdMajority on a block
func (ms *messageSet) TwoThirdMajority() (common.Hash, bool) {
	if ms == nil {
		return common.Hash{}, false
	}
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	if ms.maj23 != nil {
		return common.HexToHash(ms.maj23.Hex()), true
	}
	return common.Hash{}, false
}
