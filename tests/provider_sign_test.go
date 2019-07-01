package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

/* These tests are done on a chain with already setup account/ contracts.
To run these test, please deploy your own account/ contract and extract privatekey inorder to get the expected result
*/

func TestSendWithoutProviderSignature(t *testing.T) {
	const (
		//This should be a contract with provider address
		contractAddrStr = "0x6d88d80c9ac4bb26dac4c4bb09a61200f9cb8d75"
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

func TestSendWithProviderSignatureToContractWithoutProviderAddress(t *testing.T) {
	const (
		//This should be a contract without provider address
		contractAddrStr = "0x1a805a2735e069d1d63108f8b5d2408b30a9de4f"
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
	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)
}

func TestSendWithProviderSignatureToContractWithProviderAddress(t *testing.T) {
	const (
		//This should be a contract with provider address
		contractAddrStr = "0x6d88d80c9ac4bb26dac4c4bb09a61200f9cb8d75"
		providerPK      = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
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
	assert.NotEqual(t, nil, err)
}

func TestSendWithCorrectProviderSignatureToContractWithProviderAddress(t *testing.T) {
	const (
		//This should be a contract with provider address
		contractAddrStr = "0x6d88d80c9ac4bb26dac4c4bb09a61200f9cb8d75"
		providerPK      = "E6CFAAD68311D3A873C98750F52B2543F2C3C692A8F11E6B411B390BCD807133"
		senderPK        = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
		senderAddrStr   = "0xa091e44e0B6Adc71ce1f58B81337343597301FF6"
		//providerAddr  = "0x8359d8C955DAef81e171C13659bA3Fb0dDa144b4"

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
