package staking

// Config represents the configuration state of a whisper node.
const (
	WithdrawsStateIndexName    = "withdrawsState"
	CandidateVotersIndexName   = "candidateVoters"
	CandidateDataIndexName     = "candidateData"
	CandidatesIndexName        = "candidates"
	StartBlockIndexName        = "startBlock"
	EpochPeriodIndexName       = "epochPeriod"
	MaxValidatorSizeIndexName  = "maxValidatorSize"
	MinValidatorStakeIndexName = "minValidatorStake"
	MinVoterCapIndexName       = "minVoterCap"
	AdminIndexName             = "admin"
)

// Config represents the configuration state of a whisper node.
type Config struct {
	WithdrawsStateIndex    uint64 //1
	CandidateVotersIndex   uint64 //2
	CandidateDataIndex     uint64 //3
	CandidatesIndex        uint64 //4
	StartBlockIndex        uint64 //5
	EpochPeriodIndex       uint64 //6
	MaxValidatorSizeIndex  uint64 //7
	MinValidatorStakeIndex uint64 //8
	MinVoterCapIndex       uint64 //9
	AdminIndex             uint64 //10
}

// DefaultConfig represents he default configuration.
var DefaultConfig = &Config{
	WithdrawsStateIndex:    1,
	CandidateVotersIndex:   2,
	CandidateDataIndex:     3,
	CandidatesIndex:        4,
	StartBlockIndex:        5,
	EpochPeriodIndex:       6,
	MaxValidatorSizeIndex:  7,
	MinValidatorStakeIndex: 8,
	MinVoterCapIndex:       9,
	AdminIndex:             10,
}
