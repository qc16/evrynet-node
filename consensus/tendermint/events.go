package tendermint

// RequestEvent is posted to propose a proposal
type RequestEvent struct {
	Proposal Proposal
}

// MessageEvent is posted for Tendermint engine communication
type MessageEvent struct {
	Payload []byte
}
