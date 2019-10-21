package tendermint

import "errors"

var (
	// ErrStartedEngine is returned if the engine is already started
	ErrStartedEngine = errors.New("engine is already started")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("engine is already stopped")
	// ErrEmptyCommittedSeals is returned if the field of committed seals is zero.
	ErrEmptyCommittedSeals = errors.New("zero committed seals")
	// ErrMismatchTxhashes is returned if the TxHash in header is mismatch.
	ErrMismatchTxhashes = errors.New("mismatch transcations hashes")
)
