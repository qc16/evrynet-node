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
linklog=$root/log/$dd.current.log     # /tmp/eth/04.current.log
#password=$dd              # 04
port=311$dd               # 30304

mkdir -p "$root"/data
mkdir -p "$root"/log
ln -sf "$log" "$linklog"

#echo "---- Before create keystore ----"
#ls -la "$datadir"/geth/chaindata/ancient

# if we do not have an account, create one
# will not prompt for password, we use the double digit instance id as passwd
# NEVER EVER USE THESE ACCOUNTS FOR INTERACTING WITH A LIVE CHAIN
if [ ! -d "$root/keystore/$dd" ]; then
  echo create an account with password "$dd" [DO NOT EVER USE THIS ON LIVE]
  mkdir -p "$root"/keystore/"$dd"
  $GETH --datadir "$datadir" --password <(echo -n "$dd") account new
# create account with password 00, 01, ...
  # note that the account key will be stored also separately outside
  # datadir
  # this way you can safely clear the data directory and still keep your key
  # under `<rootdir>/keystore/dd

  cp -R "$datadir/keystore" "$root"/keystore/"$dd"
fi

# echo "copying keys $root/keystore/$dd $datadir/keystore"
# ls $root/keystore/$dd/keystore/ $datadir/keystore

# mkdir -p $datadir/keystore
if [ ! -d "$datadir/keystore" ]; then
  echo "copying keys $root/keystore/$dd $datadir/keystore"
  cp -R "$root"/keystore/"$dd"/keystore/ "$datadir"/keystore/
fi

echo "getting address for instance $id ($i/$N)"
eth="$GETH --datadir $datadir --port $port --networkid 15"
# shellcheck disable=SC2089
cmd="$eth js <(echo 'console.log(eth.coinbase); exit();') "
echo "$cmd"
# shellcheck disable=SC2046
bash -c "$cmd" 2>/dev/null | grep 0x | tee >>"$root"/validators

# Clear all data after get validators
rm -rf "$datadir"/geth/chaindata
rm -rf "$datadir"/geth/lightchaindata
rm -rf "$datadir"/geth/nodes
rm -r "$datadir"/geth/LOCK
rm -r "$datadir"/geth/transactions.rlp
rm -r "$datadir"/geth.ipc