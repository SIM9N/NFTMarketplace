package web3

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	NFT721 "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

func MintNFTs(client *ethclient.Client, address common.Address, privateKey *ecdsa.PrivateKey, baseURL string, contract *NFT721.NFT721) error {
	numOfNFTs := 5
	initialPrice := big.NewInt(1000)
	for i := 0; i < numOfNFTs; i++ {
		auth, err := PrepareTransaction(client, address, privateKey)
		if err != nil {
			return err
		}
		_, err = contract.Mint(auth, fmt.Sprintf("%s/%d.json", baseURL, i), initialPrice)
		if err != nil {
			return err
		}
	}

	return nil
}
