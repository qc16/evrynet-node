#!/bin/sh

# Build explorer
./testnet/explorer/update_explorer.sh
# Clear network bridge
yes | docker network prune
# Remove evrynet-builder
docker rmi -f img_builder img_bootnode img_node img_miner_1 img_miner_2 img_monitor_backend img_monitor_frontend img_miners
docker rm -f evrynet-builder gev-bootnode gev-node gev-miner-1 gev-miner-2 gev-monitor-backend gev-monitor-frontend gev-miners

# Run with Testnet port
#yes | gevRPCPort=22003 docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build
yes | docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build evrynet-builder

# shellcheck disable=SC2181
if [ $? -eq 0 ]
then
  echo "Building project successfully!"
  sleep 5

  echo "Copy gev, bootnode from evrynet-builder to ./testnet/builder/bin/"
  docker cp evrynet-builder:/evrynet/gev ./testnet/builder/bin/
  docker cp evrynet-builder:/evrynet/bootnode ./testnet/builder/bin/
  docker cp evrynet-builder:/evrynet/makegenesis ./testnet/builder/bin/

  # Start bootnode
  echo "Start bootnode"
  cp ./testnet/builder/bin/bootnode ./testnet/bootnode/bin/
  yes | docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build gev-bootnode

  # Start nodes
  echo "Start miners"
  cp ./testnet/builder/bin/gev ./testnet/miner/bin/
  cp ./testnet/builder/bin/makegenesis ./testnet/miner/bin/
  yes | docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build gev-miners gev-monitor-backend gev-monitor-frontend

  exit 0
else
  echo "Building project failed!"
  exit 1
fi