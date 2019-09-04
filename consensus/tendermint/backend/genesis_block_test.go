// This test to init a node with first set of validators
package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/consensus/tendermint"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/rawdb"
	"github.com/evrynet-official/evrynet-client/core/vm"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/event"
	"github.com/stretchr/testify/assert"
)

var (
	nodePKString = "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
)

func TestBackend_Genesis_block(t *testing.T) {
	config, err := makeNodeConfig()
	assert.NoError(t, err)

	nodePK, err := crypto.HexToECDSA(nodePKString)
	assert.NoError(t, err)

	dir, err := ioutil.TempDir("", "eth-chain-genesis")
	assert.NoError(t, err)

	//create db instance with implement leveldb
	db, err := rawdb.NewLevelDBDatabase(dir, 128, 1024, "")
	assert.NoError(t, err)

	//init tendermint backend
	backend := backend{
		config:             config.Tendermint,
		tendermintEventMux: new(event.TypeMux),
		privateKey:         nodePK,
		address:            crypto.PubkeyToAddress(nodePK.PublicKey),
		db:                 db,
	}

	//init tendermint engine
	engine := New(config.Tendermint, nodePK, WithDB(db))

	//set up genesis block
	chainConfig, _, err := core.SetupGenesisBlockWithOverride(db, config.Genesis, nil)
	assert.NoError(t, err)

	//init block chain with tendermint engine
	blockchain, err := core.NewBlockChain(db, nil, chainConfig, engine, vm.Config{}, nil)
	assert.NoError(t, err)

	//take snapshop at the genesis block
	genesisSnapshot, err := backend.snapshot(blockchain, 0, common.Hash{}, nil)
	assert.NoError(t, err)

	valSet := genesisSnapshot.ValSet
	validator := valSet.GetByIndex(0)

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
	return config, nil
}

func getGenesisConf() (*core.Genesis, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Read file genesis generated from pupeth
	genesisFile, err := ioutil.ReadFile(filepath.Join(workingDir, "genesis.json"))
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
