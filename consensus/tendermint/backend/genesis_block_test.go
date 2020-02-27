// This test to init a node with first set of validators
package backend

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"

	queue "github.com/enriquebris/goconcurrentqueue"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/assert"

	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/backend/fixed_valset_info"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/backend/staking"
	"github.com/Evrynetlabs/evrynet-node/consensus/tendermint/tests_utils"
	"github.com/Evrynetlabs/evrynet-node/core"
	"github.com/Evrynetlabs/evrynet-node/core/rawdb"
	coreStaking "github.com/Evrynetlabs/evrynet-node/core/state/staking"
	"github.com/Evrynetlabs/evrynet-node/core/vm"
	"github.com/Evrynetlabs/evrynet-node/crypto"
	"github.com/Evrynetlabs/evrynet-node/event"
	"github.com/Evrynetlabs/evrynet-node/log"
)

type GenesisType string

const (
	StakingSC       GenesisType = "../tests/genesis_staking_sc.json"       // 1 validator
	FixedValidators GenesisType = "../tests/genesis_fixed_validators.json" // 4 validators
)

var (
	nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
)

func TestGenesisblockWithStakingSC(t *testing.T) {
	testCases := []struct {
		name        string
		genesisType GenesisType
		validators  int
	}{
		{
			name:        "StakingSC",
			genesisType: StakingSC,
			validators:  4,
		},
		{
			name:        "FixedValidators",
			genesisType: FixedValidators,
			validators:  1,
		},
	}
	for _, tc := range testCases {
		getValidators := func(t *testing.T) {
			backend, blockchain, err := createBlockchainAndBackendFromGenesis(tc.genesisType)
			assert.NoError(t, err)

			// Tested with 4 valset but it will break the test TestBackend_HandleMsg (not enough 2f+1)
			// So I only test 1 valset
			valSet, err := backend.valSetInfo.GetValSet(blockchain, big.NewInt(0))
			assert.NoError(t, err)
			assert.Equal(t, tc.validators, len(valSet.List()))

			valSet2 := backend.Validators(big.NewInt(0))
			assert.Equal(t, tc.validators, len(valSet2.List()))

			validator := valSet.GetByIndex(0)
			assert.NotNil(t, validator)
		}
		t.Run(tc.name, getValidators)
	}
}

func TestBackendCallGetListCandidateFromSC(t *testing.T) {
	// Must init log to show error when using log.Debug
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

	backend, blockchain, err := createBlockchainAndBackendFromGenesis(StakingSC)
	assert.NoError(t, err)

	state, err := backend.chain.StateAt(backend.CurrentHeadBlock().Root())
	assert.NoError(t, err)

	header := backend.chain.CurrentHeader()
	stakingCaller := coreStaking.NewStakingCaller(state, blockchain, header, backend.chain.Config(), vm.Config{})
	validators, err := stakingCaller.GetValidators(backend.stakingContractAddr)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(validators))
}

type Config struct {
	Genesis    *core.Genesis
	Tendermint *tendermint.Config
}

func makeNodeConfig(g GenesisType) (*Config, error) {
	genesisConf, err := getGenesisConf(g)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	config.Genesis = genesisConf
	config.Tendermint = tendermint.DefaultConfig
	config.Tendermint.ProposerPolicy = tendermint.ProposerPolicy(genesisConf.Config.Tendermint.ProposerPolicy)
	config.Tendermint.Epoch = genesisConf.Config.Tendermint.Epoch
	config.Tendermint.FixedValidators = genesisConf.Config.Tendermint.FixedValidators
	config.Tendermint.StakingSCAddress = genesisConf.Config.Tendermint.StakingSCAddress
	return config, nil
}

func getGenesisConf(g GenesisType) (*core.Genesis, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Read file genesis generated from pupeth
	genesisFile, err := ioutil.ReadFile(filepath.Join(workingDir, string(g)))
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

func createBlockchainAndBackendFromGenesis(g GenesisType) (*Backend, *core.BlockChain, error) {
	config, err := makeNodeConfig(g)
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
	valSetCache, _ := lru.NewARC(inMemoryValset)
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
		computedValSetCache:  valSetCache,
	}

	if config.Tendermint.FixedValidators != nil && len(config.Tendermint.FixedValidators) > 0 {
		backend.valSetInfo = fixed_valset_info.NewFixedValidatorSetInfo(config.Tendermint.FixedValidators)
	} else {
		backend.valSetInfo = staking.NewStakingValidatorInfo(config.Tendermint.Epoch, config.Tendermint.ProposerPolicy)
		backend.stakingContractAddr = *config.Tendermint.StakingSCAddress
	}

	//backend.core = tendermintCore.New(backend, config.Tendermint)
	backend.SetBroadcaster(&tests_utils.MockProtocolManager{})
	go backend.dequeueMsgLoop()

	//set up genesis block
	chainConfig, _, err := core.SetupGenesisBlock(db, config.Genesis)
	if err != nil {
		return nil, nil, err
	}

	//init block chain with tendermint engine
	blockchain, err := core.NewBlockChain(db, nil, chainConfig, backend, vm.Config{}, nil)
	if err != nil {
		return nil, nil, err
	}
	backend.chain = blockchain
	backend.currentBlock = blockchain.CurrentBlock
	return backend, blockchain, nil
}
