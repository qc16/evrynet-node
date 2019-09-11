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
	TimeoutPrevote      time.Duration  //Duration waiting for more prevote after 2/3 received
	TimeoutPrevoteDelta time.Duration  //Increment if timeout happens at prevoteWait to reach eventually synchronous
}

var DefaultConfig = &Config{
	ProposerPolicy:      RoundRobin,
	Epoch:               30000,
	TimeoutPropose:      3000 * time.Millisecond, //This is taken from tendermint. Might need tuning
	TimeoutProposeDelta: 500 * time.Millisecond,  //This is taken from tendermint. Might need tunning
	TimeoutPrevote:      1000 * time.Millisecond,
	TimeoutPrevoteDelta: 500 * time.Millisecond,
}

//ProposeTimeout return the timeout for a specific round
//The formula is timeout= TimeoutPropose + round*TimeoutProposeDelta
func (cfg Config) ProposeTimeout(round int64) time.Duration {
	return time.Duration(
		cfg.TimeoutPropose.Nanoseconds()+cfg.TimeoutProposeDelta.Nanoseconds()*(round),
	) * time.Nanosecond
}

// PrevoteTimeout returns the amount of time to wait for straggler votes after receiving any +2/3 prevotes
func (cfg *Config) PrevoteTimeout(round int64) time.Duration {
	return time.Duration(
		cfg.TimeoutPrevote.Nanoseconds()+cfg.TimeoutPrevoteDelta.Nanoseconds()*int64(round),
	) * time.Nanosecond
}
