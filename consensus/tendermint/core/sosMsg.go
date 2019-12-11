package core

import (
	"encoding/binary"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/evrynet-official/evrynet-client/ethdb"
	"github.com/evrynet-official/evrynet-client/rlp"
	"go.uber.org/zap"
)

const (
	dbKeySOSMsgPrefix = "tendermint-rsh-"
)

// SOSMsg is the struct of SOS message
type SOSMsg struct {
	logger *zap.SugaredLogger
	db     ethdb.Database
	mu     sync.Mutex
}

// SosData contain data for message stored
type SosData struct {
	Time        time.Time     `json:"Time"`
	BlockNumber uint64        `json:"BlockNumber"`
	Step        RoundStepType `json:"Step"`
	Round       int64         `json:"Round"`
	Data        []byte        `json:"Data"`
}

// NewSOSMsg return new instance of SOSMsg
func NewSOSMsg(db ethdb.Database) *SOSMsg {
	return &SOSMsg{
		logger: zap.S(),
		db:     db,
	}
}

// StoreSentMsg stores vote/ propose to database
func (c *core) StoreSentMsg(step RoundStepType, round int64, msg interface{}) {
	var (
		blockNumber = c.currentState.BlockNumber().Uint64()
		logger      = c.sosMsg.logger.With("block", blockNumber, "step", step.String(), "round", round)
	)

	if index := c.LookupSentMsg(step, round); index >= 0 {
		logger.Warnw("message had saved at before")
		return
	}
	c.sosMsg.mu.Lock()
	defer c.sosMsg.mu.Unlock()

	msgBytes, err := rlp.EncodeToBytes(msg)
	if err != nil {
		logger.Warnw("failed to encode rlp for sentMsg data", "err", err)
		return
	}

	var (
		key   = getKey(blockNumber, step, round)
		sData = &SosData{
			Time:        time.Now(),
			BlockNumber: blockNumber,
			Round:       round,
			Step:        step,
			Data:        msgBytes,
		}
	)

	blob, err := json.Marshal(sData)
	if err != nil {
		logger.Warnw("failed to encode rlp for SOSMsg data", "err", err)
		return
	}

	logger.Infow("saving a sent Msg")
	if err = c.sosMsg.db.Put(key, blob); err != nil {
		logger.Warnw("failed write to SOSMsg file", "err", err)
	}
}

// LookupSentMsg lockups proposal/ vote messages had stored and return index of the message
// if message were not found returns -1
func (c *core) LookupSentMsg(step RoundStepType, round int64) int64 {
	c.sosMsg.mu.Lock()
	defer c.sosMsg.mu.Unlock()
	var (
		msgs        []*SosData
		prefix      []byte
		blockNumber = c.currentState.BlockNumber().Uint64()
		logger      = c.sosMsg.logger.With("block", blockNumber, "step", step.String(), "round", round)
	)

	prefix = append(prefix, []byte(dbKeySOSMsgPrefix)...)
	prefix = append(prefix, encodeUint64(blockNumber)...)
	it := c.sosMsg.db.NewIteratorWithPrefix(prefix)
	for it.Next() {
		blob, err := c.sosMsg.db.Get(it.Key())
		if err != nil {
			continue
		}
		var sData *SosData
		if err := json.Unmarshal(blob, &sData); err != nil {
			continue
		}
		msgs = append(msgs, sData)
	}

	if len(msgs) > 0 {
		// sort array
		sort.Slice(msgs, func(i, j int) bool {
			return msgs[i].Time.Before(msgs[j].Time)
		})
		// find index of message
		for i := 0; i < len(msgs); i++ {
			if msgs[i].Step == step && msgs[i].Round == round {
				return int64(i)
			}
		}
	}
	logger.Warnw("LookupSentMsg: message not found")
	return int64(-1)
}

// TruncateMsgStored removes all data stored by the block's number
func (c *core) TruncateMsgStored() {
	c.sosMsg.mu.Lock()
	defer c.sosMsg.mu.Unlock()

	var (
		key         []byte
		blockNumber = c.currentState.BlockNumber().Uint64()
		logger      = c.sosMsg.logger.With("block", blockNumber)
	)

	key = append(key, []byte(dbKeySOSMsgPrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	it := c.sosMsg.db.NewIteratorWithPrefix(key)
	for it.Next() {
		if err := c.sosMsg.db.Delete(it.Key()); err != nil {
			logger.Warnw("failed to delete SOSMsg by key", "err", err)
		}
	}
	logger.Info("truncate SOSMsg done", "block", c.currentState.BlockNumber().Uint64())
}

// UnmarshalJSON unmarshals from JSON.
func (h *SosData) UnmarshalJSON(input []byte) error {
	type sosData SosData
	var dec sosData
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	h.BlockNumber = dec.BlockNumber
	h.Step = dec.Step
	h.Round = dec.Round
	h.Data = dec.Data
	return nil
}

func encodeUint64(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

func encodeInt64(number int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, number)
	return buf[:n]
}

func getKey(blockNumber uint64, step RoundStepType, round int64) []byte {
	var key []byte
	key = append(key, []byte(dbKeySOSMsgPrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	key = append(key, []byte(step.String())...)
	key = append(key, encodeInt64(round)...)
	return key
}
