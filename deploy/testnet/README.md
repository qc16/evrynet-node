#How to use Testsnet Docker 
##Run Testnet Docker
Everything was setup. You only need to run this script  
`./deploy/testnet/deploy.sh <path_to_share_volumes> <rpc_corsdomain> <wanna_to_deploy_explorer>`    

Ex: `./deploy/testnet/deploy.sh /Volumes/Work/Projects/KyberNetwork/evrynet-client/deploy/testnet/nodes/data http://localhost:8080 y`
- `path_to_share_volumes` is a path to folder where you want to share volumes with docker. The folder must include nodekey and keystore in each node. Ex: `deploy/testnet/nodes/data` 
- `rpc_corsdomain` is a domain which was allowed to call RPC API to node  
- `wanna_to_deploy_explorer` if you wanna deploy explorer, input is `y`

## Nodes Information
Everything about 3 nodes I put at `deploy/testnet/nodes/data`.  
You can clear all data by running this file `deploy/testnet/nodes/data/clear_data.sh`

##Webs
- Explorer: http://localhost:8080

##NOTICE
- If you want to stop nodes, DON'T USE `docker stop ...`. It can make a crash to DB of nodes => can not run for next time!
- To stop nodes gracefully, USE this file `deploy/testnet/stop_dockers.sh`. It will interact with node in docker to stop gracefully.