package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//ReadSCBytecode read bytecode of SC
func ReadBytecodeOfSC(path string) (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, "EvrynetStaking.sol"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func ReadABIOfSC(path string) (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, "EvrynetStaking.abi"))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func readSCDir(path string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, info := range fileInfo {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".sol" {
			files = append(files, fmt.Sprintf("%s/%s", path, info.Name()))
		}
	}
	return files, nil
}

func readFirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			return line, nil //Only read first line
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}
