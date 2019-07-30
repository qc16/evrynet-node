#!/bin/sh
set -e

yes | docker rmi -f img_miner1 img_node

echo "Checking to restore genesis block & account"
ancient_file=~/eth-data/miner/ethereum/geth/chaindata/ancient
if [[ -d "$ancient_file" ]]; then
    echo "Backup data at $ancient_file exist => Reuse data"
	yes | needRestore=1 docker-compose -f ./dev/docker-compose.yml up --force-recreate
else
    echo "Backup data at $ancient_file does not exist => Init new data"
	yes | needRestore=0 docker-compose -f ./dev/docker-compose.yml up --force-recreate
fi