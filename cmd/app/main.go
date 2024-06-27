package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"net/http"

	"github.com/Sim9n/nft-marketplace/view"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("LISTEN_ADDRESS")

	mux := http.NewServeMux()

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		component.Index().Render(context.Background(), w)
	})
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		component.Login().Render(context.Background(), w)
	})

	logger.Info("Server Started", "port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}