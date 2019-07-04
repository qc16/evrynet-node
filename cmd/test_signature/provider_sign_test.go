package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/cyberliem/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestSendWithoutProviderSignature(t *testing.T) {
	const (
		contractAddrStr = "0xe9aABE2Ab51B068682e49126b0C58A725251932f"
		providerPK      = "87668A123F9FF917F43B9F9168BB6A30F897AA30955144C3A74FEA6AC6898BBC"
		senderPK        = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
		senderAddrStr   = "0xa091e44e0B6Adc71ce1f58B81337343597301FF6"

		testBal1     = 1000000 //1e6
		testBal2     = 2000000 //2e6
		testGasLimit = 100000000
	)
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStr)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial("http://localhost:8545")
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(1000000), 1000000, gasPrice, nil)
	// return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data)
	transaction, err = types.SignTx(transaction, signer, spk)
	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NotEqual(t, nil, err)
}

func TestSendWithProviderSignature(t *testing.T) {
	const (
		contractAddrStr = "0xe9aABE2Ab51B068682e49126b0C58A725251932f"
		providerPK      = "87668A123F9FF917F43B9F9168BB6A30F897AA30955144C3A74FEA6AC6898BBC"
		senderPK        = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
		senderAddrStr   = "0xa091e44e0B6Adc71ce1f58B81337343597301FF6"

		testBal1     = 1000000 //1e6
		testBal2     = 2000000 //2e6
		testGasLimit = 1500000
	)
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStr)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(providerPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial("http://localhost:8545")
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(1000000), testGasLimit, gasPrice, nil)
	// return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data)
	transaction, err = types.SignTx(transaction, signer, spk)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)
}

//TestProviderSignTransaction will sign a transaction with both sender's Key and Provider's Key
//Note that the account must be unlocked prior to run this test
//The JSON rpc test can be call as
//curl <rpcserver> -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_providerSignTransaction","params":["<raw_tx>", "<provider_address>"],"id":1}'
func TestProviderSignTransaction(t *testing.T) {
	const (
		//This provider should be the contract's provider or the fixed provider in tx_pool for testing purpose
		providerPK      = "181C205392D2A39453E6CDFB2839C7F0CA77ED2683F1A04B5007AC223DF9DD82"
		senderPK        = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
		senderAddrStr   = "0xa091e44e0B6Adc71ce1f58B81337343597301FF6"
		providerAddrStr = "0x6cb9612fecf88360481e69b9f1e99f833592cdf7"
		//This should be a smart contract with provider address
		scContractStr = "0xe9aabe2ab51b068682e49126b0c58a725251932f"

		testBal1     = 1000000 //1e6
		testBal2     = 2000000 //2e6
		testGasLimit = 1500000
	)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	senderAddr := common.HexToAddress(senderAddrStr)
	providerAddr := common.HexToAddress(providerAddrStr)
	ethClient, err := ethclient.Dial("http://localhost:8545")
	assert.NoError(t, err)
	id, err := ethClient.ChainID(context.Background())
	signer := types.NewEIP155Signer(id)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	tx := types.NewTransaction(nonce, common.HexToAddress(scContractStr), big.NewInt(1000000), testGasLimit, gasPrice, nil)
	txSigned, err := types.SignTx(tx, signer, spk)
	assert.NoError(t, err)
	v, r, s := txSigned.RawSignatureValues()

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
