#!/bin/bash
#deploy/testnet/deploy_by_version.sh <path_to_share_volumes> <rpc_corsdomain> <deploy_explorer> <version>

localVolumes=$1
shift
rpccorsdomain=$1
shift
deployExplorer=$1
shift
version=$1
shift

if [ "$localVolumes" == "" ] || [ "$rpccorsdomain" == "" ] || [ "$deployExplorer" == "" ] || [ "$version" == "" ]
then
  echo 'Missing params'
  exit 1
fi

BASEDIR=$(dirname "$0")
IMAGE_TAG="registry.gitlab.com/evry/evrynet-client"

# Start builder
echo "--- Cleaning executable files ..."
rm -rf "$BASEDIR"/builder/bin
rm -rf "$BASEDIR"/bootnode/bin
rm -rf "$BASEDIR"/nodes/bin/gev

# Stop all dockers gracefully to avoid DB crash
echo "--- Stop docker containers of nodes"
"$BASEDIR"/stop_dockers.sh

# Stop bootnode
sudo docker stop bootnode

# Clear network bridge
yes | sudo docker network prune

echo "--- Pulling $IMAGE_TAG image"
sudo docker pull "$IMAGE_TAG:$version"

echo "--- Checking $IMAGE_TAG:$version image"
GrepInfo=$(sudo docker images -a | grep "$IMAGE_TAG" | grep "$version" | awk '{print $1}')
echo "$GrepInfo"

# shellcheck disable=SC2181
if [ "$GrepInfo" != "" ];
then
  echo "=> Pulling $IMAGE_TAG:$version image successfully!"
  sleep 3

  # Remove old services
  echo '--- Removing old services'
  sudo docker rmi -f img_bootnode img_node_1 img_node_2 img_node_3
  sudo docker rm -f gev-bootnode gev-node-1 gev-node-2 gev-node-3

  echo "--- Creating temporary docker container from $IMAGE_TAG to get bin files"
  TempBuilderContainer=$(docker create "$IMAGE_TAG:$version")

  echo "--- Copy gev, bootnode from $IMAGE_TAG to ./builder/bin/"
  mkdir "$BASEDIR"/builder/bin/
  sudo docker cp "$TempBuilderContainer":/evrynet/gev "$BASEDIR"/builder/bin/
  sudo docker cp "$TempBuilderContainer":/evrynet/bootnode "$BASEDIR"/builder/bin/
  docker rm -v "$TempBuilderContainer"

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