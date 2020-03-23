#!/usr/bin/env bash
set -euo pipefail

HOME=$("pwd")/consensus/staking_contracts
docker run -v $HOME:/staking_contracts ethereum/solc:0.5.11 --overwrite -o /staking_contracts/EvrynetStaking.bin --optimize --abi --bin /staking_contracts/EvrynetStaking.sol
./abigen --bin=$HOME/EvrynetStaking.bin/EvrynetStaking.bin --abi=$HOME/EvrynetStaking.bin/EvrynetStaking.abi --out $HOME/EvrynetStaking.go --pkg=staking_contracts