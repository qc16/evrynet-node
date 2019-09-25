package tendermint

import "errors"

var (
	// ErrStartedEngine is returned if the engine is already started
	ErrStartedEngine = errors.New("engine is already started")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("engine is already stopped")
)
