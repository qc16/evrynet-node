#!/bin/sh
# ./deploy/testnet/deploy.sh <path_to_share_volumes> <>

localVolumes=$1
shift
rpccorsdomain=$1
shift
deployExplorer=$1
shift

BASEDIR=$(dirname "$0")

# Stop all dockers gracefully to avoid DB crash
echo "------ Stop dockers ------"
"$BASEDIR"/stop_dockers.sh

# Clear network bridge
yes | docker network prune
# Remove old services
docker rm -f gev-bootnode gev-node-1 gev-node-2 gev-node-3

# Start bootnode
echo "------ Start bootnode ------"
yes | docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-bootnode

# Start nodes
echo "------ Start nodes ------"
yes | shareVolumes=$localVolumes rpccorsdomain=$rpccorsdomain docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-node-1 gev-node-2 gev-node-3
#
## Start explorer
## shellcheck disable=SC2039
#if [[ "$deployExplorer" == "y" ]]; then
#  echo "------ Start explorer ------"
#  rm -rf "$BASEDIR"/explorer/web
#  echo "Cloning explorer from master branch ..."
#  git clone git@github.com:evrynet-official/explorer.git "$BASEDIR"/explorer/web
#
#  docker rmi -f img_explorer
#  docker rm -f gev-explorer
#  yes | gevRPCPort=22001 docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-explorer
#  rm -rf "$BASEDIR"/explorer/web
#fi
