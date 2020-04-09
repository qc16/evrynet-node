package staking_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/Evrynetlabs/evrynet-node/core/state/staking"
)

type variableConfig struct {
	Label  string `json:"label"`
	Slot   uint64 `json:"slot,string"`
	Offset uint64 `json:"offset"`
}

type storageLayout struct {
	Storage       []variableConfig `json:"Storage"`
	StructConfigs map[string]struct {
		Label   string           `json:"label"`
		Members []variableConfig `json:"members"`
	} `json:"types"`
}

func TestDefaultConfig(t *testing.T) {
	f, err := os.Open("../../../consensus/staking_contracts/Storage-layout.json")
	require.NoError(t, err)

	data, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	storageJson := gjson.Get(string(data), `contracts.\./EvrynetStaking\.sol.EvrynetStaking.storageLayout`)
	var storageLayout storageLayout
	require.NoError(t, json.Unmarshal([]byte(storageJson.Raw), &storageLayout))
	// test layout config with generated layout
	for _, layout := range storageLayout.Storage {
		switch layout.Label {
		case "withdrawsState":
			require.Equal(t, staking.DefaultConfig.WithdrawsStateLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "candidateVoters":
			require.Equal(t, staking.DefaultConfig.CandidateVotersLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "candidateData":
			require.Equal(t, staking.DefaultConfig.CandidateDataLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "candidates":
			require.Equal(t, staking.DefaultConfig.CandidatesLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "startBlock":
			require.Equal(t, staking.DefaultConfig.StartBlockLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "epochPeriod":
			require.Equal(t, staking.DefaultConfig.EpochPeriodLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "maxValidatorSize":
			require.Equal(t, staking.DefaultConfig.MaxValidatorSizeLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "minValidatorStake":
			require.Equal(t, staking.DefaultConfig.MinValidatorStakeLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "minVoterCap":
			require.Equal(t, staking.DefaultConfig.MinVoterCapLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		case "admin":
			require.Equal(t, staking.DefaultConfig.AdminLayout.Slot, layout.Slot)
			require.Equal(t, uint64(0), layout.Offset)
		}
	}

	//test layout position inside struct
	for _, structCfg := range storageLayout.StructConfigs {
		if structCfg.Label != "struct EvrynetStaking.CandidateData" {
			continue
		}

		for _, member := range structCfg.Members {
			switch member.Label {
			case "totalStake":
				require.Equal(t, staking.DefaultConfig.CandidateDataStruct.TotalStake.Slot, member.Slot)
				require.Equal(t, uint64(0), member.Offset)
			case "owner":
				require.Equal(t, staking.DefaultConfig.CandidateDataStruct.Owner.Slot, member.Slot)
				require.Equal(t, uint64(0), member.Offset)
			case "voterStake":
				require.Equal(t, staking.DefaultConfig.CandidateDataStruct.VotersStakes.Slot, member.Slot)
				require.Equal(t, uint64(0), member.Offset)
			}
		}
	}

}
