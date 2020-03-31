package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/chilts/sid"
	"github.com/tidwall/gjson"
)

const (
	// ShFile returns the path of file to gen storage layout
	ShFile = "./getStorageLayput.sh"
	// StorageLayoutConfigPath returns the name of index state variables
	StorageLayoutConfigPath = "./indexGenerator.json"
	// KeyPath returns the path of the storage data in json file
	KeyPath = "contracts.EvrynetStaking\\.sol.EvrynetStaking.storageLayout.storage"
)

type storageLayout struct {
	Label  string `json:"label"`
	Offset uint16 `json:"offset"`
	Slot   uint16 `json:"slot,string"`
}

func main() {
	rawFile, err := genJSONFile()
	if err != nil {
		log.Fatalf("generate data from the sol file failed with %s\n", err)
		return
	}
	items, err := readJSONData(rawFile)
	if err != nil {
		log.Fatalf("reads data from json failed with %s\n", err)
		return
	}
	err = saveToFile(items)
	if err != nil {
		log.Fatalf("save data to json failed with %s\n", err)
		return
	}
}

func saveToFile(items []storageLayout) error {
	file, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(StorageLayoutConfigPath, file, 0644)
}

func readJSONData(rawFile string) ([]storageLayout, error) {
	jsonData, err := ioutil.ReadFile(rawFile)
	if err != nil {
		return nil, err
	}
	var items []storageLayout
	data := gjson.Get(string(jsonData), KeyPath)
	if data.Exists() {
		err = json.Unmarshal([]byte(data.Raw), &items)
		if err != nil {
			return nil, err
		}
	}
	if err = os.Remove(rawFile); err != nil {
		log.Printf("remove the temporary file error %s\n", err.Error())
	}

	return items, nil
}

func genJSONFile() (string, error) {
	outputFile := sid.Id()
	cmd := exec.Command("sh", ShFile, outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputFile, nil
}
