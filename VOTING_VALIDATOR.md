# How to get the list validator from tendermint API?

1. Attach gev console with module tendermint to a running node. Please ensure the node is a validator in network and miner is started.
    ```shell
    $ ./gev --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,tendermint attach ipc:[path_to_geth_ipc]/geth.ipc     
    ```

2. Check validator list.

    ```shell
    > tendermint.getValidators()
    ```

    Returns the list validator of the heighest block of local node. 

    ```shell
    > tendermint.getValidators(blockNumber int)
    ```

    Returns the list validator for a block number. 
