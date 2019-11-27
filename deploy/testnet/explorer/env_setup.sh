#!/bin/sh

echo "----- Setup environment variables -----"
echo "$PWD"
ls -la
sed -i -e \
  's/ENV_GEV_RPCPORT/'"$GETH_RPCPORT"'/g' \
  ./app/app.js