package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/ethdb/memorydb"
	"github.com/evrynet-official/evrynet-client/rlp"
)

func TestCore_StoreAndLookupSentMsg(t *testing.T) {
	var (
		sentMsgStorage = NewMsgStorageData(memorydb.New())
		blockNumber    = uint64(1)
		step           = RoundStepPropose
		step1          = RoundStepPrevote
		round          = int64(0)
		round1         = int64(1)
		round2         = int64(2)
		payload        = []byte{}
	)
	msgData, _ := rlp.EncodeToBytes(payload)
	sentMsgStorage.storeSentMsg(blockNumber, step1, round, msgData)

	sentMsgStorage.storeSentMsg(blockNumber, step, round, msgData)
	sentMsgStorage.storeSentMsg(blockNumber, step, round2, msgData)
	sentMsgStorage.storeSentMsg(blockNumber, step, round1, msgData)

	index := sentMsgStorage.lookupSentMsg(blockNumber, step, round)
	assert.Equal(t, int64(0), index)
	index = sentMsgStorage.lookupSentMsg(blockNumber, step, round1)
	assert.Equal(t, int64(1), index)
	index = sentMsgStorage.lookupSentMsg(blockNumber, step, round2)
	assert.Equal(t, int64(2), index)
	index = sentMsgStorage.lookupSentMsg(blockNumber, step1, round)
	assert.Equal(t, int64(3), index)
}

func TestCore_TruncateSentMsg(t *testing.T) {
	var (
		sentMsgStorage = NewMsgStorageData(memorydb.New())
		blockNumber    = uint64(1)
		blockNumber2   = uint64(2)
		step           = RoundStepPropose
		round          = int64(0)
		payload        = []byte{}
	)
	msgData, _ := rlp.EncodeToBytes(payload)
	sentMsgStorage.storeSentMsg(blockNumber, step, round, msgData)
	sentMsgStorage.storeSentMsg(blockNumber2, step, round, msgData)

	index := sentMsgStorage.lookupSentMsg(blockNumber, step, round)
	assert.Equal(t, int64(0), index)
	index = sentMsgStorage.lookupSentMsg(blockNumber2, step, round)
	assert.Equal(t, int64(0), index)

	sentMsgStorage.truncateMsgStored(blockNumber)
	index = sentMsgStorage.lookupSentMsg(blockNumber, step, round)
	assert.Equal(t, int64(-1), index)

	index = sentMsgStorage.lookupSentMsg(blockNumber2, step, round)
	assert.Equal(t, int64(0), index)
}
