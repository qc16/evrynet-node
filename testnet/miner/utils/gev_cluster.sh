#!/bin/bash

root=$1 #./testnet
shift
network_id=$1
dir=$root/$network_id #./testnet/15
mkdir -p "$dir"/data
mkdir -p "$dir"/log
shift
N=$1
shift

for ((i=1;i<=N;++i)); do
  # shellcheck disable=SC2006
  id=$(printf "%02d" $i)
  echo "Launching node $i/$N ---> tail-f $dir/log/$id.log"
  echo GEV="$GEV" bash ./gev_wrapper_init.sh "$dir" "$id" --networkid "$network_id" $*
  GEV=$GEV bash ./gev_wrapper_init.sh "$dir" "$id" --networkid "$network_id" $*
done

echo "List coinbase"
cat ./coinbase

echo "List alloc"
cat ./alloc

echo "Createing genesis.json"
./makegenesis
cat genesis.json



for ((i=1;i<=N;++i)); do
  # shellcheck disable=SC2006
  id=$(printf "%02d" $i)
  echo "Launching node $i/$N ---> tail-f $dir/log/$id.log"
  echo GEV="$GEV" bash ./gev_wrapper_run.sh "$dir" "$id" --networkid "$network_id" $*
  GEV=$GEV bash ./gev_wrapper_run.sh "$dir" "$id" --networkid "$network_id" $*
done