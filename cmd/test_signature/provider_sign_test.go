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
