#!/bin/sh
set -e

echo "Checking to restore genesis block & account"
ancient_file=~/eth-data/miner/ethereum/geth/chaindata/ancient
if [[ -d "$ancient_file" ]]; then
    echo "Backup data at $ancient_file exist => Reuse data"
    yes | make local_existed
else
    echo "Backup data at $ancient_file does not exist => Init new data"
    yes | make local_first_time
fi