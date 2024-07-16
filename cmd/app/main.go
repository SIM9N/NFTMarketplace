package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"net/http"

	NFT721 "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/Sim9n/nft-marketplace/handlers"
	"github.com/Sim9n/nft-marketplace/services"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		port         = os.Getenv("LISTEN_ADDRESS")
		url          = os.Getenv("ETHER_URL")
		contractAddr = os.Getenv("CONTRACT_ADDRESS")
	)

	mux := http.NewServeMux()

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	client, err := ethclient.DialContext(context.TODO(), url)
	if err != nil {
		log.Fatalf("Failed to connect to ether client: %v", err)
	}
	defer client.Close()

	nft721Svc := services.NewNFT721Service(client, NFT721.NFT721ABI, contractAddr)
	handler := handlers.New(nft721Svc)

	mux.HandleFunc("GET /", handler.HandleIndex)
	mux.HandleFunc("GET /market", handler.HandleMarketPage)
	mux.HandleFunc("GET /my-nft", handler.HandleMyNFTPage)

	mux.HandleFunc("POST /events/onAccountConnected", handler.HandleAccountChangedEvent)

	logger.Info("Server Started", "port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
