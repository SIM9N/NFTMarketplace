package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
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
		url           = os.Getenv("ETHER_URL")
		privateKeyHex = os.Getenv("PRIVATE_KEY")
		nftBaseURL    = os.Getenv("NFT_BASE_URL")
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

	txAddr, tx, contract, err := NFT721.DeployNFT721(auth, client, "SIMON_MOK_NFT", "SM_NFT")
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	numOfNFTs := 5
	initialPrice := big.NewInt(1000)
	for i := 0; i < numOfNFTs; i++ {
		auth, err := web3.PrepareTransaction(client, address, privateKey)
		if err != nil {
			log.Fatalf("Failed Prepare Transaction: %v", err)
		}
		_, err = contract.Mint(auth, fmt.Sprintf("%s/%d.json", nftBaseURL, i), initialPrice)
		if err == nil {
			log.Printf("Minted NFT %d", i)
		}
	}

	log.Println("Contract Address", txAddr.Hex())
	log.Println("Transaction", tx.Hash().Hex())
}
