#!/bin/sh
# bash <path_to_this_file> <number_of_miners> <path_to_share_volumes>

pathShareVolumes=$1
shift

# Clear network bridge
yes | docker network prune
# Remove evrynet-builder
docker rmi -f img_bootnode img_node_1 img_node_2 img_node_3
docker rm -f gev-bootnode gev-node-1 gev-node-2 gev-node-3

yes | shareVolumes=$pathShareVolumes docker-compose -f ./docker-compose-testnet.yml up -d --force-recreate --build gev-bootnode gev-node-1 gev-node-2 gev-node-3
exit