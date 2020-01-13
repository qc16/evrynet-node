package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeoutInfo(t *testing.T) {
	ti1 := timeoutInfo{
		BlockNumber: big.NewInt(209),
		Round:       1,
		Step:        RoundStepPropose,
		Duration:    time.Second,
		Retry:       0,
	}

	ti2 := timeoutInfo{
		BlockNumber: big.NewInt(209),
		Round:       0,
		Step:        RoundStepPrecommitWait,
		Duration:    time.Second,
		Retry:       0,
	}

	ti3 := timeoutInfo{
		BlockNumber: big.NewInt(209),
		Round:       1,
		Step:        RoundStepPrecommit,
		Duration:    time.Second,
		Retry:       0,
	}

	ti4 := timeoutInfo{
		BlockNumber: big.NewInt(208),
		Round:       1,
		Step:        RoundStepPrecommit,
		Duration:    time.Second,
		Retry:       0,
	}

	require.Equal(t, false, ti1.earlierOrEqual(ti2))
	require.Equal(t, true, ti1.earlierOrEqual(ti3))
	require.Equal(t, false, ti1.earlierOrEqual(ti4))
}
