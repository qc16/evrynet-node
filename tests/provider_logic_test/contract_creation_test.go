package test

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	provideraddrArr := []*common.Address{&provideraddr}
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	var option types.CreateAccountOption
	option.ProviderAddresses = provideraddrArr

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
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	assert.NoError(t, ethClient.SendTransaction(context.Background(), tx))

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
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
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
	provideraddrArr := []*common.Address{&provideraddr}
	var option types.CreateAccountOption
	option.ProviderAddresses = provideraddrArr
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
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
	provideraddrArr := []*common.Address{&provideraddr}
	payLoadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	var option types.CreateAccountOption
	option.ProviderAddresses = provideraddrArr
	option.OwnerAddress = &sender

	ethClient, err := ethclient.Dial(ethRPCEndpoint)
	assert.NoError(t, err)
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
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
	nonce, err := ethClient.NonceAt(context.Background(), sender, nil)
	assert.NoError(t, err)
	tx := types.NewContractCreation(nonce, big.NewInt(0), testGasLimit, big.NewInt(testGasPrice), payLoadBytes)
	tx, err = types.SignTx(tx, types.HomesteadSigner{}, spk)
	assert.NoError(t, err)
	assert.Nil(t, tx.Owner())
	assert.Nil(t, tx.Provider())
}
