#!/bin/sh
set -e

#./gev --datadir ./data init ./genesis.json

# shellcheck disable=SC2039
if [[ $HAS_METRIC = 1 ]]; then
  echo "Start node $ID with metric! BOOTNODE_ID: $BOOTNODE_ID | BOOTNODE_IP: $BOOTNODE_IP"
  ./gev --datadir ./data --identity "$ID" --verbosity 4 --tendermint.blockperiod 1 --syncmode full --networkid 15 \
    --rpc --rpcaddr 0.0.0.0 --rpcvhosts "*" --rpcport 2200"$ID" --port 3030"$ID" \
    --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3 \
    --bootnodes "enode://$BOOTNODE_ID@$BOOTNODE_IP:30300" \
    --metrics --metrics.influxdb --metrics.influxdb.username test --metrics.influxdb.password test 2>>./log/node_"$ID".log
else
  echo "Start node $ID! BOOTNODE_ID: $BOOTNODE_ID | BOOTNODE_IP: $BOOTNODE_IP"
  nohup ./gev --datadir ./data --identity "$ID" --verbosity 4 --tendermint.blockperiod 1 --syncmode full --networkid 15 \
    --rpc --rpcaddr 0.0.0.0 --rpccorsdomain "http://52.220.52.16:8080" --rpcvhosts "*" --rpcport 2200"$ID" --port 3030"$ID" \
    --bootnodes "enode://$BOOTNODE_ID@$BOOTNODE_IP:30300" \
    --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3 2>>./log/node_"$ID".log
fi

