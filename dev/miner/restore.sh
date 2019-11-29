#!/bin/sh
set -e

echo "Checking to create genesis block & account"
ancient_file=/root/.ethereum/geth/chaindata/ancient
if [[ $NEED_RESTORE = 1 ]]; then
    echo "Backup data at $ancient_file exist => Reuse data"
else
    echo "Backup data at $ancient_file does not exist => Init new data"
    geth init /root/genesis.json
    geth account import --password /root/.accountpassword /root/.privatekey
fi