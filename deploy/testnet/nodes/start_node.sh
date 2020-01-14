#!/bin/bash

./gev --datadir ./data init ./genesis.json

# shellcheck disable=SC2039
if [[ $HAS_METRIC ]]; then
  if [[ ! $METRICS_ENDPOINT ]]; then
    METRICS_ENDPOINT=http://52.220.52.16:8086
  fi
  if [[ ! $METRICS_USER ]]; then
    # shellcheck disable=SC2209
    METRICS_USER=test
  fi
  if [[ ! $METRICS_PASS ]]; then
    # shellcheck disable=SC2209
    METRICS_PASS=test
  fi

  echo "Start node $NODE_ID with metric!"
  ./gev --datadir ./data --identity "$NODE_ID" --verbosity 4 --tendermint.blockperiod 1 --syncmode full --networkid 15 --mine \
    --rpc --rpcaddr 0.0.0.0 --rpcvhosts "*" --rpcport 8545 --port 30303 \
    --rpcapi admin,db,evr,debug,miner,net,shh,txpool,personal,web3 \
    --bootnodes "enode://$BOOTNODE_ID@$BOOTNODE_IP:30300" \
    --allow-insecure-unlock --unlock "$UNLOCK_ACCOUNT" --password <(echo -n "$UNLOCK_PASS") \
    --pprof --pprofaddr 0.0.0.0 --pprofport 6060 \
    --nodekeyhex "$NODEKEYHEX"\
    --metrics --metrics.influxdb --metrics.influxdb.endpoint "$METRICS_ENDPOINT" --metrics.influxdb.username $METRICS_USER --metrics.influxdb.password $METRICS_PASS 2>>./log/node.log
else
  echo "Start node $NODE_ID!"
  ./gev --datadir ./data --identity "$NODE_ID" --verbosity 4 --tendermint.blockperiod 1 --syncmode full --networkid 15 --mine \
    --rpc --rpcaddr 0.0.0.0 --rpccorsdomain http://"$RPC_CORSDOMAIN":8080 --rpcvhosts "*" --rpcport 8545 --port 30303 \
    --bootnodes "enode://$BOOTNODE_ID@$BOOTNODE_IP:30300" \
    --allow-insecure-unlock --unlock "$UNLOCK_ACCOUNT" --password <(echo -n "$UNLOCK_PASS") \
    --pprof --pprofaddr 0.0.0.0 --pprofport 6060 \
    --nodekeyhex "$NODEKEYHEX"\
    --rpcapi admin,db,evr,debug,miner,net,shh,txpool,personal,web3 2>>./log/node.log
fi