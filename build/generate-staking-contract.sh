#!/usr/bin/env bash
set -euo pipefail

pushd ./consensus/staking_contracts
git submodule update --init staking-contract
docker run -v $(pwd):/staking_contracts ethereum/solc:0.5.11 --overwrite -o /staking_contracts/EvrynetStaking.bin --optimize --abi --bin /staking_contracts/staking-contract/contracts/EvrynetStaking.sol
abigen --bin=./EvrynetStaking.bin/EvrynetStaking.bin --abi=./EvrynetStaking.bin/EvrynetStaking.abi --out ./EvrynetStaking.go --pkg=staking_contracts
popd