package main

import (
	"context"
	"log"
	"os"

	NFT721 "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/Sim9n/nft-marketplace/utils/web3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)
 
func main() {
	log.Println("Deploying NFT721 smart contract")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		url = os.Getenv("ETHER_URL")
		privateKeyHex = os.Getenv("PRIVATE_KEY")
	)
	
	address, privateKey, err := web3.ImportWallet(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to import wallet: %v", err)
	}

	//Connect to ether net
	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to connect to ether client: %v", err)
	}
	defer client.Close()
	
	auth, err := web3.PrepareTransaction(client, address, privateKey)
	if err != nil {
		log.Fatalf("Failed Prepare Transaction: %v", err)
	}

	txAddr, tx, _, err := NFT721.DeployNFT721(auth, client, address)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	log.Println("Contract Address", txAddr.Hex())
	log.Println("Transaction", tx.Hash().Hex())
}