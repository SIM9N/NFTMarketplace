package test

import (
	"context"
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

	contractAddr, _, contract, err := NFT721.DeployNFT721(auth, client, addr)
	if err != nil {
		t.Fatalf("Failed to deploy contract: %v", err)
	}
	t.Logf("Deployed NFT721 at %s", contractAddr)

	auth, err = web3.PrepareTransaction(client, addr, pk)
	if err != nil {
		t.Fatalf("Failed to Prepare Transaction: %v", err)
	}
	_, err = contract.Mint(auth, "testUrl", big.NewInt(int64(100000)))
	if err != nil {
		t.Fatalf("Failed to Mint: %v", err)
	}

	nft721Svc := services.NewNFT721Service(client, NFT721.NFT721ABI, contractAddr.String())
	count, err := nft721Svc.TokenCount()
	if err != nil {
		t.Fatalf("TokenCount error: %v", err)
	}
	t.Logf("tokenCount: %d", count)

	for i := 1; i <= int(count); i++ {
		data := nft721Svc.GetItemData(int64(i))
		t.Logf("token(%d) %+v", i, data)
	}
}