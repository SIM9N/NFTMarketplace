package main

import (
	"context"
	"os"

	"log"
	"net/http"

	"github.com/Sim9n/nft-marketplace/view"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("LISTEN_ADDRESS")

	mux := http.NewServeMux()

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		view.Index().Render(context.Background(), w)
	})

	log.Printf("Listening on %v", port)
	log.Fatal(http.ListenAndServe(port, mux))
}