package staking_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind/backends"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/staking_contracts"
	"github.com/Evrynetlabs/evrynet-node/core"
	"github.com/Evrynetlabs/evrynet-node/core/state/staking"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/params"
)

func TestGetValidatorsFromStateDb(t *testing.T) {
	var (
		candidates = []common.Address{
			common.HexToAddress("0x560089aB68dc224b250f9588b3DB540D87A66b7a"),
			common.HexToAddress("0x954e4BF2C68F13D97C45db0e02645D145dB6911f"),
		}
		epoch             = big.NewInt(40)
		startBlock        = common.Big0
		maxValidatorSize  = big.NewInt(100)
		minValidatorStake = big.NewInt(20)
		minVoteCap        = big.NewInt(10)
		adminAddr         = common.HexToAddress("0x94F5B16552DCEaCbAdABA146D6e3235f4A8617a8")
		indexLayout       = staking.DefaultConfig
	)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)
	publicKey := privateKey.Public()
	addr := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	be := backends.NewSimulatedBackend(core.GenesisAlloc{
		addr: core.GenesisAccount{
			Balance: big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil),
		},
	}, gasLimit)

	authOpts := bind.NewKeyedTransactor(privateKey)
	authOpts.Nonce = big.NewInt(0)
	authOpts.GasPrice = big.NewInt(params.GasPriceConfig)

	scAddress, tx, _, err := staking_contracts.DeployStakingContracts(authOpts, be, candidates, candidates, epoch, startBlock, maxValidatorSize, minValidatorStake, minVoteCap, adminAddr)
	require.NoError(t, err)

	be.Commit()

	receipt, err := be.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)

	stateDBCaller, err := be.GetStakingCaller(indexLayout)
	require.NoError(t, err)

	validators, err := stateDBCaller.GetValidators(scAddress)
	require.NoError(t, err)
	assert.Equal(t, 2, len(validators))
	for _, val := range validators {
		fmt.Println(val.Hex())
	}
}

func TestCheckIndex(t *testing.T) {
	var (
		candidates = []common.Address{
			common.HexToAddress("0x560089aB68dc224b250f9588b3DB540D87A66b7a"),
			common.HexToAddress("0x954e4BF2C68F13D97C45db0e02645D145dB6911f"),
		}
		epoch             = big.NewInt(40)
		startBlock        = big.NewInt(1)
		maxValidatorSize  = big.NewInt(100)
		minValidatorStake = big.NewInt(20)
		minVoteCap        = big.NewInt(10)
		adminAddr         = common.HexToAddress("0x94F5B16552DCEaCbAdABA146D6e3235f4A8617a8")
	)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)
	publicKey := privateKey.Public()
	addr := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	be := backends.NewSimulatedBackend(core.GenesisAlloc{
		addr: core.GenesisAccount{
			Balance: big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil),
		},
	}, gasLimit)

	authOpts := bind.NewKeyedTransactor(privateKey)
	authOpts.Nonce = big.NewInt(0)
	authOpts.GasPrice = big.NewInt(params.GasPriceConfig)

	scAddress, tx, _, err := staking_contracts.DeployStakingContracts(authOpts, be, candidates, candidates, epoch, startBlock, maxValidatorSize, minValidatorStake, minVoteCap, adminAddr)
	require.NoError(t, err)

	be.Commit()

	receipt, err := be.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)

	stateDB, err := be.CurrentStateDb()
	require.NoError(t, err)

	// startBlock 5
	startBlockData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(5)))
	assert.Equal(t, startBlockData.Big(), startBlock)

	// epoch 6
	epochData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(6)))
	assert.Equal(t, epochData.Big(), epoch)

	// maxValidatorSize 7
	maxValidatorSizeData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(7)))
	assert.Equal(t, maxValidatorSizeData.Big(), maxValidatorSize)

	// minValidatorStake 8
	minValidatorStakeData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(8)))
	assert.Equal(t, minValidatorStakeData.Big(), minValidatorStake)

	// minVoteCap 9
	minVoteCapData := stateDB.GetState(scAddress, common.BigToHash(new(big.Int).SetUint64(9)))
	assert.Equal(t, minVoteCapData.Big(), minVoteCap)
}
