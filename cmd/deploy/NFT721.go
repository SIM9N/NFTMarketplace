package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	NFT721 "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// Import Wallet
func ImportWallet(privateKey string) (common.Address, *ecdsa.PrivateKey) {
	importedPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
 
	publicKey := importedPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
	 log.Fatal("error casting public key to ECDSA")
	}
 
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, importedPrivateKey
}

func PrepareTransaction(client *ethclient.Client, address common.Address, privateKey *ecdsa.PrivateKey) *bind.TransactOpts {
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice((context.Background()))
	if err != nil {
		log.Fatalf("Failed to get gasPrice: %v", err)
	}

	chainID, err := client.ChainID((context.Background()))
	if err != nil {
		log.Fatalf("Failed to get chainID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to auth: %v", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000
	auth.Nonce = big.NewInt(int64(nonce))
 
	return auth
 }
 
func main() {
	log.Println("Deploying NFT721 smart contract")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		url = os.Getenv("ETHER_URL")
		privateKeyHex = os.Getenv("PRIVATE_KEY")
		address, privateKey = ImportWallet(privateKeyHex)
	)

	//Connect to ether net
	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to connect to ether client: %v", err)
	}
	defer client.Close()
	
	auth := PrepareTransaction(client, address, privateKey)

	txAddr, tx, _, err := NFT721.DeployNFT721(auth, client, address)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	log.Println("Contract Address", txAddr.Hex())
	log.Println("Transaction", tx.Hash().Hex())
}