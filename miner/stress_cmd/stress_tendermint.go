// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// This file contains a miner stress test based on the Clique consensus engine.
package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/evrynet-official/evrynet-client/accounts/keystore"
	"github.com/evrynet-official/evrynet-client/common/fdlimit"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/eth"
	"github.com/evrynet-official/evrynet-client/eth/downloader"
	"github.com/evrynet-official/evrynet-client/log"
	"github.com/evrynet-official/evrynet-client/miner"
	"github.com/evrynet-official/evrynet-client/node"
	"github.com/evrynet-official/evrynet-client/p2p"
	"github.com/evrynet-official/evrynet-client/p2p/enode"
	"github.com/evrynet-official/evrynet-client/params"
)

func main() {
	var err error
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	fdlimit.Raise(2048)

	enodes, faucets := parseTestConfig("stress_config.json")

	nodePriKey, _ := crypto.GenerateKey()
	// Create a Clique network based off of the Rinkeby config
	genesis, err := makeGenesis("./genesis_testnet.json")
	if err != nil {
		panic(err)
	}

	//make node
	node, err := makeNode(genesis)
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for node.Server().NodeInfo().Ports.Listener == 0 {
		time.Sleep(250 * time.Millisecond)
	}
	// Connect the node to al the previous ones
	for _, n := range enodes {
		node.Server().AddPeer(n)
	}

	// Inject the signer key and start sealing with it
	store := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	signer, err := store.ImportECDSA(nodePriKey, "")
	if err != nil {
		panic(err)
	}
	if err := store.Unlock(signer, ""); err != nil {
		panic(err)
	}

	// wait until node is synced
	time.Sleep(3 * time.Second)
	var ethereum *eth.Ethereum
	if err := node.Service(&ethereum); err != nil {
		panic(err)
	}
	bc := ethereum.BlockChain()
	for !ethereum.Synced() {
		log.Warn("node is not synced, sleeping", "current_block", bc.CurrentHeader().Number)
		time.Sleep(3 * time.Second)
	}

	nonces := make([]uint64, len(faucets))
	// wait for nonce is not change
	for {
		for i, faucet := range faucets {
			log.Info("faucet addr", "addr", faucet)
			addr := crypto.PubkeyToAddress(*(faucet.Public().(*ecdsa.PublicKey)))
			nonces[i] = ethereum.TxPool().State().GetNonce(addr)
		}
		time.Sleep(time.Second * 10)
		var diff = false
		for i, faucet := range faucets {
			log.Info("faucet addr", "addr", faucet)
			addr := crypto.PubkeyToAddress(*(faucet.Public().(*ecdsa.PublicKey)))
			tmp := ethereum.TxPool().State().GetNonce(addr)
			if tmp != nonces[i] {
				diff = true
			}
		}
		if !diff {
			break
		}
	}

	maxBlockNumber := ethereum.BlockChain().CurrentHeader().Number.Uint64()
	numTxs := 0
	start := time.Now()
	// Start injecting transactions from the faucet like crazy
	go func() {
		for {
			currentBlk := bc.CurrentHeader().Number.Uint64()
			for currentBlk > maxBlockNumber {
				maxBlockNumber++
				numTxs += len(bc.GetBlockByNumber(maxBlockNumber).Body().Transactions)
			}
			duration := time.Since(start)
			log.Warn("num tx info", "txs", numTxs, "duration", time.Since(start),
				"txs_per_seconds", float64(numTxs)/duration.Seconds(), "block", currentBlk)
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		index := rand.Intn(len(faucets))

		// Fetch the accessor for the relevant signer
		var ethereum *eth.Ethereum
		if err := node.Service(&ethereum); err != nil {
			panic(err)
		}
		// Create a self transaction and inject into the pool
		tx, err := types.SignTx(types.NewTransaction(nonces[index], crypto.PubkeyToAddress(faucets[index].PublicKey), new(big.Int), 21000, big.NewInt(params.GasPriceConfig), nil), types.HomesteadSigner{}, faucets[index])
		if err != nil {
			panic(err)
		}
		if err := ethereum.TxPool().AddLocal(tx); err != nil {
			panic(err)
		}
		nonces[index]++

		// Wait if we're too saturated
		for {
			pend, queue := ethereum.TxPool().Stats()
			if pend < 2048 {
				break
			}
			log.Info("sleeping tx_pool is full", "pend", pend, "queue", queue)
			time.Sleep(100 * time.Millisecond)
		}

	}
}

type stressConfig struct {
	EnodeStrings  []string `json:"enodes"`
	FaucetStrings []string `json:"faucets"`
}

func parseTestConfig(fileName string) ([]*enode.Node, []*ecdsa.PrivateKey) {
	var cfg stressConfig
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		panic(err)
	}

	log.Info("test config", "enodes", cfg.EnodeStrings, "faucet", cfg.FaucetStrings)
	var (
		enodes  []*enode.Node
		faucets []*ecdsa.PrivateKey
	)
	for _, enodeS := range cfg.EnodeStrings {
		enodes = append(enodes, enode.MustParse(enodeS))
	}

	for _, faucetS := range cfg.FaucetStrings {
		faucetPriKey, err := crypto.HexToECDSA(faucetS)
		if err != nil {
			panic(err)
		}
		faucets = append(faucets, faucetPriKey)
	}
	return enodes, faucets
}

// makeGenesis creates a custom Clique genesis block based on some pre-defined
// signer and faucet accounts.
func makeGenesis(fileName string) (*core.Genesis, error) {
	// Create a Clique network based off of the Rinkeby config
	// Read file genesis generated from pupeth
	genesisFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	config := &core.Genesis{}
	err = json.Unmarshal(genesisFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func makeNode(genesis *core.Genesis) (*node.Node, error) {
	// Define the basic configurations for the Ethereum node
	datadir := "./test_data"

	config := &node.Config{
		Name:    "geth",
		Version: params.Version,
		DataDir: datadir,

		P2P: p2p.Config{
			ListenAddr:  "0.0.0.0:0",
			NoDiscovery: true,
			MaxPeers:    25,
		},
		NoUSB: true,
	}
	// Start the node and configure a full Ethereum node on it
	stack, err := node.New(config)
	if err != nil {
		return nil, err
	}
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &eth.Config{
			Genesis:         genesis,
			NetworkId:       genesis.Config.ChainID.Uint64(),
			GasPrice:        big.NewInt(params.GasPriceConfig),
			SyncMode:        downloader.FullSync,
			DatabaseCache:   256,
			DatabaseHandles: 256,
			TxPool:          core.DefaultTxPoolConfig,
			GPO:             eth.DefaultConfig.GPO,
			Miner: miner.Config{
				GasFloor: genesis.GasLimit * 9 / 10,
				GasCeil:  genesis.GasLimit * 11 / 10,
				GasPrice: genesis.Config.GasPrice,
				Recommit: time.Second,
			},
		})
	}); err != nil {
		return nil, err
	}
	// Start the node and return if successful
	return stack, stack.Start()
}
