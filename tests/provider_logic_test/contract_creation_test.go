package test

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Evrynetlabs/evrynet-client/common"
	"github.com/Evrynetlabs/evrynet-client/common/hexutil"
	"github.com/Evrynetlabs/evrynet-client/core/types"
	"github.com/Evrynetlabs/evrynet-client/crypto"
	"github.com/Evrynetlabs/evrynet-client/ethclient"
)

/* These tests are done on a chain with already setup account/ contracts.
To run these test, please deploy your own account/ contract and extract privatekey inorder to get the expected result
Adjust these params to match deployment on local machine:
*/

/*
	Test Send ETH to a normal address
		- No provider signature is required
*/
func TestCreateContractWithProviderAddress(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	provideraddr := common.HexToAddress(providerAddrStr)
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	var option types.CreateAccountOption
	option.ProviderAddress = &provideraddr

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes, option)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	err = errors.New("owner is required")
	assert.Error(t, err, ethClient.SendTransaction(context.Background(), tx))
}

func TestCreateContractWithProviderAndOwner(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	provideraddr := common.HexToAddress(providerAddrStr)
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	var option types.CreateAccountOption
	option.OwnerAddress = &sender
	option.ProviderAddress = &provideraddr

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes, option)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	assert.NoError(t, ethClient.SendTransaction(context.Background(), tx))
}

func TestCreateContractWithoutProviderAddress(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)

	require.NoError(t, ethClient.SendTransaction(context.Background(), tx))
	for i := 0; i < 10; i++ {
		var receipt *types.Receipt
		receipt, err = ethClient.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil {
			assert.Equal(t, uint64(1), receipt.Status)
			assert.NotEqual(t, receipt.ContractAddress, common.Address{})
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func TestCreateContractWithProviderSignature(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	ppk, err := crypto.HexToECDSA(providerPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	tx, err = types.ProviderSignTx(tx, types.HomesteadSigner{}, ppk)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, ethClient.SendTransaction(context.Background(), tx))
}

func TestCreateContractWithProviderAddressWithoutGas(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	provideraddr := common.HexToAddress(providerWithoutGasAddr)
	var option types.CreateAccountOption
	option.ProviderAddress = &provideraddr
	option.OwnerAddress = &sender
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes, option)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	assert.NoError(t, ethClient.SendTransaction(context.Background(), tx))
}

func TestCreateContractWithProviderAddressMustHaveOwnerAddress(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	provideraddr := common.HexToAddress(providerAddrStr)
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	var option types.CreateAccountOption
	option.ProviderAddress = &provideraddr
	option.OwnerAddress = &sender

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes, option)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	assert.Equal(t, strings.ToLower(senderAddrStr), strings.ToLower(tx.Owner().Hex()))
	assert.Equal(t, strings.ToLower(providerAddrStr), strings.ToLower(tx.Provider().Hex()))
}

func TestCreateNormalContractMustHaveNoOwnerAndProviderAddress(t *testing.T) {
	spk, err := crypto.HexToECDSA(senderPK)
	assert.NoError(t, err)
	sender := common.HexToAddress(senderAddrStr)
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	assert.Nil(t, tx.Owner())
	assert.Nil(t, tx.Provider())
}
