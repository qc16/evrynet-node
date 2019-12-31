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
