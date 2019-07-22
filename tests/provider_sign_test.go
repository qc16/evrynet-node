package tests

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

/* These tests are done on a chain with already setup account/ contracts.
To run these test, please deploy your own account/ contract and extract privatekey inorder to get the expected result
Adjust these params to match deployment on local machine:
*/

const (
	normalAddress                  = "0x11c93c29591ba613852ac2c9278faec2d7e7ea59"
	senderPK                       = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
	senderAddrStr                  = "0xa091e44e0B6Adc71ce1f58B81337343597301FF6"
	contractAddrStrWithoutProvider = "0xe9aABE2Ab51B068682e49126b0C58A725251932f"
	contractAddrStrWithProvider    = "0xA014d882aA3bd232c96c2AacbCCEcb334eE48B5b"
	providerPK                     = "E6CFAAD68311D3A873C98750F52B2543F2C3C692A8F11E6B411B390BCD807133"
	invadlidProviderPK             = "5564a4ddd059ba6352aae637812ea6be7d818f92b5aff3564429478fcdfe4e8a"
	providerAddrStr                = "0x8359d8C955DAef81e171C13659bA3Fb0dDa144b4"

	// make sure this contract allows to receive ETH
	payload = "0x6080604052348015600f57600080fd5b5060a68061001e6000396000f3fe60806040526004361060265760003560e01c80633fa4f2451460285780635524107714604c575b005b348015603357600080fd5b50603a6066565b60408051918252519081900360200190f35b602660048036036020811015606057600080fd5b5035606c565b60005481565b60005556fea265627a7a72305820bc1f0b3fba5fd519e0e123a58363a4aa98115b675e2ca69adbd2166e9d06872364736f6c634300050a0032"

	providerWithoutGasAddr     = "0x6D4ke 9e3Ba3f77a19e9dF7EceD8AA7154Fe372ea27"
	providerWithoutGasPK       = "34b377a903b4a01c55062d978160084992271c4f89797caa18fd4e1d61123fbb"
	contractProviderWithoutGas = "0x4d988Aebd3Ee0e30426c2C0b003A515991De8657"

	senderWithoutGasPK      = "CD79C18795A866C4A7FA8D3A88494F618AB0E69B1493382D638A6483538EEA97"
	senderWithoutGasAddrStr = "0xBBD9e63B95308358AAfb20d6606701A4b6429f5e"

	testGasLimit   = 1000000
	testGasPrice   = 1000000000
	testAmountSend = 1000000
	ethRPCEndpoint = "http://localhost:8545"
)

/*
	Test Send ETH to a normal address
		- No provider signature is required
*/
func TestSendToNormalAddress(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	normalAddr := common.HexToAddress(normalAddress)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, normalAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)
	// Check gas payer, should be sender's address
	for {
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), transaction.Hash())
		if err == nil {
			assert.Equal(t, receipt.GasPayer, senderAddr)
			assert.Equal(t, receipt.Status, uint64(1))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

/*
	Test send to a normal address with provider's signature
		- Expect to get error with redundant provider's signature
*/
func TestSendToNormalAddressWithProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	normalAddr := common.HexToAddress(normalAddress)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(providerPK)
	assert.NoError(t, err)
	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, normalAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))
}

/*
	Test Send ETH to a Smart Contract without provider's signature
		- Provider's signature is not required
*/
func TestSendToNonEnterpriseSmartContractWithoutProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithoutProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, nil)
	// return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data)
	transaction, err = types.SignTx(transaction, signer, spk)
	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)

	// Check gasPayer, should be sender's address
	for {
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), transaction.Hash())
		if err == nil {
			assert.Equal(t, receipt.GasPayer, senderAddr)
			assert.Equal(t, receipt.Status, uint64(1))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

/*
	Test send ETH to a Non-enterprise Smart Contract with provider's signature
		- Expect to get error as provider's signature is redundant
*/
func TestSendToNonEnterpriseSmartContractWithProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithoutProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	ppk, err := crypto.HexToECDSA(providerPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))
}

/*
	Test interact with Non-Enterprise Smart Contract
		- Update value inside Smart Contract and expect to get no error (skip provider check)
	Note: Please change data to your own function data
*/
func TestInteractWithNonEnterpriseSmartContractWithoutProviderSignature(t *testing.T) {
	//This should be a contract with provider address
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithoutProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	// data to interact with a function of this contract
	dataBytes := []byte("0x552410770000000000000000000000000000000000000000000000000000000000000002")
	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)
}

/*
	Test Send ETH to an Enterprise Smart Contract with invalid provider's signature
*/
func TestSendToEnterPriseSmartContractWithInvalidProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(invadlidProviderPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))
}

/*
	Test Send ETH to an enterprise Smart Contract with valid provider's signature
*/
func TestSendToEnterPriseSmartContractWithValidProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(providerPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)

	for {
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), transaction.Hash())
		if err == nil {
			assert.Equal(t, receipt.GasPayer, common.HexToAddress(providerAddrStr))
			assert.Equal(t, receipt.Status, uint64(1))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

/*
	Test interact with Enterprise Smart Contract
		- Update value inside Smart Contract and expect to get error with invalid provider signature
	Note: Please change data to your own function data
*/
func TestInteractToEnterpriseSmartContractWithInvalidProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(invadlidProviderPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	// data to interact with a function of this contract
	dataBytes := []byte("0x552410770000000000000000000000000000000000000000000000000000000000000004")
	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))
}

/*
	Test interact with Enterprise Smart Contract
		- Update value inside Smart Contract and expect to successfully update data with valid provider signature
	Note: Please change data to your own function data
*/
func TestInteractToEnterpriseSmartContractWithValidProviderSignature(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractAddrStrWithProvider)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(providerPK)
	assert.NoError(t, err)

	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	// data to interact with a function of this contract
	dataBytes := []byte("0x552410770000000000000000000000000000000000000000000000000000000000000004")
	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(testAmountSend), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)

	for {
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), transaction.Hash())
		if err == nil {
			assert.Equal(t, receipt.GasPayer, common.HexToAddress(providerAddrStr))
			assert.Equal(t, receipt.Status, uint64(1))
			break
		}
		time.Sleep(1 * time.Second)
	}
}
