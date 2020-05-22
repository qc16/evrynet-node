// +build none

/*

   The mkalloc tool creates the genesis allocation constants in genesis_alloc.go for tendermint genesis configuration
   It outputs a const declaration that contains (address, balance and precompiled Staking Smart Contract).

       go run mkalloc_tendermint.go genesis.json

*/
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	// AllocPath returns the path of alloc data in the genesis config file
	AllocPath = "alloc"
)

// readJSONData reads data of alloc config
func readJSONData(genesisPath string) (string, error) {
	jsonData, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return "", err
	}
	data := gjson.Get(string(jsonData), AllocPath)
	if !data.Exists() {
		return "", errors.New("data at `alloc` not existed")
	}

	input := strings.Replace(data.Raw, "\n", "", -1)
	input = strings.Replace(input, " ", "", -1)
	input = strconv.Quote(input)
	return input, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: mkalloc genesis.json")
		os.Exit(1)
	}

	content, err := readJSONData(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("const allocData =", content)
}
