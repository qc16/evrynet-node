package staking_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind/backends"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/staking_contracts"
	"github.com/Evrynetlabs/evrynet-node/core"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/params"
)

const (
	privateKeyHex = "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1"
	gasLimit      = 10000000
)

func TestGetValidators(t *testing.T) {
	var (
		candidates = []common.Address{
			common.HexToAddress("0x560089aB68dc224b250f9588b3DB540D87A66b7a"),
			common.HexToAddress("0x954e4BF2C68F13D97C45db0e02645D145dB6911f"),
		}
		epoch             = big.NewInt(300000)
		startBlock        = common.Big0
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

	addr, tx, _, err := staking_contracts.DeployStakingContracts(authOpts, be, candidates, candidates[0], epoch, startBlock, maxValidatorSize, minValidatorStake, minVoteCap, adminAddr)
	require.NoError(t, err)

	be.Commit()

	receipt, err := be.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)

	stakingCaller, err := be.GetStakingCaller()
	require.NoError(t, err)

	validators, err := stakingCaller.GetValidators(addr)
	require.NoError(t, err)
	for _, val := range validators {
		fmt.Println(val.Hex())
	}
}
