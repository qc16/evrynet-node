#!/bin/sh
./testnet/explorer/update_explorer.sh
# Clear network bridge
yes | docker network prune
# Run with Testnet port
yes | gevRPCPort=22003 docker-compose -f ./testnet/docker-compose.yml up -d --force-recreate --build