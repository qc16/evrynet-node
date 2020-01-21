#!/bin/bash
#deploy/testnet/build_image_by_version.sh <version>

version=$1
shift

if [ "$version" == "" ]
then
  echo 'Missing params'
  exit 1
fi

BASEDIR=$(dirname "$0")

echo "--- Cloning evrynet-node from tag $version ..."
git fetch --all --tags --prune
git clone -b "$version" git@github.com:Evrynetlabs/evrynet-node.git "$BASEDIR"/builder/project

echo "--- Building builder container for version $version"
yes | sudo version="kybernetwork/evrynet-builder:$version" docker-compose -f "$BASEDIR"/docker-compose.yml up -d --force-recreate --build gev-builder

rm -rf "$BASEDIR"/builder/project

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
  echo "=> Building version $version successfully!"
  exit 0
else
  echo "Building version $version failed!"
  exit 1
fi