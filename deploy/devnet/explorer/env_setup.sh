#!/bin/bash

echo "----- Setup environment variables -----"
ls -la

# shellcheck disable=SC2039
if [[ ! $GETH_RPCPORT ]]; then
  GETH_RPCPORT="22001"
fi

sed -i -e \
  's/ENV_GEV_RPCPORT/'"$GETH_RPCPORT"'/g' \
  ./app/app.js