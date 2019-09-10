package tendermint

import (
	"time"
)

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
)

//Config store all the configuration required for a Tendermint consensus
type Config struct {
	ProposerPolicy      ProposerPolicy `toml:",omitempty"` // The policy for proposer selection
	Epoch               uint64         `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
	BlockPeriod         uint64         `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	TimeoutPropose      time.Duration  //Duration waiting a propose
	TimeoutProposeDelta time.Duration  //Increment if timeout happens at propose step to reach eventually synchronous
}

var DefaultConfig = &Config{
	ProposerPolicy:      RoundRobin,
	Epoch:               30000,
	TimeoutPropose:      3000 * time.Millisecond, //This is taken from tendermint. Might need tuning
	TimeoutProposeDelta: 500 * time.Millisecond,  //This is taken from tendermint. Might need tunning
}

//ProposeTimeout return the timeout for a specific round
//The formula is timeout= TimeoutPropose + round*TimeoutProposeDelta
func (cfg Config) ProposeTimeout(round int64) time.Duration {
	return time.Duration(
		cfg.TimeoutPropose.Nanoseconds()+cfg.TimeoutProposeDelta.Nanoseconds()*(round),
	) * time.Nanosecond
}
