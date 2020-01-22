package test

import (
	"context"
	"testing"
	"time"

	"github.com/Evrynetlabs/evrynet-node/rpc"
	"github.com/stretchr/testify/assert"
)

const (
	RPCEndpoint = "http://0.0.0.0:22001"
)

func TestAllApis(t *testing.T) {

	rpcClient, err := rpc.Dial(RPCEndpoint)
	assert.NoError(t, err)

	var (
		prefix  = "eth_"
		methods = []string{
			"protocolVersion", "syncing", "coinbase", "mining", "hashrate", "gasPrice", "accounts",
			"blockNumber", "getStorageAt", "getTransactionCount", "getBlockTransactionCountByHash",
			"getBlockTransactionCountByNumber", "getUncleCountByBlockHash", "getUncleCountByBlockNumber",
			"getCode", "sign", "providerSignTransaction", "sendTransaction", "sendRawTransaction",
			"estimateGas", "getBlockByHash", "getBlockSignerByHash", "getBlockByNumber",
			"getBlockSignerByNumber", "getTransactionByHash", "getTransactionByBlockHashAndIndex",
			"getTransactionByBlockNumberAndIndex", "getTransactionReceipt", "getUncleByBlockHashAndIndex",
			"getCompilers", "compileSolidity", "newFilter",
			"newBlockFilter", "newPendingTransactionFilter", "uninstallFilter", "getFilterChanges",
			"getFilterLogs", "getLogs",
		}
	)

	for _, element := range methods {
		if element == "" {
			continue
		}
		method := prefix + element
		err = rpcClient.CallContext(context.Background(), nil, method)
		if err != nil {
			assert.NotContains(t, err.Error(), "the method "+method+" does not exist/is not available")
		}

		time.Sleep(200 * time.Millisecond)
	}

}
