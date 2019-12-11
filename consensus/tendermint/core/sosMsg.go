package core

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

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
	if data, _ := c.LookupSentMsg(step, round); data != nil {
		c.sosMsg.logger.Warnw("message had saved at before", "Step", step.String(), "Round", round)
		return
	}
	c.sosMsg.mu.Lock()
	defer c.sosMsg.mu.Unlock()

	msgBytes, err := rlp.EncodeToBytes(msg)
	if err != nil {
		c.sosMsg.logger.Warnw("failed to encode rlp for sentMsg data", "err", err)
		return
	}

	var (
		key         []byte
		blockNumber = c.currentState.BlockNumber().Uint64()
		sData       = &SosData{
			BlockNumber: blockNumber,
			Round:       round,
			Step:        step,
			Data:        msgBytes,
		}
	)

	blob, err := json.Marshal(sData)
	if err != nil {
		c.sosMsg.logger.Warnw("failed to encode rlp for SOSMsg data", "err", err)
		return
	}

	fmt.Printf("StoreSentMsg block number %d  in Step %s at Round %d \n", blockNumber, step.String(), round)
	key = append(key, []byte(dbKeySOSMsgPrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	key = append(key, []byte(step.String())...)
	key = append(key, encodeInt64(round)...)

	if err = c.sosMsg.db.Put(key, blob); err != nil {
		c.sosMsg.logger.Warnw("failed write to SOSMsg file", "err", err)
	}
}

// LookupSentMsg lockups proposal/ vote messages had stored
func (c *core) LookupSentMsg(step RoundStepType, round int64) (*SosData, error) {
	c.sosMsg.mu.Lock()
	defer c.sosMsg.mu.Unlock()
	var (
		key         []byte
		sData       *SosData
		blockNumber = c.currentState.BlockNumber().Uint64()
	)

	key = append(key, []byte(dbKeySOSMsgPrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	key = append(key, []byte(step.String())...)
	key = append(key, encodeInt64(round)...)

	blob, err := c.sosMsg.db.Get(key)
	if err != nil {
		c.sosMsg.logger.Warnw("failed to get SOSMsg by key", "err", err)
		return nil, err
	}
	if err := json.Unmarshal(blob, &sData); err != nil {
		c.sosMsg.logger.Warnw("failed to decode SOSData", "err", err)
		return nil, err
	}
	return sData, nil
}

// TruncateMsgStored removes all data stored by the block's number
func (c *core) TruncateMsgStored() {
	c.sosMsg.mu.Lock()
	defer c.sosMsg.mu.Unlock()

	var (
		key         []byte
		blockNumber = c.currentState.BlockNumber().Uint64()
	)

	key = append(key, []byte(dbKeySOSMsgPrefix)...)
	key = append(key, encodeUint64(blockNumber)...)
	it := c.sosMsg.db.NewIteratorWithPrefix(key)
	for it.Next() {
		if err := c.sosMsg.db.Delete(it.Key()); err != nil {
			c.sosMsg.logger.Warnw("failed to delete SOSMsg by key", "err", err)
		}
	}
	c.sosMsg.logger.Info("truncate SOSMsg done", "block", c.currentState.BlockNumber().Uint64())
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
