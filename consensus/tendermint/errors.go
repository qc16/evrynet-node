package tendermint

import "errors"

var (
	// ErrStartedEngine is returned if the engine is already started
	ErrStartedEngine = errors.New("started engine")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
	// ErrFailedDecodeProposal is returned when the proposal message is malformed.
	ErrFailedDecodeProposal = errors.New("failed to decode Proposal")
	// ErrIncorrectProposer is returned when received message is from incorrect proposer
	ErrIncorrectProposer = errors.New("message does not come from correct proposer")
	// ErrInvalidProposal is returned when a proposal is malformed.
	ErrInvalidProposal = errors.New("invalid proposal")
	// ErrMismatchTxhashes is returned if the TxHash in header is mismatch.
	ErrMismatchTxhashes = errors.New("mismatch transactions hashes")
	// ErrInvalidUncleHash is returned if a block contains an non-empty uncle list.
	ErrInvalidUncleHash = errors.New("non empty uncle hash")
	// ErrFutureMessage is returned when current view is earlier than the
	// view of the received message.
	ErrFutureMessage = errors.New("future message")
	// ErrOldMessage is returned when the received message's view is earlier
	// than current view.
	ErrOldMessage = errors.New("old message")
	// ErrInvalidMessage is returned when the message is malformed.
	ErrInvalidMessage = errors.New("invalid message")
)
