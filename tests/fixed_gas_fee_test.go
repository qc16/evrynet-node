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

// TestSendNormalTxWithFixedFee
func TestSendNormalTxWithFixedFee(t *testing.T) {
	const (
		normalAddress = "0xBBD9e63B95308358AAfb20d6606701A4b6429f5e"
		senderPK      = "112CD7FA616EF6499DA9FA0A227AC73B4B109CC3F7F94C2BEFB3346CCB18CD08"
		senderAddrStr = "0xa091e44e0B6Adc71ce1f58B81337343597301FF6"

		testBal1     = 1000000 //1e6
		testBal2     = 2000000 //2e6
		testGasLimit = 100000000
	)

	var (
		senderAddr    = common.HexToAddress(senderAddrStr)
		normalAddr    = common.HexToAddress(normalAddress)
		fixedGasPrice = big.NewInt(1000000000)
	)

	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	signer := types.HomesteadSigner{}
	ethClient, err := ethclient.Dial("http://localhost:9015")
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), senderAddr)
	assert.NoError(t, err)

	//SuggestGasPrice will return fixedGasPrice
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, gasPrice, fixedGasPrice)

	//this transaction should be reject since its gas price is not the fixed gas price
	transaction := types.NewTransaction(nonce, normalAddr, big.NewInt(1000000), 1000000, big.NewInt(2000000), nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), transaction))

	//only transaction with gixedGasPrice/nil gas price is success
	transaction = types.NewTransaction(nonce, normalAddr, big.NewInt(1000000), 1000000, fixedGasPrice, nil)
	transaction, err = types.SignTx(transaction, signer, spk)
	assert.NoError(t, ethClient.SendTransaction(context.Background(), transaction))
}
