package tendermint

import (
	"time"

	"github.com/Evrynetlabs/evrynet-node/common"
)

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
)

//FaultyMode is the config mode to enable fauty node
type FaultyMode uint64

const (
	// Disabled disables the faulty mode
	Disabled FaultyMode = iota
	// SendFakeProposal sends the proposal with the fake info
	SendFakeProposal
	// RandomlyStopSendingMsg randomly stop message sending
	RandomlyStopSendingMsg
)

func (f FaultyMode) Uint64() uint64 {
	return uint64(f)
}

//Config store all the configuration required for a Tendermint consensus
type Config struct {
	ProposerPolicy        ProposerPolicy   `toml:",omitempty"` // The policy for proposer selection
	Epoch                 uint64           `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
	StakingSCAddress      *common.Address  `toml:",omitempty"` // The staking SC address for validating when deploy SC
	BlockPeriod           uint64           `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	TimeoutPropose        time.Duration    //Duration waiting a propose
	TimeoutProposeDelta   time.Duration    //Increment if timeout happens at propose step to reach eventually synchronous
	TimeoutPrevote        time.Duration    //Duration waiting for more pre-vote after 2/3 received
	TimeoutPrevoteDelta   time.Duration    //Increment if timeout happens at pre-voteWait to reach eventually synchronous
	TimeoutPrecommit      time.Duration    //Duration waiting for more pre-commit after 2/3 received
	TimeoutPrecommitDelta time.Duration    //Duration waiting to increase if pre-commit wait expired to reach eventually synchronous
	TimeoutCommit         time.Duration    //Duration waiting to start round with new height
	FixedValidators       []common.Address // The fixed validators

	FaultyMode uint64 `toml:",omitempty"` // The faulty node indicates the faulty node's behavior
}

var DefaultConfig = &Config{
	ProposerPolicy:        RoundRobin,
	Epoch:                 30000,
	StakingSCAddress:      &common.Address{},
	BlockPeriod:           1,                       // 1 seconds
	TimeoutPropose:        3000 * time.Millisecond, //This is taken from tendermint. Might need tuning
	TimeoutProposeDelta:   500 * time.Millisecond,  //This is taken from tendermint. Might need tuning
	TimeoutPrevote:        1000 * time.Millisecond,
	TimeoutPrevoteDelta:   500 * time.Millisecond,
	TimeoutPrecommit:      1000 * time.Millisecond,
	TimeoutPrecommitDelta: 500 * time.Millisecond,
	TimeoutCommit:         1000 * time.Millisecond,
	FaultyMode:            Disabled.Uint64(),
}

//ProposeTimeout return the timeout for a specific round
//The formula is timeout= TimeoutPropose + round*TimeoutProposeDelta
func (cfg Config) ProposeTimeout(round int64) time.Duration {
	return time.Duration(
		cfg.TimeoutPropose.Nanoseconds()+cfg.TimeoutProposeDelta.Nanoseconds()*(round),
	) * time.Nanosecond
}

// PrevoteTimeout returns the amount of time to wait for straggler votes after receiving any +2/3 pre-votes
func (cfg *Config) PrevoteTimeout(round int64) time.Duration {
	return time.Duration(
		cfg.TimeoutPrevote.Nanoseconds()+cfg.TimeoutPrevoteDelta.Nanoseconds()*int64(round),
	) * time.Nanosecond
}

//Pre-voteCatchupTimeout returns the amount of time to wait for pre-vote msgs before sending catchup msg
//Notes: if node 1 did not receive propose msg, it will delay max = proposeTimeout before sending pre-vote
// So node 2 received propose msg, it will entered pre-vote earlier than node 1 by proposeTimeout
// In here, node 2 sleep about 2 times of proposeTimeout before assuming that sending pre-vote message of node 1 has problem
func (cfg *Config) PrevoteCatchupTimeout(round int64) time.Duration {
	return time.Duration(cfg.ProposeTimeout(round).Nanoseconds() * int64(2))
}

// PrecommitTimeout returns the amount of time to wait for straggler votes after receiving any +2/3 pre-commits
func (cfg *Config) PrecommitTimeout(round int64) time.Duration {
	return time.Duration(
		cfg.TimeoutPrecommit.Nanoseconds()+cfg.TimeoutPrecommitDelta.Nanoseconds()*int64(round),
	) * time.Nanosecond
}

//Pre-commitCatchupTimeout returns the amount of time to wait for pre-commit msgs before sending catchup msg
//Notes: if node 1 did not receive a polka of pre-vote msg, it will delay max = pre-voteWaitTimeout before sending pre-commit
// So node 2 received a polka of pre-vote msg, it will entered pre-commit earlier than node 1 by pre-voteWaitTimeout
// In here, node 2 sleep about 2 times of pre-voteWaitTimeout before assuming that sending pre-commit message of node 1 has problem
func (cfg *Config) PrecommitCatchupTimeout(round int64) time.Duration {
	return time.Duration(cfg.PrevoteTimeout(round).Nanoseconds() * int64(2))
}

// Commit returns the amount of time to wait for straggler votes after receiving +2/3 pre-commits for a single block (ie. a commit).
func (cfg *Config) Commit(t time.Time) time.Time {
	return t.Add(cfg.TimeoutCommit)
}
