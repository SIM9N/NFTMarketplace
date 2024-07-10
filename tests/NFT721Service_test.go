package test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	NFT721 "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/Sim9n/nft-marketplace/services"
	"github.com/Sim9n/nft-marketplace/utils/web3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func Test(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	var (
		url = os.Getenv("ETHER_URL")
		privateKeyHex = os.Getenv("PRIVATE_KEY")
		nftBaseURL = os.Getenv("NFT_BASE_URL")
	)

	addr, pk, err := web3.ImportWallet(privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to import wallet: %v", err)
	}

	//Connect to ether net
	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		t.Fatalf("Failed to connect to ether client: %v", err)
	}
	defer client.Close()

	auth, err := web3.PrepareTransaction(client, addr, pk)
	if err != nil {
		t.Fatalf("Failed to Prepare Transaction: %v", err)
	}

	contractAddr, _, contract, err := NFT721.DeployNFT721(auth, client, "MyNFT", "MyNFT")
	if err != nil {
		t.Fatalf("Failed to deploy contract: %v", err)
	}
	t.Logf("Deployed NFT721 at %s", contractAddr)

	numOfNFTs := 5
	initialPrice := big.NewInt(1000)
	for i := 0; i < numOfNFTs; i++{
		auth, err := web3.PrepareTransaction(client, addr, pk)
		if err != nil {
			t.Fatalf("Failed Prepare Transaction: %v", err)
		}
		_, err =contract.Mint(auth, fmt.Sprintf("%s%d.json", nftBaseURL, i), initialPrice)
		if err == nil {
			t.Logf("Minted NFT %d", i)
		}
	}

	nft721Svc := services.NewNFT721Service(client, NFT721.NFT721ABI, contractAddr.String())
	items := nft721Svc.ListAll()
	for _, item := range items {
		t.Logf("%+v", item)
	}
}