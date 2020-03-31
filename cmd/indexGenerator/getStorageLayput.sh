#!/usr/bin/env bash
output=$1
echo '{"language":"Solidity","sources":{"EvrynetStaking.sol":{"urls":["EvrynetStaking.sol"]}},"settings":{"optimizer":{"enabled":true,"runs":200},"evmVersion":"homestead","outputSelection":{"*":{"*":["storageLayout"]}}}}' | solc --optimize --standard-json --allow-paths *, > ${output}