#!/bin/bash

root=$1
shift
network_id=$1
dir=$root/$network_id
mkdir -p "$dir"/data
mkdir -p "$dir"/log
shift
N=$1
shift

# Export enode list
#if [ ! -f "$dir/nodes" ]; then
#
##  echo "[" >>"$dir"/nodes
#  for ((i = 1; i <= N; ++i)); do
#    # shellcheck disable=SC2006
#    id=$(printf "%02d" $i)
##    if [ ! "$ip_addr" = "" ]; then
##      ip_addr="[::]"
##    fi
#
##    echo "getting enode for instance $id ($i/$N)"
##    eth="$GETH --datadir $dir/data/$id --port 311$id --networkid $network_id"
##    # shellcheck disable=SC2089
##    cmd="$eth js <(echo 'console.log(admin.nodeInfo.enode); exit();') "
##    echo "$cmd"
##    bash -c "$cmd" 2>/dev/null | grep enode | perl -pe "s/\[\:\:\]/$ip_addr/g" | perl -pe "s/^/\"/; s/\s*$/\"/;" | tee >>"$dir"/nodes
##    if ((i < N)); then
##      echo "," >>"$dir"/nodes
##    fi
#
#
##    echo "getting address for instance $id ($i/$N)"
##    eth="$GETH --datadir $dir/data/$id --port 311$id --networkid $network_id"
##    # shellcheck disable=SC2089
##    cmd="$eth js <(echo 'console.log(eth.coinbase); exit();') "
##    echo "$cmd"
##    bash -c "$cmd" 2>/dev/null | grep 0x | tee >>"$dir"/validators
##    if ((i < N)); then
##      echo "," >>"$dir"/validators
##    fi
#
#    # Clear all data after get enode list
#    #    rm -rfv "$dir"/data/"$id"/geth/!("nodekey")
#    #    find "$dir"/data/"$id"/geth -mindepth 1 ! -regex '^/nodekey\(/.*\)?' -delete
#    rm -rf "$dir"/data/"$id"/geth/chaindata
#    rm -rf "$dir"/data/"$id"/geth/lightchaindata
#    rm -rf "$dir"/data/"$id"/geth/nodes
#    rm -r "$dir"/data/"$id"/geth/LOCK
#    rm -r "$dir"/data/"$id"/geth/transactions.rlp
#    rm -r "$dir"/data/"$id"/geth.ipc
#  done
##  echo "]" >>"$dir"/nodes
#fi



for ((i=1;i<=N;++i)); do
  # shellcheck disable=SC2006
  id=$(printf "%02d" $i)
#  echo "copy $dir/data/$id/static-nodes.json"
#  mkdir -p "$dir"/data/"$id"
#  cp "$dir"/nodes "$dir"/data/"$id"/static-nodes.json
  echo "launching node $i/$N ---> tail-f $dir/log/$id.log"
  echo GETH="$GETH" bash ./gev_wrapper_init.sh "$dir" "$id" --networkid "$network_id" $*
  GETH=$GETH bash ./gev_wrapper_init.sh "$dir" "$id" --networkid "$network_id" $*
done

echo "---- Createing genesis.json ----"
cat "$dir"/validators
cp "$dir"/validators ./validators

echo "---- After Createing genesis.json ----"
./makegenesis
cat genesis.json



for ((i=1;i<=N;++i)); do
  # shellcheck disable=SC2006
  id=$(printf "%02d" $i)
  echo "launching node $i/$N ---> tail-f $dir/log/$id.log"
  echo GETH="$GETH" bash ./gev_wrapper_run.sh "$dir" "$id" --networkid "$network_id" $*
  GETH=$GETH bash ./gev_wrapper_run.sh "$dir" "$id" --networkid "$network_id" $*
done