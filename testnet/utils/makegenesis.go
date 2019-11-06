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
	bodyCoinbase, err := ioutil.ReadFile("coinbase")
	if err != nil {
		fmt.Print(err)
		return
	}
	listValidator := strings.Split(string(bodyCoinbase), "\n")
	var valAddreses []common.Address
	for _, addr := range listValidator {
		if len(addr) > 0 {
			fmt.Printf("Append validator address: %s", common.HexToAddress(addr).Hex())
			valAddreses = append(valAddreses, common.HexToAddress(addr))
		}
	}

	bodyAlloc, err := ioutil.ReadFile("alloc")
	if err != nil {
		fmt.Print(err)
		return
	}
	listAlloc := strings.Split(string(bodyAlloc), "\n")
	var allocAddreses []common.Address
	for _, addr := range listAlloc {
		if len(addr) > 0 {
			fmt.Printf("Append alloc address: %s", common.HexToAddress(addr).Hex())
			allocAddreses = append(allocAddreses, common.HexToAddress(addr))
		}
	}
	makeGenesis(valAddreses, allocAddreses)
}

// makeGenesis creates a new genesis struct based on some user input.
func makeGenesis(valAddrs []common.Address, allocAddrs []common.Address) {
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
	for _, addr := range allocAddrs {
		// Read the address of the account to fund
		genesis.Alloc[addr] = core.GenesisAccount{
			Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
		}
	}

	// Query the user for some custom extras
	genesis.Config.ChainID = new(big.Int).SetUint64(15)

	// All done, store the genesis and flush to disk
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
