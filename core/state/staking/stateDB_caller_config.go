package staking

// Constants represents the configuration name of all state variables.
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

// StorageLayout represents the struct of object its get from a json data file
type StorageLayout struct {
	Label  string `json:"label"`
	Offset uint16 `json:"offset"`
	Slot   uint64 `json:"slot,string"`
}

// LayOut represents the Offset and Slot order of a state variable
type LayOut struct {
	Offset uint16
	Slot   uint64
}

// IndexConfigs represents the configuration index of state variables.
type IndexConfigs struct {
	WithdrawsStateIndex    LayOut //1
	CandidateVotersIndex   LayOut //2
	CandidateDataIndex     LayOut //3
	CandidatesIndex        LayOut //4
	StartBlockIndex        LayOut //5
	EpochPeriodIndex       LayOut //6
	MaxValidatorSizeIndex  LayOut //7
	MinValidatorStakeIndex LayOut //8
	MinVoterCapIndex       LayOut //9
	AdminIndex             LayOut //10
}

// DefaultConfig represents he default configuration.
var DefaultConfig = &IndexConfigs{
	WithdrawsStateIndex:    NewLayOut(1, 0),
	CandidateVotersIndex:   NewLayOut(2, 0),
	CandidateDataIndex:     NewLayOut(3, 0),
	CandidatesIndex:        NewLayOut(4, 0),
	StartBlockIndex:        NewLayOut(5, 0),
	EpochPeriodIndex:       NewLayOut(6, 0),
	MaxValidatorSizeIndex:  NewLayOut(7, 0),
	MinValidatorStakeIndex: NewLayOut(8, 0),
	MinVoterCapIndex:       NewLayOut(9, 0),
	AdminIndex:             NewLayOut(10, 0),
}

// NewLayOut returns new instance of a LayOut
func NewLayOut(slot uint64, offset uint16) LayOut {
	return LayOut{Offset: offset, Slot: slot}
}
