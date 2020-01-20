# How to use Testsnet Docker  

## 1. Using Existing Config
### Before building Docker image
- To allow docker can run `go get` for private repository when building evrynet-client. You must change the `login` & `password` (is your token) in `deploy/testnet/builder/token` to yours.  
- [Here](https://github.com/settings/tokens) setup your token.  

### Run Testnet Docker
Everything was setup. After updating Github token, running this script to deploy:  
`./deploy/testnet/deploy.sh <path_to_share_volumes> <rpc_corsdomain> <wanna_to_deploy_explorer>`    

Ex: `./deploy/testnet/deploy.sh /Volumes/Work/Projects/KyberNetwork/evrynet-client/deploy/testnet/nodes/data localhost y`
- `path_to_share_volumes` is a path to folder where you want to share volumes with docker. The folder must include nodekey and keystore in each node. Ex: `deploy/testnet/nodes/data` 
- `rpc_corsdomain` is a domain which was allowed to call RPC API to node  
- `wanna_to_deploy_explorer` if you wanna deploy explorer, input is `y`

### Nodes Information
Everything about 3 nodes I put at `deploy/testnet/nodes/data`.  
You can clear all data by running this file `deploy/testnet/nodes/data/clear_data.sh`

### Webs
- Explorer: http://localhost:8080

### **NOTICE!
- If you want to stop nodes, DON'T USE `docker stop ...`. It can make a crash to DB of nodes => can not run for next time!
- To stop nodes gracefully, USE this file `deploy/testnet/stop_dockers.sh`. It will interact with node in docker to stop gracefully.

---

## 2. Setup From Zero
1. Build and export to `PATH`
    ```shell script
    $ go build ./cmd/gev
    $ go build ./cmd/bootnode
    $ go build ./cmd/puppeth
    $ export PATH=$(pwd):$PATH
    ```
2. Move to folder where you want to store nodes data. Creating a working directory for 3 validator nodes  
    ```shell script
    $ mkdir node_1 node_2 node_3
    ```  

3. In each node’s working directory, create a log & data directory called `data`, and inside `data` create the `geth` directory   
    ```shell script
    $ mkdir -p node_1/log
    $ mkdir -p node_2/log
    $ mkdir -p node_3/log
    $ mkdir -p node_1/data/geth
    $ mkdir -p node_2/data/geth
    $ mkdir -p node_3/data/geth
    ```

4. Generate node key and copy it into folder `node_1`, `node_2` `node_3`
    ```shell script
    $ bootnode --genkey=nodekey1
    $ bootnode --genkey=nodekey2
    $ bootnode --genkey=nodekey3
    ```
   
5. Now we will generate initial accounts for any of the nodes in the required node’s working directory. The resulting public account address printed in the terminal should be recorded. Repeat as many times as necessary. A set of funded accounts may be required depending what you are trying to accomplish  
    ```shell script
    $ gev --datadir node_1/data account new
    INFO [06-11|16:05:53.672] Maximum peer count                       ETH=25 LES=0 total=25
    Your new account is locked with a password. Please give a password. Do not forget this password.
    Passphrase: 
    Repeat passphrase: 
    Public address of the key:   0x106674Ec8dc5eAA1fB69A3adD61Da9ADdC78cC34
    Path of the secret key file: nodes/node_1/data/keystore/UTC--2019-11-27T10-54-06.993216000Z--106674ec8dc5eaa1fb69a3add61da9addc78cc34

    $ gev --datadir node_2/data account new
    ...
   
    $ gev --datadir node_3/data account new
    ... 
    ```
  
 6. Now we will get address for 3 nodes. Using the content of nodekey1,2,3 files (Ex: `node_1/data/geth/nodekey`) as Private Key to get Address of each nodes at [myetherwallet](myetherwallet.com) 
 7. Last step, we need to update `deploy/testnet/nodes/bin/genesis.json` by new address of validators.  
 Run `puppeth` and full fill data of your chain. Using 3 Addresses we just get from 3 Private Keys above.   
    ```shell script
    Please specify a network name to administer (no spaces, hyphens or capital letters please)
    > testnet
    
    Sweet, you can set this via --network=testnet next time!
    
    INFO [11-27|18:06:28.498] Administering Evrynet network           name=testnet
    INFO [11-27|18:06:28.501] No remote machines to gather stats from
    
    What would you like to do? (default = stats)
     1. Show network stats
     2. Configure new genesis
     3. Track new remote server
     4. Deploy network components
    > 2
    
    What would you like to do? (default = create)
     1. Create new genesis from scratch
     2. Import already existing genesis
    > 1
    
    Which consensus engine to use? (default = clique)
     1. Ethash - proof-of-work
     2. Clique - proof-of-authority
     3. Tendermint - practical-byzantine-fault-tolerance
    > 3
    How many block (Epoch) after which to checkpoint and reset the pending votes (default 30000)
    > 30000
    What is poclicy to select proposer (default 0 - roundrobin)
    > 0
    
    Which accounts are validators? (mandatory at least one)
    > 0xaddress_node_1 (Put 3 Address of nodes here)
    > 0xaddress_node_2
    > 0xaddress_node_2
    
    Which accounts should be pre-funded? (advisable at least one)
    > 0xaddress_node_1 (Put 3 Address of nodes here)
        > 0xaddress_node_2
        > 0xaddress_node_2
    
    Should the precompile-addresses (0x1 .. 0xff) be pre-funded with 1 wei? (advisable yes)
    > yes
    
    Specify your chain/network ID if you want an explicit one (default = random)
    > 15
    
    INFO [11-27|18:08:51.838] Configured new genesis block
    
    What would you like to do? (default = stats)
     1. Show network stats
     2. Manage existing genesis
     3. Track new remote server
     4. Deploy network components
    > 2
    
     1. Modify existing fork rules
     2. Export genesis configurations
     3. Remove genesis configuration
    > 2
    
    Which folder to save the genesis specs into? (default = current)
      Will create testnet.json, testnet-aleth.json, testnet-harmony.json, testnet-parity.json
    >
    INFO [11-27|18:09:06.979] Saved native genesis chain spec          path=testnet.json
    ERROR[11-27|18:09:06.980] Failed to create Aleth chain spec        err="unsupported consensus engine"
    ERROR[11-27|18:09:06.980] Failed to create Parity chain spec       err="unsupported consensus engine"
    INFO [11-27|18:09:06.981] Saved genesis chain spec                 client=harmony path=testnet-harmony.json
    ```
8. Replace all content of `deploy/testnet/nodes/bin/genesis.json` with `testnet.json` (new genesis file just created)
9. Repeat steps at session **1. Using Existing Config**