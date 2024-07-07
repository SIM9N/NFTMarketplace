package main

import (
	"log"
	"log/slog"
	"os"

	"net/http"

	"github.com/Sim9n/nft-marketplace/handlers"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("LISTEN_ADDRESS")

	mux := http.NewServeMux()

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /", handlers.HandleIndex)

	mux.HandleFunc("POST /events/onAccountConnected", handlers.HandleAccountChangedEvent)

	logger.Info("Server Started", "port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}