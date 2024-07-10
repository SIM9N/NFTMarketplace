package web3

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Import Wallet
func ImportWallet(privateKey string) (common.Address, *ecdsa.PrivateKey, error) {
	importedPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return common.Address{}, nil, err
	}
 
	publicKey := importedPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, nil, err
	}
 
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, importedPrivateKey, nil
}

func PrepareTransaction(client *ethclient.Client, address common.Address, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice((context.Background()))
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID((context.Background()))
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000
	auth.Nonce = big.NewInt(int64(nonce))
 
	return auth, nil
}