#!/bin/sh

# Build explorer
./testnet/explorer/update_explorer.sh
# Clear network bridge
yes | docker network prune
# Remove evrynet-builder
docker rmi -f img_builder img_bootnode
docker rm -f evrynet-builder gev-bootnode

# Run with Testnet port
#yes | gevRPCPort=22003 docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build
yes | docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build evrynet-builder

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
  echo "Building project successfully!"
  sleep 5

  echo "Copy gev, bootnode from evrynet-builder to ./testnet/builder/bin/"
#  rm ./testnet/builder/bin/gev ./testnet/builder/bin/bootnode
  docker cp evrynet-builder:/evrynet/gev ./testnet/builder/bin/
  echo "Copied gev"
  docker cp evrynet-builder:/evrynet/bootnode ./testnet/builder/bin/
  echo "Copied bootnode"

  # Start bootnode
  cp ./testnet/builder/bin/bootnode ./testnet/bootnode/bin/
  sleep 5
  yes | docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build gev-bootnode

  exit 0
else
  echo "Building project failed!"
  exit 1
fi