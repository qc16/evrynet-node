#!/bin/sh
echo "------------Clear Data for 4 Test Nodes------------"
# Kill all apps are using port: 30301, 30302, 30303, 30304
sh ./stop_test_nodes.sh

# Init genesis block & Run test node
for i in 1 2 3 4
do
  echo "--- Clear data for node $i ..."
  rm -rf ./tests/test_nodes/node"$i"/data/geth/chaindata
  rm -rf ./tests/test_nodes/node"$i"/data/geth/lightchaindata
  rm -rf ./tests/test_nodes/node"$i"/data/geth/nodes
  rm -r ./tests/test_nodes/node"$i"/data/geth/LOCK
  rm -r ./tests/test_nodes/node"$i"/data/geth/transactions.rlp
  rm -r ./tests/test_nodes/node"$i"/data/geth.ipc
  rm -r ./node"$i".log
done 