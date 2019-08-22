package tendermint

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
)

type Config struct {
	ProposerPolicy ProposerPolicy `toml:",omitempty"` // The policy for proposer selection
	BlockPeriod    uint64         `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	Epoch          uint64         `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
}

var DefaultConfig = &Config{
	ProposerPolicy: RoundRobin,
	Epoch:          30000,
}
