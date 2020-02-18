package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/Evrynetlabs/evrynet-node/accounts/abi/bind"
	"github.com/Evrynetlabs/evrynet-node/common"
	"github.com/Evrynetlabs/evrynet-node/consensus/staking_contracts"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/evrclient"
	"github.com/Evrynetlabs/evrynet-node/params"
)

//TODO: remove this after create staking SC in genesis block is done

const (
	privateKeyHex = "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1"
)

// create a smart contract for staking test
// addr: 0x14229bC33F417cd5470e09b41DF41ce23dE0D06b
func main() {
	var (
		candidates = []common.Address{
			common.HexToAddress("0x560089aB68dc224b250f9588b3DB540D87A66b7a"),
			common.HexToAddress("0x954e4BF2C68F13D97C45db0e02645D145dB6911f"),
		}
		epoch             = big.NewInt(300000)
		maxValidatorSize  = big.NewInt(100)
		minValidatorStake = big.NewInt(20)
		minVoteCap        = big.NewInt(10)
	)

	client, err := evrclient.Dial("http://127.0.0.1:22001")
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.Public()
	fromAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))
	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		panic(err)
	}
	if nonce != 0 {
		panic("nonce should be zero")
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(0)
	auth.GasPrice = big.NewInt(params.GasPriceConfig)

	addr, tx, _, err := staking_contracts.DeployStakingContracts(auth, client, candidates, candidates[0], epoch, maxValidatorSize, minValidatorStake, minVoteCap)
	if err != nil {
		panic(err)
	}

	for {
		if receipt, err := client.TransactionReceipt(context.Background(), tx.Hash()); err == nil {
			if receipt.Status == 1 {
				break
			} else {
				panic("transaction is not success")
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Println("staking SC is deployed, addr", addr.Hex())
}
