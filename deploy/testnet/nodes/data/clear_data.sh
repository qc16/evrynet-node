#!/bin/sh
echo "------------Clear Data for 3 Test Nodes------------"
BASEDIR=$(dirname "$0")

for i in 1 2 3
do
  echo "--- Clear data for node $i ..."
  rm -rf "$BASEDIR"/node_"$i"/data/geth/chaindata
  rm -rf "$BASEDIR"/node_"$i"/data/geth/lightchaindata
  rm -rf "$BASEDIR"/node_"$i"/data/geth/nodes
  rm -r "$BASEDIR"/node_"$i"/data/geth/LOCK
  rm -r "$BASEDIR"/node_"$i"/data/geth/transactions.rlp
  rm -r "$BASEDIR"/node_"$i"/data/geth.ipc
  rm -r "$BASEDIR"/node_"$i"/log/*.log
done 
