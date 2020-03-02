# How to deploy Testnet

## 1. Preparing Images 
### Building images
First of all, you need to change the `login` & `password` (is your token) at file `deploy/testnet/builder/token` to yours. You can get a token from [Here](https://github.com/settings/tokens).

You use this file `deploy/testnet/build_images_for_bootnode_node_explorer.sh` with appropriate to params to build images for bootnode, node & explorer.  
The params are `<tag_version_or_develop_branch> <evironment> <build_explorer>`:
- `tag_version_or_develop_branch`: is specific tag or branch you wan to build (1.1.2, develop ...)
- `environment`: to use as suffix of image tag (testnet, local, ...)
- `build_explorer`: to build image for explorer or not (y ,n)

Ex:  `deploy/testnet/build_images_for_bootnode_node_explorer.sh develop testnet n`

After building successfully, you can check the images by command `docker images`. The result should be like this:
```
REPOSITORY                      TAG                 IMAGE ID            CREATED             SIZE
kybernetwork/evrynet-node       develop-testnet     96b4226f096f        6 hours ago         85.2MB
kybernetwork/evrynet-bootnode   develop-testnet     1f48c40245c7        6 hours ago         35.1MB
kybernetwork/evrynet-builder    develop-testnet     b5b539596145        6 hours ago         1.08GB
...
```

### Pushing images to Docker Hub
Now you must push images which you just built above to docker hub.   
Make sure you already login to the docker hub by `docker login`.  
Use `docker push <REPOSITORY>:<TAG>` to push the image to Docker Hub.  

Ex:
```
docker push kybernetwork/evrynet-node:develop-testnet
docker push kybernetwork/evrynet-bootnode:develop-testnet
docker push kybernetwork/evrynet-explorer:develop-testnet
```

## 2. Deploy Testnet
Before deploying Testnet, you must stop all nodes gracefully by this file `deploy/testnet/stop_nodes.sh` (avoid crashing data).  

You must use file `deploy/testnet/deploy_bootnode_nodes_explorer.sh` to deploy Testnet with suitable params.   
The params:
- `path_to_share_volumes`: the path to nodes folder. On Testnet is `/home/ubuntu/testnet/nodes`
- `rpc_corsdomain`: to allow URL can call API
- `tag_version_or_develop_branch`: is a specific tag or branch you want to build (1.1.2, develop ...)
- `environment`: to use as a suffix of image tag (testnet, local, ...)
- `genesis_path`: to inject genesis file to image
- `deploy_explorer`: to ask to deploy explorer or not

Ex:
`deploy/testnet/deploy_bootnode_nodes_explorer.sh /home/ubuntu/testnet/nodes "*" develop testnet /home/ubuntu/evrynet-node/deploy/testnet/nodes/bin/genesis.json n`

## 3. Deploy A Dependant Node
You can deploy only one node in a new instance by creating a new bash file as `start_node.sh` with content like this:
```shell script
#!/bin/bash

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
bootnodeIP=$1
shift
nodeID=$1
shift

if [[ "$localVolumes" == "" || "$rpccorsdomain" == "" || "$version" == "" || "$env" == "" || "$genesisPath" == "" || "$bootnodeIP" == "" || "$nodeID" == "" ]]
then
  echo 'Missing params'
  exit 1
fi

NODE_REPOSITORY="kybernetwork/evrynet-node"
NODE_TAG_ENV="$NODE_REPOSITORY:$version-$env"


echo "--- Pulling images $NODE_TAG_ENV ..."
yes | sudo docker pull "$NODE_TAG_ENV"

echo "--- Stopping nodes $nodeID ..."
yes | sudo docker exec -it gev-node-"$nodeID" /bin/sh ./stop_node.sh
sleep 3

sudo docker stop gev-node-"$nodeID"
sudo docker rm -f gev-node-"$nodeID"

# Clear network bridge
yes | sudo docker network prune

echo "--- Starting normal node $nodeID ..."
yes | sudo docker run --name gev-node-"$nodeID" -d \
  -v "$genesisPath":/node/genesis.json \
  -v "$localVolumes"/data:/node/data \
  -v "$localVolumes"/log:/node/log \
  --publish 6060:6060 \
  --publish 8545:8545 \
  --publish 30303:30303 \
  -e NODE_ID="$nodeID" \
  -e NODEKEYHEX="64a56099c703a5fd52c53b046852136c2ab6798130a68c440885c95cd0c2d069" \
  -e BOOTNODE_ID="aa8d839e6dbe3524e8c0a62aefae7cefa3880f9c7394ddaaa31cc8679fe3a25396e014c5c48814d0fe18d7f03d72a7971fd50b7dd689bd04498d98902dd0d82f" \
  -e BOOTNODE_IP="$bootnodeIP" \
  -e RPC_CORSDOMAIN="$rpccorsdomain" \
  "$NODE_TAG_ENV"
```
The params:
- `path_to_share_volumes`: path to nodes folder. 
- `rpc_corsdomain`: to allow URL can call API
- `tag_version_or_develop_branch`: is specific tag or branch you wan to build (1.1.2, develop ...)
- `environment`: to use as suffix of image tag (testnet, local, ...)
- `genesis_path`: to inject genesis file to image
- `bootnode_ip`: to register node to bootnode
- `node_id`: is ID of node

Ex:
`./start_node.sh /home/ubuntu/node "*" develop testnet /home/ubuntu/node/genesis.json 52.220.52.16 4`