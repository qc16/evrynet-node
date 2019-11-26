package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/common/hexutil"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/crypto"
	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/stretchr/testify/assert"
)

//TestProviderSignTransaction will sign a transaction with both sender's Key and Providers's Key
//Note that the account must be unlocked prior to run this test
//The JSON rpc test can be call as
//curl <rpcserver> -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_providerSignTransaction","params":["<raw_tx>", "<provider_address>"],"id":1}'
func TestProviderSignTransaction(t *testing.T) {
	contractAddr := prepareNewContract(false)
	assert.NotNil(t, contractAddr)

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

	tx := types.NewTransaction(nonce, *contractAddr, big.NewInt(1000000), testGasLimit, gasPrice, nil)
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
	if err != nil {
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
}

func prepareNewContract(hasProvider bool) *common.Address {
	var (
		tx           *types.Transaction
		providerAddr = common.HexToAddress(providerAddrStr)
		ownerAddr    = common.HexToAddress(senderAddrStr)
		sender       = common.HexToAddress(senderAddrStr)
	)

	spk, err := crypto.HexToECDSA(senderPK)
	if err != nil {
		return nil
	}
	payLoadBytes, err := hexutil.Decode(payload)
	if err != nil {
		return nil
	}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	if err != nil {
		return nil
	}
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return nil
	}

	if hasProvider {
		var option types.CreateAccountOption
		option.ProviderAddress = &providerAddr
		option.OwnerAddress = &ownerAddr
		tx = types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes, option)
	} else {
		tx = types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes)
	}

	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	if err != nil {
		return nil
	}

	err = ethClient.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil
	}
	if err == nil {
		var (
			maxTrie = 10
			trie    = 1
		)

		for {
			if trie > maxTrie {
				break
			}
			var receipt *types.Receipt
			receipt, err = ethClient.TransactionReceipt(context.Background(), tx.Hash())
			if err == nil && receipt.Status == uint64(1) {
				return &receipt.ContractAddress
			}

			time.Sleep(1 * time.Second)
			trie = trie + 1
		}
	}
	return nil
}
