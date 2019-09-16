package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/stretchr/testify/assert"
)

// TestGetExtraData run unit test for the ExtraDataByBlockNumber and ExtraDataByBlockHash
// expected get address of proposer and commit signers
// you can run via curl with
// curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getExtraDataByBlockNumber","params":["0x5"],"id":1}' http://localhost:8545
// curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getExtraDataByBlockHash","params":["0x4cddd578050781138591684e16674e7056ce89fb2706cf72b9f946e7e168d73b"],"id":1}' http://localhost:8545
func TestGetExtraData(t *testing.T) {
	ethClient, err := ethclient.Dial("http://localhost:8545")
	assert.NoError(t, err)
	fakeBlockNumber := big.NewInt(5)

	trie := 10
	for {
		if trie <= 0 {
			t.Error("You have to run a node and start miner before")
			break
		}
		block, err := ethClient.BlockByNumber(context.Background(), fakeBlockNumber)
		if err == nil {
			extra, err := ethClient.ExtraDataByBlockNumber(context.Background(), block.Number())
			assert.NoError(t, err)
			assert.NotNil(t, extra)
			assert.NotNil(t, extra.BlockProposer)
			assert.NotEqual(t, common.Address{}, extra.BlockProposer)

			extra, err = ethClient.ExtraDataByBlockHash(context.Background(), block.Hash())
			assert.NoError(t, err)
			assert.NotNil(t, extra)
			assert.NotNil(t, extra.BlockProposer)
			assert.NotEqual(t, common.Address{}, extra.BlockProposer)
			break
		}
		time.Sleep(2 * time.Second)
		trie--
	}

}
