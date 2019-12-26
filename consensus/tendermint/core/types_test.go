package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/evrynet-official/evrynet-client/rlp"
)

func TestCatchUpMsg_DecodeRLP(t *testing.T) {
	msg := CatchUpMsg{
		BlockNumber: big.NewInt(5),
		Round:       3,
		Step:        RoundStepPrecommit,
	}
	data, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	var decodedMsg CatchUpMsg
	require.NoError(t, rlp.DecodeBytes(data, &decodedMsg))
	require.Equal(t, msg, decodedMsg)
}

func TestResendMsg_DecodeRLP(t *testing.T) {
	var (
		payload1 = []byte("aaaaa")
		payload2 = []byte("cc")
	)
	resendMsg := ResendMsg{
		Payloads:    append([][]byte{}, payload1, payload2),
		BlockNumber: big.NewInt(24),
	}
	bs, err := rlp.EncodeToBytes(&resendMsg)
	require.NoError(t, err)
	var newResendMsg ResendMsg
	require.NoError(t, rlp.DecodeBytes(bs, &newResendMsg))

	require.Equal(t, uint64(24), newResendMsg.BlockNumber.Uint64())
	require.Equal(t, payload1, newResendMsg.Payloads[0])
	require.Equal(t, payload2, newResendMsg.Payloads[1])
}
