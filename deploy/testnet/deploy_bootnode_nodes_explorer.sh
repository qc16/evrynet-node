#!/bin/bash
#deploy/testnet/deploy_bootnode_nodes_explorer.sh <path_to_share_volumes> <rpc_corsdomain> <tag_version_or_develop_branch> <environment> <genesis_path> <deploy_explorer>
# Ex: deploy/testnet/deploy_bootnode_nodes_explorer.sh /Volumes/Work/Kyber/evrynet-node/tests/nodes "*" develop testnet /Volumes/Work/Kyber/evrynet-node/deploy/testnet/nodes/bin/genesis.json n

localVolumes=$1
shift
rpccorsdomain=$1
shift
version=$1
shift
env=$1
shift
genesisPath=$1
shift
deployExplorer=$1
shift

if [[ "$localVolumes" == "" || "$rpccorsdomain" == "" || "$deployExplorer" == "" || "$version" == "" || "$env" == "" || "$genesisPath" == "" ]]
then
  echo 'Missing params'
  exit 1
fi

BOOTNODE_REPOSITORY="kybernetwork/evrynet-bootnode"
BOOTNODE_TAG_ENV="$BOOTNODE_REPOSITORY:$version-$env"

NODE_REPOSITORY="kybernetwork/evrynet-node"
NODE_TAG_ENV="$NODE_REPOSITORY:$version-$env"

echo "--- Pulling images $BOOTNODE_TAG_ENV ..."
yes | sudo docker pull "$BOOTNODE_TAG_ENV"

echo "--- Pulling images $NODE_TAG_ENV ..."
yes | sudo docker pull "$NODE_TAG_ENV"


echo "--- Stopping nodes 1,2,3 ..."
for i in 1 2 3
do
  yes | sudo docker exec -it gev-node-"$i" /bin/sh ./stop_node.sh
  sleep 3
done

# Stop bootnode & nodes container
sudo docker stop gev-bootnode gev-node-1 gev-node-2 gev-node-3
# Remove bootnode nodes container
sudo docker rm -f gev-bootnode gev-node-1 gev-node-2 gev-node-3

# Clear network bridge
yes | sudo docker network prune

# shellcheck disable=SC2028
echo "--- Starting bootnode ..."
yes | sudo docker run --name gev-bootnode -d \
  -p 30300:30300/tcp -p 30300:30300/udp \
  -e NODE_HEX_KEY='9dbcbd49f9f4e1b4949178d7e413142267050377ff99d81c08e371cdea712f09' \
  "$BOOTNODE_TAG_ENV"

bootnodeIP=$(sudo docker inspect -f "{{ .NetworkSettings.IPAddress }}" gev-bootnode)

if [[ "$bootnodeIP" == "" ]]; then
  echo "=> Bootnode ID is empty!"
  exit 1
fi

nodes=(1 2 3)
nodekeys=(
  'ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1'
  'e74f3525fb69f193b51d33f4baf602c4572d81ede57907c61a62eaf9ed95374a'
  '276cd299f350174a6005525a523b59fccd4c536771e4876164adb9f1459b79e4'
)
bootnodeID='aa8d839e6dbe3524e8c0a62aefae7cefa3880f9c7394ddaaa31cc8679fe3a25396e014c5c48814d0fe18d7f03d72a7971fd50b7dd689bd04498d98902dd0d82f'
for i in "${!nodes[@]}"; do
  nodeID="${nodes[i]}"

  # shellcheck disable=SC2181
  if [ "$nodeID" != "3" ];
  then
    echo "--- Starting normal node $nodeID ..."
    yes | sudo docker run --name gev-node-"$nodeID" -d \
      -v "$genesisPath":/node/genesis.json \
      -v "$localVolumes"/node_"$nodeID"/data:/node/data \
      -v "$localVolumes"/node_"$nodeID"/log:/node/log \
      -p 2200"$nodeID":8545/tcp -p 2200"$nodeID":8545/udp \
      -p 606"$nodeID":6060 \
      -p 3030"$nodeID":30303/tcp -p 3030"$nodeID":30303/udp \
      -e NODE_ID="${nodes[i]}" \
      -e NODEKEYHEX="${nodekeys[i]}" \
      -e BOOTNODE_ID=$bootnodeID \
      -e BOOTNODE_IP="$bootnodeIP" \
      -e RPC_CORSDOMAIN="$rpccorsdomain" \
      "$NODE_TAG_ENV"
  else
    echo "--- Starting metrics node $nodeID ..."
    yes | sudo docker run --name gev-node-"$nodeID" -d \
      -v "$genesisPath":/node/genesis.json \
      -v "$localVolumes"/node_"$nodeID"/data:/node/data \
      -v "$localVolumes"/node_"$nodeID"/log:/node/log \
      -p 2200"$nodeID":8545/tcp -p 2200"$nodeID":8545/udp \
      -p 606"$nodeID":6060 \
      -p 3030"$nodeID":30303/tcp -p 3030"$nodeID":30303/udp \
      -e NODE_ID="${nodes[i]}" \
      -e NODEKEYHEX="${nodekeys[i]}" \
      -e BOOTNODE_ID=$bootnodeID \
      -e BOOTNODE_IP="$bootnodeIP" \
      -e RPC_CORSDOMAIN="$rpccorsdomain" \
      -e HAS_METRIC=1 \
      -e METRICS_ENDPOINT="http://52.220.52.16:8086" \
      -e METRICS_USER='test' \
      -e METRICS_PASS='test' \
      "$NODE_TAG_ENV"
  fi
done


if [[ "$deployExplorer" == "y" ]]; then
  EXPLORER_REPOSITORY="kybernetwork/evrynet-explorer"
  EXPLORER_TAG_ENV="$EXPLORER_REPOSITORY:$version-$env"

  echo "--- Pulling images $EXPLORER_TAG_ENV ..."
  yes | sudo docker pull "$EXPLORER_TAG_ENV"

  sudo docker stop gev-explorer
  sudo docker rm -f gev-explorer

  echo "--- Starting explorer on image $EXPLORER_TAG_ENV..."
  yes | sudo docker run --name gev-explorer -d \
      --publish 8080:8080 \
      -e GETH_RPCPORT=22001 \
      "$EXPLORER_TAG_ENV"
fi