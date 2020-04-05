package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/chilts/sid"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"

	"github.com/Evrynetlabs/evrynet-node/core/state/staking"
)

const (
	// StorageLayoutConfigPath returns the name of index state variables
	StorageLayoutConfigPath = "./indexGenerator.json"
	// KeyPath returns the path of the storage data in json file
	KeyPath = "contracts.EvrynetStaking\\.sol.EvrynetStaking.storageLayout.storage"
)

var (
	shFilePathFlag = cli.StringFlag{
		Name:  "shfilepath",
		Usage: "The path of file to generates storage layout (there are commands to generates storage data layout in this file)",
		Value: "./getStorageLayput.sh",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "indenGenerator"
	app.Usage = "use indexGenerator to generates the storage layout of contract's state variables"
	app.Version = "0.0.1"
	app.Commands = indexGeneratorCmds()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func indexGeneratorCmds() []cli.Command {
	runCmd := cli.Command{
		Action:      run,
		Name:        "run",
		Usage:       "use run command to generates the storage layout of contract's state variables",
		Description: "this tool to generates the storage layout of contract's state variables",
		Flags:       []cli.Flag{shFilePathFlag},
	}

	return []cli.Command{runCmd}
}

func run(ctx *cli.Context) error {
	var (
		shFilePath = ctx.String(shFilePathFlag.Name)
	)

	rawFile, err := genJSONFile(shFilePath)
	if err != nil {
		log.Fatalf("generate data from the sol file failed with %s\n", err)
		return err
	}
	items, err := readJSONData(rawFile)
	if err != nil {
		log.Fatalf("reads data from json failed with %s\n", err)
		return err
	}
	err = saveToFile(items)
	if err != nil {
		log.Fatalf("save data to json failed with %s\n", err)
		return err
	}

	log.Printf("\nThe processing to generate index layout data was successful. You can view data in this file: %s", StorageLayoutConfigPath)
	return nil
}

func saveToFile(items []staking.StorageLayout) error {
	file, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(StorageLayoutConfigPath, file, 0644)
}

func readJSONData(rawFile string) ([]staking.StorageLayout, error) {
	jsonData, err := ioutil.ReadFile(rawFile)
	if err != nil {
		return nil, err
	}
	var items []staking.StorageLayout
	data := gjson.Get(string(jsonData), KeyPath)
	if !data.Exists() {
		return nil, errors.New("storageLayout's data not found")
	}
	err = json.Unmarshal([]byte(data.Raw), &items)
	if err != nil {
		return nil, err
	}
	if err = os.Remove(rawFile); err != nil {
		log.Printf("remove the temporary file error %s\n", err.Error())
	}

	return items, nil
}

func genJSONFile(shFilePath string) (string, error) {
	outputFile := sid.Id()
	cmd := exec.Command("sh", shFilePath, outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputFile, nil
}
