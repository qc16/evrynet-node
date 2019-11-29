#!/bin/sh
# ./deploy/testnet/deploy.sh <path_to_share_volumes> <>

localVolumes=$1
shift
rpccorsdomain=$1
shift
deployExplorer=$1
shift

BASEDIR=$(dirname "$0")

# Remove gev-builder
docker rmi -f img_builder
docker rm -f gev-builder

# Start builder
echo "------ Building project ------"
rm -rf "$BASEDIR"/builder/project
rm -rf "$BASEDIR"/builder/bin
rm -rf "$BASEDIR"/bootnode/bin
rm -rf "$BASEDIR"/nodes/bin/gev
echo "Cloning evrynet-client from develop branch ..."
git clone --single-branch -b develop --single-branch git@github.com:evrynet-official/evrynet-client.git "$BASEDIR"/builder/project
yes | docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-builder

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
  rm -rf "$BASEDIR"/builder/project
  echo "=> Building project successfully!"
  sleep 3

  # Stop all dockers gracefully to avoid DB crash
  echo "------ Stop dockers ------"
  "$BASEDIR"/stop_dockers.sh

  # Clear network bridge
  yes | docker network prune
  # Remove old services
  docker rmi -f img_bootnode img_node_1 img_node_2 img_node_3
  docker rm -f gev-bootnode gev-node-1 gev-node-2 gev-node-3

  echo "Copy gev, bootnode from gev-builder to ./builder/bin/"
  mkdir "$BASEDIR"/builder/bin/
  docker cp gev-builder:/evrynet/gev "$BASEDIR"/builder/bin/
  docker cp gev-builder:/evrynet/bootnode "$BASEDIR"/builder/bin/

  # Start bootnode
  echo "------ Start bootnode ------"
  mkdir "$BASEDIR"/bootnode/bin/
  cp "$BASEDIR"/builder/bin/bootnode "$BASEDIR"/bootnode/bin/
  yes | docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-bootnode

  # Start nodes
  echo "------ Start nodes ------"
  mkdir "$BASEDIR"/nodes/bin/
  cp "$BASEDIR"/builder/bin/gev "$BASEDIR"/nodes/bin/
  yes | shareVolumes=$localVolumes rpccorsdomain=$rpccorsdomain docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-node-1 gev-node-2 gev-node-3

  # Start explorer
  # shellcheck disable=SC2039
  if [[ "$deployExplorer" == "y" ]]; then
    echo "------ Start explorer ------"
    rm -rf "$BASEDIR"/explorer/web
    echo "Cloning explorer from master branch ..."
    git clone git@github.com:evrynet-official/explorer.git "$BASEDIR"/explorer/web

    docker rmi -f img_explorer
    docker rm -f gev-explorer
    yes | gevRPCPort=22001 docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-explorer
    rm -rf "$BASEDIR"/explorer/web
  fi

  exit 0
else
  rm -rf "$BASEDIR"/builder/project
  echo "Building project failed!"
  exit 1
fi