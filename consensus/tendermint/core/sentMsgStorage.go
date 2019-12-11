package core

import (
	"encoding/binary"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/evrynet-official/evrynet-client/ethdb"

	"go.uber.org/zap"
)

const (
	dbKeyMsgStoragePrefix = "tendermint-ms-"
)

// MsgStorage is the struct of SOS message
type MsgStorage struct {
	logger *zap.SugaredLogger
	db     ethdb.Database
	mu     sync.Mutex
}

// MsgStorageData contain data for message stored
type MsgStorageData struct {
	Time        time.Time     `json:"Time"`
	BlockNumber uint64        `json:"BlockNumber"`
	Step        RoundStepType `json:"Step"`
	Round       int64         `json:"Round"`
	Data        []byte        `json:"Data"`
}

// NewMsgStorageData return new instance of MsgStorage
func NewMsgStorageData(db ethdb.Database) *MsgStorage {
	return &MsgStorage{
		logger: zap.S(),
		db:     db,
	}
}

// storeSentMsg stores vote/ propose to database
func (c *core) storeSentMsg(step RoundStepType, round int64, msgData []byte) {
	var (
		blockNumber = c.currentState.BlockNumber().Uint64()
		logger      = c.sentMsgStorage.logger.With("block", blockNumber, "step", step.String(), "round", round)
	)

	if index := c.lookupSentMsg(step, round); index >= 0 {
		logger.Warnw("message had saved at before")
		return
	}
	c.sentMsgStorage.mu.Lock()
	defer c.sentMsgStorage.mu.Unlock()

	var (
		key   = getKey(blockNumber, step, round)
		sData = &MsgStorageData{
			Time:        time.Now(),
			BlockNumber: blockNumber,
			Round:       round,
			Step:        step,
			Data:        msgData,
		}
	)

	blob, err := json.Marshal(sData)
	if err != nil {
		logger.Warnw("failed to encode rlp for MsgStorage data", "err", err)
		return
	}

	logger.Infow("saving a sent Msg")
	if err = c.sentMsgStorage.db.Put(key, blob); err != nil {
		logger.Warnw("failed write to MsgStorage file", "err", err)
	}
}

// lookupSentMsg lockups proposal/ vote messages had stored and return index of the message
// if message were not found returns -1
func (c *core) lookupSentMsg(step RoundStepType, round int64) int64 {
	c.sentMsgStorage.mu.Lock()
	defer c.sentMsgStorage.mu.Unlock()
	var (
		msgs        []*MsgStorageData
		prefix      []byte
		blockNumber = c.currentState.BlockNumber().Uint64()
		logger      = c.sentMsgStorage.logger.With("block", blockNumber, "step", step.String(), "round", round)
	)

	prefix = append(prefix, []byte(dbKeyMsgStoragePrefix)...)
	prefix = append(prefix, encodeUint64(blockNumber)...)
	it := c.sentMsgStorage.db.NewIteratorWithPrefix(prefix)
	for it.Next() {
		blob, err := c.sentMsgStorage.db.Get(it.Key())
		if err != nil {
			continue
		}
		var sData *MsgStorageData
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
	logger.Warnw("lookupSentMsg: message not found")
	return int64(-1)
}

// truncateMsgStored removes all data stored by the block's number
func (c *core) truncateMsgStored() {
	c.sentMsgStorage.mu.Lock()
	defer c.sentMsgStorage.mu.Unlock()

	var (
		key         []byte
		blockNumber = c.currentState.BlockNumber().Uint64()
		logger      = c.sentMsgStorage.logger.With("block", blockNumber)
	)

	key = append(key, []byte(dbKeyMsgStoragePrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	it := c.sentMsgStorage.db.NewIteratorWithPrefix(key)
	for it.Next() {
		if err := c.sentMsgStorage.db.Delete(it.Key()); err != nil {
			logger.Warnw("failed to delete MsgStorage by key", "err", err)
		}
	}
	logger.Info("truncate MsgStorage done", "block", c.currentState.BlockNumber().Uint64())
}

// UnmarshalJSON unmarshals from JSON.
func (h *MsgStorageData) UnmarshalJSON(input []byte) error {
	type sosData MsgStorageData
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
	key = append(key, []byte(dbKeyMsgStoragePrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	key = append(key, []byte(step.String())...)
	key = append(key, encodeInt64(round)...)
	return key
}
