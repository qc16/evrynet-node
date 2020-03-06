#!/bin/bash

version=
until [[ $version ]]; do read -rp "- Tag Version/Branch Name you want to deploy: " version; done
env=
until [[ $env ]]; do read -rp "- Environment of Image: " env; done

EXPLORER_REPOSITORY="kybernetwork/evrynet-explorer"
EXPLORER_TAG_ENV="$EXPLORER_REPOSITORY:$version-$env"

# Check status of explorer image
if [[ "$(sudo docker images -q "$EXPLORER_TAG_ENV" 2>/dev/null)" == "" ]]; then
  pullExplorerImage=
  until [[ $pullExplorerImage ]]; do read -rp "** Image $EXPLORER_TAG_ENV doesn't exist on your local! Do you want to pull this image? " pullExplorerImage; done

  if [[ "$pullExplorerImage" == "y" ]]; then
    yes | sudo docker pull "$EXPLORER_TAG_ENV"

    # Check status of explorer image
    if [[ "$(sudo docker images -q "$EXPLORER_TAG_ENV" 2>/dev/null)" == "" ]]; then
      echo "=> Can not pull Image $EXPLORER_TAG_ENV . Make sure this image has existed on your hub!"
      exit 1
    fi
  else
    echo "=> You must build explorer image on you local. Read README.md to know the way to build!"
    exit 1
  fi
fi

yes | sudo docker rm -f gev-explorer

echo "--- Starting explorer on image $EXPLORER_TAG_ENV..."
yes | sudo docker run --name gev-explorer -d \
  --publish 8080:8080 \
  -e GETH_RPCPORT=22001 \
  "$EXPLORER_TAG_ENV"
