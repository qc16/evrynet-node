package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/evrynet-official/evrynet-client/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/evrynet-official/evrynet-client/common"
	"github.com/evrynet-official/evrynet-client/core/types"
	"github.com/evrynet-official/evrynet-client/ethclient"
)

/* These tests are done on a chain with already setup account/ contracts.
To run these test, please deploy your own account/ contract and extract privatekey inorder to get the expected result
*/

// TestInteractToEnterpriseSmartContractWithValidProviderSignatureFromAccountWithoutGas
// Will attempt to reproduce logic of provider paying gas fee.
// It should be send from address without any native token
// The balance of provider should be check prior and after the transaction is mined to
// assure the correctness of the program.
func TestInteractToEnterpriseSmartContractWithValidProviderSignatureFromAccountWithoutGas(t *testing.T) {
	var (
		senderAddr   = common.HexToAddress(senderWithoutGasAddrStr)
		contractAddr = prepareNewContract(true)
	)

	spk, err := crypto.HexToECDSA(senderWithoutGasPK)
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
	dataBytes := []byte("0x3fb5c1cb0000000000000000000000000000000000000000000000000000000000000002")
	transaction := types.NewTransaction(nonce, *contractAddr, big.NewInt(0), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)

	var (
		maxTrie = 10
		trie    = 1
	)

	for {
		if trie > maxTrie {
			break
		}
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), transaction.Hash())
		if err == nil {
			assert.Equal(t, common.HexToAddress(providerAddrStr), receipt.GasPayer)
			assert.Equal(t, uint64(1), receipt.Status)
			break
		}
		time.Sleep(1 * time.Second)
		trie = trie + 1
	}
}

// Interact with a payable function and sending some native token along with transaction
// Please make sure the sender does not have any funds
// expected to get revert as sender's balance is not enough for transaction amount
func TestInteractWithAmountToEnterpriseSmartContractWithValidProviderSignatureFromAccountWithoutGas(t *testing.T) {
	senderAddr := common.HexToAddress(senderWithoutGasAddrStr)
	contractAddr := prepareNewContract(false)
	assert.NotNil(t, contractAddr)

	spk, err := crypto.HexToECDSA(senderWithoutGasPK)
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
	dataBytes := []byte("0x3fb5c1cb0000000000000000000000000000000000000000000000000000000000000002")
	transaction := types.NewTransaction(nonce, *contractAddr, big.NewInt(1000000), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))
}

// Interact with a payable function and sending some native token along with transaction
// Please make sure sender has enough balance to cover transaction amount
// expected to get passed as sender's balance is enough for transaction amount
func TestInteractWithAmountToEnterpriseSmartContractWithValidProviderSignatureFromAccountWithEnoughBalance(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := prepareNewContract(true)
	assert.NotNil(t, contractAddr)

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
	dataBytes := []byte("0x3fb5c1cb0000000000000000000000000000000000000000000000000000000000000002")
	transaction := types.NewTransaction(nonce, *contractAddr, big.NewInt(1000000), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), transaction)
	assert.NoError(t, err)

	var (
		maxTrie = 10
		trie    = 1
	)

	for {
		if trie > maxTrie {
			break
		}
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), transaction.Hash())
		if err == nil {
			assert.Equal(t, common.HexToAddress(providerAddrStr), receipt.GasPayer)
			assert.Equal(t, uint64(1), receipt.Status)
			break
		}
		time.Sleep(1 * time.Second)
		trie = trie + 1
	}
}

// Interact with enterprise contract where provider has zero gas
// Please make sure sender has balance and provider has zero balance
// Expected to get failure as provider's balance is not enough for transaction fee
// Please check error message
func TestInteractEnterpriseSmartContractWithValidProviderSignatureWithoutGas(t *testing.T) {
	senderAddr := common.HexToAddress(senderAddrStr)
	contractAddr := common.HexToAddress(contractProviderWithoutGas)
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)

	ppk, err := crypto.HexToECDSA(providerWithoutGasPK)
	assert.NoError(t, err)
	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)

	// data to interact with a function of this contract
	dataBytes := []byte("0x3fb5c1cb0000000000000000000000000000000000000000000000000000000000000002")
	transaction := types.NewTransaction(nonce, contractAddr, big.NewInt(0), testGasLimit, gasPrice, dataBytes)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	transaction, err = types.ProviderSignTx(transaction, signer, ppk)
	assert.NoError(t, err)

	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))
}
