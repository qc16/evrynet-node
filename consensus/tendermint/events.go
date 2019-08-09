package tendermint

// RequestEvent is posted to propose a proposal
type RequestEvent struct {
	Proposal Proposal
}
