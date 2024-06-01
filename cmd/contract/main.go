package main

import (
	"context"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	url := os.Getenv("ETHER_URL")
	
	log.Print("Hello from contract main")
	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to connect to ether client: %v", err)
	}
	defer client.Close()


	addrHex := "0x4075CC5BF83Fc7c0Ed3005ceFd0161A56FD7300E"
	addr := common.HexToAddress(addrHex)

	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	log.Printf(`address: %s, balance: %d`, addrHex, balance)
}