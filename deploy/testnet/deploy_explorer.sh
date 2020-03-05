#!/bin/bash
#deploy/testnet/deploy_explorer.sh <tag_version_or_develop_branch> <environment>
# Ex: deploy/testnet/deploy_explorer.sh develop testnet

version=$1
shift
env=$1
shift

if [[ "$version" == "" || "$env" == "" ]]
then
  echo 'Missing params'
  exit 1
fi

EXPLORER_REPOSITORY="kybernetwork/evrynet-explorer"
EXPLORER_TAG_ENV="$EXPLORER_REPOSITORY:$version-$env"

echo "--- Pulling images $EXPLORER_TAG_ENV ..."
yes | sudo docker pull "$EXPLORER_TAG_ENV"

sudo docker rm -f gev-explorer

echo "--- Starting explorer on image $EXPLORER_TAG_ENV..."
yes | sudo docker run --name gev-explorer -d \
  --publish 8080:8080 \
  -e GETH_RPCPORT=22001 \
  "$EXPLORER_TAG_ENV"
