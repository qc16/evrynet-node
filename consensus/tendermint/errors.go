package tendermint

import "errors"

var (
	// ErrStartedEngine is returned if the engine is already started
	ErrStartedEngine = errors.New("started engine")
)
