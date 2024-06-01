package main

import (
	"context"
	"log"
	"math/big"
	"os"

	contract "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	url := os.Getenv("ETHER_URL")
	addrHex := os.Getenv("ETHER_ADDRESS")
	addr := common.HexToAddress(addrHex)
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	
	log.Print("Hello from contract main")
	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to connect to ether client: %v", err)
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(context.Background(), addr)
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
	auth.GasLimit = 300000
	auth.Nonce = big.NewInt(int64(nonce))

	txAddr, tx, c, err := contract.DeployContract(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	log.Println("Deployed contract")
	log.Println(txAddr.Hex())
	log.Println(tx.Hash().Hex())

	greet, err := c.Greet(&bind.CallOpts{
		From: addr,
	})
	if err != nil {
		log.Fatalf("Failed to call Greet: %v", err)
	}
	log.Printf("greet %v", greet)
}