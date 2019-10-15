#!/bin/bash
# Usage:
# bash /path/to/eth-utils/gethup.sh <datadir> <instance_name>

root=$1  # base directory to use for datadir and logs
shift
dd=$1  # double digit instance id like 00 01 02
shift

# shellcheck disable=SC2006
datetag=`date "+%c%y%m%d-%H%M%S"|cut -d ' ' -f 6`
datadir=$root/data/$dd        # /tmp/eth/04
log=$root/log/$dd.$datetag.log     # /tmp/eth/04.2019191101-135434.log
stablelog=$root/log/$dd.log     # /tmp/eth/04.log
#password=$dd              # 04
port=311$dd               # 30304
rpcport=82$dd             # 8104

# shellcheck disable=SC2006
account=`$GETH --datadir="$datadir" account list|head -n1|perl -ne '/([a-f0-9]{40})/ && print $1'`

echo "--- Create genesis block for node $dd ..."
ls -la
$GETH --datadir "$datadir" init ./genesis.json

#echo "--- Static nodes for $dd ..."
#cat "$datadir"/static-nodes.json

# bring up node `dd` (double digit)
# - using <rootdir>/<dd>
# - listening on port 303dd, (like 30300, 30301, ...)
# - with the account unlocked
# - launching json-rpc server on port 81dd (like 8100, 8101, 8102, ...)
echo "$GETH --datadir=$datadir \
  --identity="$dd" \
  --port=$port \
  --unlock=$account \
  --password=<(echo -n $dd) \
  --rpc --rpcport=$rpcport --rpccorsdomain='*' \
  --allow-insecure-unlock $* \
  2>&1 | tee "$stablelog" > "$log" &
"

$GETH --datadir "$datadir" --syncmode "full" --gasprice 1000000000 \
  --identity "$dd" \
  --port "$port" \
  --unlock "$account" \
  --password <(echo -n $dd) \
  --rpc --rpcaddr "0.0.0.0" --rpccorsdomain "*" --rpcvhosts "*" --rpcport "$rpcport" \
  --rpcapi "admin,db,eth,debug,miner,net,shh,txpool,personal,web3" \
  --allow-insecure-unlock $* \
   2>&1 | tee "$stablelog" > "$log" &
