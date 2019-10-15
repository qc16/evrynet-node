#!/bin/sh

ancient_file=./testnet/explorer/explorer/app
# shellcheck disable=SC2039
if [[ -d "$ancient_file" ]]; then
  echo "Update explorer"
	cd ./testnet/explorer/explorer && git checkout -f && git pull
else
  echo "Clone all explorer"
	git clone git@github.com:evrynet-official/explorer.git ./testnet/explorer/explorer
fi
