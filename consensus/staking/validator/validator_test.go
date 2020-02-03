package staking

import (
	"reflect"
	"testing"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
)

func TestValidatorCaller_GetValidators(t *testing.T) {
	type fields struct {
		contract *bind.BoundContract
	}
	type args struct {
		opts        *bind.CallOpts
		blockNumber uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []common.Address
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := &ValidatorCaller{
				contract: tt.fields.contract,
			}
			got, err := val.GetValidators(tt.args.opts, tt.args.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatorCaller.GetValidators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidatorCaller.GetValidators() = %v, want %v", got, tt.want)
			}
		})
	}
}
