package main

import (
	"fmt"
	"io/ioutil"
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
