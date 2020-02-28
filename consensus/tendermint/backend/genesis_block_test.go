// This test to init a node with first set of validators
package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/Evrynetlabs/evrynet-node/common"

	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/stretchr/testify/assert"

	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/backend/fixed_valset_info"
	tendermintCore "github.com/Evrynetlabs/evrynet-node/consensus/tendermint/core"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/tests_utils"
	"github.com/Evrynetlabs/evrynet-node/core"
	"github.com/Evrynetlabs/evrynet-node/core/rawdb"
	"github.com/Evrynetlabs/evrynet-node/core/vm"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/event"
)

var (
	nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
)

func TestBackend_Genesis_block(t *testing.T) {
	backend, blockchain, err := createBlockchainAndBackendFromGenesis()
	assert.NoError(t, err)

	valSet, err := backend.valSetInfo.GetValSet(blockchain, big.NewInt(0))
	assert.NoError(t, err)

	validator := valSet.GetByIndex(0)
	assert.NotNil(t, validator)

	fmt.Println("First set validators")
	fmt.Println(validator)

}

type Config struct {
	Genesis    *core.Genesis
	Tendermint *tendermint.Config
}

func makeNodeConfig() (*Config, error) {
	genesisConf, err := getGenesisConf()
	if err != nil {
		return nil, err
	}
	config := &Config{}
	config.Genesis = genesisConf
	config.Tendermint = &tendermint.Config{}
	config.Tendermint.Epoch = genesisConf.Config.Tendermint.Epoch
	config.Tendermint.FixedValidators = genesisConf.Config.Tendermint.FixedValidators
	config.Tendermint.StakingSCAddress = &common.Address{}
	return config, nil
}

func getGenesisConf() (*core.Genesis, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Read file genesis generated from pupeth
	genesisFile, err := ioutil.ReadFile(filepath.Join(workingDir, "../../../genesis.json"))
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

func createBlockchainAndBackendFromGenesis() (*Backend, *core.BlockChain, error) {
	config, err := makeNodeConfig()
	if err != nil {
		return nil, nil, err
	}

	nodePK, err := crypto.HexToECDSA(nodePKString)
	if err != nil {
		return nil, nil, err
	}

	dir, err := ioutil.TempDir("", "eth-chain-genesis")
	if err != nil {
		return nil, nil, err
	}

	//create db instance with implement leveldb
	db, err := rawdb.NewLevelDBDatabase(dir, 128, 1024, "")
	if err != nil {
		return nil, nil, err
	}

	//init tendermint backend
	backend := &Backend{
		config:               config.Tendermint,
		tendermintEventMux:   new(event.TypeMux),
		privateKey:           nodePK,
		address:              crypto.PubkeyToAddress(nodePK.PublicKey),
		db:                   db,
		mutex:                &sync.RWMutex{},
		storingMsgs:          queue.NewFIFO(),
		dequeueMsgTriggering: make(chan struct{}, 1000),
		broadcastCh:          make(chan broadcastTask),
		valSetInfo:           fixed_valset_info.NewFixedValidatorSetInfo(config.Tendermint.FixedValidators),
	}
	backend.core = tendermintCore.New(backend, config.Tendermint)
	backend.SetBroadcaster(&tests_utils.MockProtocolManager{})
	go backend.dequeueMsgLoop()

	//set up genesis block
	chainConfig, _, err := core.SetupGenesisBlockWithOverride(db, config.Genesis, nil)
	if err != nil {
		return nil, nil, err
	}

	//init block chain with tendermint engine
	blockchain, err := core.NewBlockChain(db, nil, chainConfig, backend, vm.Config{}, nil)
	if err != nil {
		return nil, nil, err
	}
	return backend, blockchain, nil
}
