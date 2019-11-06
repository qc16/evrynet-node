#!/bin/bash

root=$1  # base directory to use for datadir and logs ./testnet/15
shift
dd=$1  # double digit instance id like 00 01 02
shift

# shellcheck disable=SC2006
datetag=`date "+%c%y%m%d-%H%M%S"|cut -d ' ' -f 6`
datadir=$root/data/$dd        # ./testnet/15/data/01
log=$root/log/$dd.$datetag.log     # ./testnet/15/log/04.2019191101-135434.log
linklog=$root/log/$dd.current.log     # ./testnet/15/log/04.current.log
port=311$dd               # 30304

mkdir -p "$root"/data # ./testnet/15/data
mkdir -p "$root"/log  # ./testnet/15/log
mkdir -p "$root"/nodekey  # ./testnet/15/nodekey
ln -sf "$log" "$linklog"

ancient_file=$datadir/geth/chaindata/ancient
if [ "$(ls -A "$ancient_file")" ]; then
  echo "Backup data at $ancient_file exist => Reuse data"
else
  echo "Backup data at $ancient_file does not exist => Init new data!"

  mkdir -p "$datadir"/geth
  # Reuse nodekey
  if [ "$(ls -A "$root"/nodekey/"$dd")" ]; then
    echo "Has nodekey at $root/nodekey/$dd => Reuse"
    cp "$root"/nodekey/"$dd"/nodekey "$datadir"/geth/
  else
    echo "Has no nodekey at $root/nodekey/$dd => Create new"
    $BOOTNODE --genkey="$datadir"/nodekey
    cp "$datadir"/nodekey "$datadir"/geth/

    echo "Copying keys $datadir/nodekey -> $root/nodekey/$dd"
    mkdir -p "$root"/nodekey/"$dd"
    cp "$datadir"/nodekey "$root"/nodekey/"$dd"/ # ./testnet/15/data/01/nodekey -> ./testnet/15/nodekey/01/
  fi


  # Fist run to create eth.coinbase
  $GEV --datadir "$datadir" init ./fake_genesis.json
  bash -c "$GEV --datadir $datadir --port $port --networkid 15 js <(echo 'console.log(eth.coinbase); exit();') " 2>/dev/null | grep 0x | tee >>./coinbase
  cat ./coinbase

  # Reuse keystore
  if [ "$(ls -A "$root"/keystore/"$dd")" ]; then
    echo "Has keystore at $root/keystore/$dd => Reuse"
    cp -a "$root"/keystore/"$dd"/. "$datadir"/keystore/
  else
    echo "Create an account with password $dd"
    mkdir -p "$root"/keystore/"$dd"
    $GEV --datadir "$datadir" --password <(echo -n "$dd") account new 2>/dev/null | grep 0x | cut -d ' ' -f 8 | tee >>./alloc

    echo "Copying keys $datadir/keystore -> $root/keystore/$dd"
    cp -a "$datadir"/keystore/. "$root"/keystore/"$dd" # ./testnet/15/data/01/keystore -> ./testnet/15/keystore/01
  fi

  # Clear all data except nodekey to keyy old address
  find "$datadir"/geth/ ! -name nodekey -delete
fi