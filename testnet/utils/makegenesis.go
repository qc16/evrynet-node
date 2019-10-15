package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/params"
	"github.com/evrynet-official/evrynet-client/rlp"
)

type config struct {
	path    string        // File containing the configuration values
	Genesis *core.Genesis `json:"genesis,omitempty"` // Genesis block to cache for node deploys
}

func main() {
	body, err := ioutil.ReadFile("validators")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(string(body))
	listValidator := strings.Split(string(body), "\n")
	fmt.Printf("---- listValidator len: %d", len(listValidator))
	var valAddreses []common.Address
	for _, addr := range listValidator {
		if len(addr) > 0 {
			valAddreses = append(valAddreses, common.HexToAddress(addr))
		}
	}
	makeGenesis(valAddreses)
}

// makeGenesis creates a new genesis struct based on some user input.
func makeGenesis(valAddrs []common.Address) {
	// Construct a default genesis block
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(524288),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			HomesteadBlock:      big.NewInt(1),
			EIP150Block:         big.NewInt(2),
			EIP155Block:         big.NewInt(3),
			EIP158Block:         big.NewInt(3),
			ByzantiumBlock:      big.NewInt(4),
			ConstantinopleBlock: big.NewInt(5),
			Tendermint: &params.TendermintConfig{
				Epoch:          uint64(30000),
				ProposerPolicy: uint64(0),
			},
		},
	}

	// In the case of Tendermint, configure the consensus parameters
	genesis.Difficulty = big.NewInt(1)

	// We also need the initial list of validators
	fmt.Println()
	fmt.Println("Which accounts are validators? (mandatory at least one)")

	tendermintExtra := types.TendermintExtra{
		Validators: valAddrs,
	}
	extraData, err := rlp.EncodeToBytes(&tendermintExtra)
	if err != nil {
		fmt.Println("rlp encode got error", "error", err)
		return
	}
	tendermintExtraVanity := bytes.Repeat([]byte{0x00}, types.TendermintExtraVanity)
	genesis.ExtraData = append(tendermintExtraVanity, extraData...)

	// Consensus all set, just ask for initial funds and go
	fmt.Println()
	fmt.Println("Which accounts should be pre-funded? (advisable at least one)")
	for _, addr := range valAddrs {
		// Read the address of the account to fund
		genesis.Alloc[addr] = core.GenesisAccount{
			Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
		}
	}

	// Query the user for some custom extras
	fmt.Println()
	fmt.Println("Specify your chain/network ID if you want an explicit one (default = random)")
	genesis.Config.ChainID = new(big.Int).SetUint64(15)

	// All done, store the genesis and flush to disk
	fmt.Println("Configured new genesis block")

	config := config{
		path:    "./genesis.json",
		Genesis: genesis}
	config.flush()
}

// flush dumps the contents of config to disk.
func (c config) flush() {
	os.MkdirAll(filepath.Dir(c.path), 0755)

	out, _ := json.MarshalIndent(c.Genesis, "", "  ")
	if err := ioutil.WriteFile(c.path, out, 0644); err != nil {
		fmt.Println("Failed to save puppeth configs", "file", c.path, "err", err)
	}
}
