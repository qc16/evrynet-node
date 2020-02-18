package backend

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/tests_utils"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/validator"
	"github.com/Evrynetlabs/evrynet-node/crypto"
)

func TestStakingValidatorInfo_GetValSet(t *testing.T) {
	var (
		nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
		nodeAddr     = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
		validators   = []common.Address{
			nodeAddr,
		}
		genesisHeader = tests_utils.MakeGenesisHeader(validators)
	)
	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	//create New test backend and newMockChain
	cfg := tendermint.DefaultConfig
	cfg.Epoch = 1000
	chain, engine := mustStartTestChainAndBackend(nodePK, genesisHeader, cfg, WithValsetAddresses([]common.Address{}))
	assert.NotNil(t, chain)
	assert.NotNil(t, engine)

	type fields struct {
		Epoch uint64
	}
	type args struct {
		chainReader consensus.ChainReader
		blockNumber *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    tendermint.ValidatorSet
		wantErr bool
	}{
		{
			name: "test with block-number is lower than the epoch",
			args: args{
				chainReader: chain,
				blockNumber: new(big.Int).SetInt64(900),
			},
			want:    validator.NewSet(validators, tendermint.RoundRobin, 900),
			wantErr: false,
		},
		{
			name: "test with block-number is greater than the epoch",
			args: args{
				chainReader: chain,
				blockNumber: new(big.Int).SetInt64(1999),
			},
			want:    validator.NewSet(validators, tendermint.RoundRobin, 1999),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewStakingValidatorInfo(engine)
			got, err := v.GetValSet(tt.args.chainReader, tt.args.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValSetData.GetValSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.GetAddresses(), tt.want.GetAddresses()) {
				t.Errorf("ValSetData.GetValSet() = %v, want %v", got.GetAddresses(), tt.want.GetAddresses())
			}
			if !reflect.DeepEqual(got.Policy(), tt.want.Policy()) {
				t.Errorf("ValSetData.GetValSet() = %v, want %v", got.Policy(), tt.want.Policy())
			}
			if !reflect.DeepEqual(got.GetProposer(), tt.want.GetProposer()) {
				t.Errorf("ValSetData.GetValSet() = %v, want %v", got.GetProposer(), tt.want.GetProposer())
			}
			assert.Equal(t, tt.want.Height(), got.Height())
		})
	}
}
