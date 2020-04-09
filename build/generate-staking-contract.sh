#!/usr/bin/env bash
set -euo pipefail

pushd ./consensus/staking_contracts
# Generate bytes code and abi files
echo "Generate bytes code and abi files"
docker run -v $(pwd):/staking_contracts ethereum/solc:0.5.13 --overwrite -o /staking_contracts/EvrynetStaking.bin --optimize \
    --abi --bin /staking_contracts/staking-contract/contracts/EvrynetStaking.sol
# Generate go file
echo "Generate go file"
abigen --bin=./EvrynetStaking.bin/EvrynetStaking.bin --abi=./EvrynetStaking.bin/EvrynetStaking.abi --out ./EvrynetStaking.go --pkg=staking_contracts
# Generate storage layout file
echo "Generate storage layout file"
cat solc-input.json | docker run -i  -v $(pwd):/staking_contracts -w /staking_contracts ethereum/solc:0.5.13 --standard-json --allow-paths *, > storage-layout.json
popd