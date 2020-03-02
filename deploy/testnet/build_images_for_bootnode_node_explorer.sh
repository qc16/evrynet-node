#!/bin/bash

# This script will clone, build & push images of bootnode, node to docker hub by version
# ./build_images_for_bootnode_node_explorer.sh <tag_version_or_develop_branch> <environment> <build_explorer>
# Ex: deploy/testnet/build_images_for_bootnode_node_explorer.sh develop testnet n


version=$1
shift
env=$1
shift
buildExplorer=$1
shift


if [[ "$version" == "" || "$env" == "" || "$buildExplorer" == "" ]]
then
  echo 'Missing params'
  exit 1
fi

BASEDIR=$(dirname "$0")
BUILDER_REPOSITORY="kybernetwork/evrynet-builder"
BUILDER_TAG_ENV="$BUILDER_REPOSITORY:$version-$env"

BOOTNODE_REPOSITORY="kybernetwork/evrynet-bootnode"
BOOTNODE_TAG_ENV="$BOOTNODE_REPOSITORY:$version-$env"

NODE_REPOSITORY="kybernetwork/evrynet-node"
NODE_TAG_ENV="$NODE_REPOSITORY:$version-$env"

EXPLORER_REPOSITORY="kybernetwork/evrynet-explorer"
EXPLORER_TAG_ENV="$EXPLORER_REPOSITORY:$version-$env"


git fetch --all --tags --prune
rm -rf "$BASEDIR"/builder/project
if [ "$version" == "develop" ];
then
  echo "--- Cloning evrynet-node from develop branch ..."
  git clone --single-branch -b develop --single-branch git@github.com:Evrynetlabs/evrynet-node.git "$BASEDIR"/builder/project
else
  echo "--- Cloning evrynet-node from tag $version ..."
  git clone -b "$version" git@github.com:Evrynetlabs/evrynet-node.git "$BASEDIR"/builder/project
fi


echo "--- Building builder container for version $version"
#yes | sudo version="$BUILDER_TAG_ENV" docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-builder
yes | sudo docker build -f ./deploy/testnet/builder/Dockerfile -t "$BUILDER_TAG_ENV" ./deploy/testnet/builder/

echo "--- Creating temporary docker container from $BUILDER_TAG_ENV to get bin files"
TempBuilderContainer=$(sudo docker create "$BUILDER_TAG_ENV")

echo "--- Copy gev, bootnode from $BUILDER_TAG_ENV to ./builder/bin/"
rm -rf "$BASEDIR"/nodes/bin/gev
rm -rf "$BASEDIR"/bootnode/bin/bootnode
sudo docker cp "$TempBuilderContainer":/evrynet/gev "$BASEDIR"/nodes/bin/
sudo docker cp "$TempBuilderContainer":/evrynet/bootnode "$BASEDIR"/bootnode/bin/
sudo docker rm -v "$TempBuilderContainer"

echo "--- Building bootnode image $BOOTNODE_TAG_ENV "
yes | sudo docker build -f ./deploy/testnet/bootnode/Dockerfile -t "$BOOTNODE_TAG_ENV" ./deploy/testnet/bootnode/

echo "--- Building node image $NODE_TAG_ENV "
yes | sudo docker build -f ./deploy/testnet/nodes/Dockerfile -t "$NODE_TAG_ENV" ./deploy/testnet/nodes/

rm -rf "$BASEDIR"/builder/project


if [[ "$buildExplorer" == "y" ]]; then
  rm -rf "$BASEDIR"/explorer/web
  echo "--- Cloning explorer from master branch ..."
  git clone git@github.com:Evrynetlabs/explorer.git "$BASEDIR"/explorer/web

  sudo docker rmi -f img_explorer
  sudo docker rm -f gev-explorer

  echo "--- Building explorer image $EXPLORER_TAG_ENV "
  yes | sudo docker build -f ./deploy/testnet/explorer/Dockerfile -t "$EXPLORER_TAG_ENV" ./deploy/testnet/explorer/

  rm -rf "$BASEDIR"/explorer/web
fi
