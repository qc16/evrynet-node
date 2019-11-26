# How to a validator node can vote adding a new validator, or vote removing a validator?

1. Attach gev console with module tendermint to a running node. Please ensure the node is a validator in network and miner is started.
    ```shell
    $ ./gev --rpcapi admin,db,evr,debug,miner,net,shh,txpool,personal,web3,tendermint attach ipc:[path_to_geth_ipc]/geth.ipc     
    ```

2. Run command voting in gev console 

    For vote adding a new validator

    ```shell
    > tendermint.proposeValidator("0x3Cf628d49Ae46b49b210F0521Fbd9F82B461A9E1", true)
    ```

    For vote removing a validator
    ```shell
    > tendermint.proposeValidator("0x3Cf628d49Ae46b49b210F0521Fbd9F82B461A9E1", false)
    ```

3. Check whether having a pending validator, prepared for vote:

    ```shell
    > tendermint.getPendingProposedValidator()
    ```

    If there is no pending validator, the reponse must be
    ```
    {
        validator: "0x0000000000000000000000000000000000000000",
        vote: false
    }
    ```

    Notice: If modified validator in previsous step is added in the block when the node is proposer, pending validator will be removed. The response is similar like above.

4. Remove pending validator.

    If a modified validator have been not added in the block yet, owner of node can remove it by the command:
    ```shell
    > tendermint.clearPendingProposedValidator()
    ```

5. Check validator list.

    ```shell
    > tendermint.getValidators()
    ```

    Returns the list validator of the heighest block of local node. 

    ```shell
    > tendermint.getValidators(blockNumber int)
    ```

    Returns the list validator for a block number. 
