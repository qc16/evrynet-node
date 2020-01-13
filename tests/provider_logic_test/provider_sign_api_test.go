package test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Evrynetlabs/evrynet-client/common"
	"github.com/Evrynetlabs/evrynet-client/core/types"
	"github.com/Evrynetlabs/evrynet-client/crypto"
	"github.com/Evrynetlabs/evrynet-client/ethclient"
	"github.com/stretchr/testify/assert"
)

//TestProviderSignTransaction will sign a transaction with both sender's Key and Providers's Key
//Note that the account must be unlocked prior to run this test
//The JSON rpc test can be call as
//curl <rpcserver> -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_providerSignTransaction","params":["<raw_tx>", "<provider_address>"],"id":1}'
func TestProviderSignTransaction(t *testing.T) {
	contractAddr := common.HexToAddress(contractAddrStrWithProvider)

	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	senderAddr := common.HexToAddress(senderAddrStr)
	providerAddr := common.HexToAddress(providerAddrStr)
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	id, err := ethClient.ChainID(context.Background())
	signer := types.NewEIP155Signer(id)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(1000000), testGasLimit, gasPrice, nil)
	txSigned, err := types.SignTx(tx, signer, spk)
	assert.NoError(t, err)
	v, r, s := txSigned.RawSignatureValues()

	ppk, err := crypto.HexToECDSA(providerPK)
	// Check Tx for existion
	_, err = types.ProviderSignTx(txSigned, signer, ppk)
	assert.NoError(t, err)

	// Get Tx via RPC
	pTxSigned, err := ethClient.ProviderSignTx(context.Background(), txSigned, &providerAddr)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, pTxSigned)

	v2, r2, s2 := pTxSigned.RawSignatureValues()
	pv, pr, ps := pTxSigned.RawProviderSignatureValues()
	assert.Equal(t, v, v2)
	assert.Equal(t, r, r2)
	assert.Equal(t, s, s2)
	assert.NotEqual(t, nil, pv)
	assert.NotEqual(t, nil, pr)
	assert.NotEqual(t, nil, ps)

	//The transaction should be able to process
	err = ethClient.SendTransaction(context.Background(), pTxSigned)
	assert.NoError(t, err)
}
