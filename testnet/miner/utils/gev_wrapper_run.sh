#!/bin/bash

root=$1  # base directory to use for datadir and logs
shift
dd=$1  # double digit instance id like 00 01 02
shift

# shellcheck disable=SC2006
datetag=`date "+%c%y%m%d-%H%M%S"|cut -d ' ' -f 6`
datadir=$root/data/$dd        # /tmp/eth/04
log=$root/log/$dd.$datetag.log     # /tmp/eth/04.2019191101-135434.log
stablelog=$root/log/$dd.log     # /tmp/eth/04.log
port=311$dd               # 30304
rpcport=82$dd             # 8104


ancient_file=$datadir/geth/chaindata/ancient
if [ "$(ls -A "$ancient_file")" ]; then
  echo "Backup data at $ancient_file exist => Reuse data"
else
  # NOTE: after this step coinbase will be changed
  echo "Create genesis block for node $dd ..."
  $GEV --datadir "$datadir" init ./genesis.json
fi


echo "$GEV --datadir=$datadir --syncmode full --gasprice 1000000000  \
  --identity $dd \
  --port $port \
  --rpc --rpcaddr 0.0.0.0 --rpccorsdomain * --rpcvhosts * --rpcport $rpcport  \
  --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3 \
  --allow-insecure-unlock $* \
  2>&1 | tee "$stablelog" > "$log" &
"

$GEV --datadir "$datadir" --syncmode "full" --gasprice 1000000000 \
  --identity "$dd" \
  --port "$port" \
  --rpc --rpcaddr "0.0.0.0" --rpccorsdomain "*" --rpcvhosts "*" --rpcport "$rpcport" \
  --rpcapi "admin,db,eth,debug,miner,net,shh,txpool,personal,web3" \
  --allow-insecure-unlock $* \
   2>&1 | tee "$stablelog" > "$log" &
