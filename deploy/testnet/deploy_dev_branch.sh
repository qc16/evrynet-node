#!/bin/bash
#deploy/testnet/deploy_dev_branch.sh <path_to_share_volumes> <rpc_corsdomain> <deploy_explorer>

localVolumes=$1
shift
rpccorsdomain=$1
shift
deployExplorer=$1
shift

BASEDIR=$(dirname "$0")

# Remove gev-builder
sudo docker rmi -f img_builder
sudo docker rm -f gev-builder

# Start builder
echo "------ Building project ------"
rm -rf "$BASEDIR"/builder/project
rm -rf "$BASEDIR"/builder/bin
rm -rf "$BASEDIR"/bootnode/bin
rm -rf "$BASEDIR"/nodes/bin/gev

echo "--- Cloning evrynet-node from develop branch ..."
git clone --single-branch -b develop --single-branch git@github.com:Evrynetlabs/evrynet-node.git "$BASEDIR"/builder/project

# Stop all dockers gracefully to avoid DB crash
echo "--- Stop docker containers of nodes"
"$BASEDIR"/stop_dockers.sh

# Stop bootnode
sudo docker stop bootnode

# Clear network bridge
yes | sudo docker network prune

echo "--- Building builder container"
yes | sudo version="img_builder" docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-builder

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
  rm -rf "$BASEDIR"/builder/project
  echo "=> Building project successfully!"
  sleep 3

  # Remove old services
  sudo docker rmi -f img_bootnode img_node_1 img_node_2 img_node_3
  sudo docker rm -f gev-bootnode gev-node-1 gev-node-2 gev-node-3

  echo "--- Copy gev, bootnode from gev-builder to ./builder/bin/"
  mkdir "$BASEDIR"/builder/bin/
  sudo docker cp gev-builder:/evrynet/gev "$BASEDIR"/builder/bin/
  sudo docker cp gev-builder:/evrynet/bootnode "$BASEDIR"/builder/bin/

  # Start bootnode
  echo "------ Start bootnode ------"
  mkdir "$BASEDIR"/bootnode/bin/
  cp "$BASEDIR"/builder/bin/bootnode "$BASEDIR"/bootnode/bin/
  yes | sudo docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-bootnode

  # Start nodes
  echo "------ Start nodes ------"
  mkdir "$BASEDIR"/nodes/bin/
  cp "$BASEDIR"/builder/bin/gev "$BASEDIR"/nodes/bin/
  yes | sudo shareVolumes="$localVolumes" rpccorsdomain="$rpccorsdomain" docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-node-1 gev-node-2 gev-node-3

  # Start explorer
  # shellcheck disable=SC2039
  if [[ "$deployExplorer" == "y" ]]; then
    echo "------ Start explorer ------"
    rm -rf "$BASEDIR"/explorer/web
    echo "--- Cloning explorer from master branch ..."
    git clone git@github.com:Evrynetlabs/explorer.git "$BASEDIR"/explorer/web

    sudo docker rmi -f img_explorer
    sudo docker rm -f gev-explorer
    yes | sudo gevRPCPort=22001 docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-explorer
    rm -rf "$BASEDIR"/explorer/web
  fi

  exit 0
else
  rm -rf "$BASEDIR"/builder/project
  echo "--- Building project failed!"
  exit 1
fi