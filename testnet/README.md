#How to use Testsnet Docker 
##Run Testnet Docker
Everything was setup. You only need to run this script  
`testnet/start_testnet.sh <number_of_miners> <path_to_share_volumes>`    

Ex: `./testnet/start_testnet.sh 4 /Volumes/Work/Projects/KyberNetwork/lab/testnet/15`
- `number_of_miners` is the number of miner you want to run. Maximum miner is 99.
- `path_to_share_volumes` is a path to place where you want to share volumes with docker.  

##Configuring Netstat for your cluster 
Use this script to generate app.json to define which nodes will be shown on Netstat  
`bash /path/to/eth-utils/netstatconf.sh <number_of_clusters> <name_prefix> <rpc_host> <ws_server> <ws_secret> > ./app.json`  

Ex: `bash testnet/utils/netstat_config.sh 5 everynet 172.25.0.102 ws://gev-monitor-frontend:3000 testnet > ./app.json`  
- will output resulting app.json to use in testnet/monitor/backend/app.json
- `number_of_clusters` is the number of nodes in the cluster.
- `name_prefix` is a prefix for the node names as will appear in the listing.
- `rpc_host` is the RPC host
- `ws_server` is the eth-netstats server. Make sure you write the full URL.
- `ws_secret` is the eth-netstats secret.

##Web
- Explorer: http://localhost:8080
- Netstat: http://localhost:3000