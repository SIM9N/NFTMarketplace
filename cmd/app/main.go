package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"regexp"

	"net/http"
	"net/url"

	component "github.com/Sim9n/nft-marketplace/view"
	"github.com/joho/godotenv"
)


func isJson(input string) bool {
	match, _ := regexp.MatchString(`^\s*(\{.*\}|\[.*\])\s*$`, input)
	return match
}

func DecodeHTMXValue(r *http.Request) (map[string]interface{}, error) {
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	urlValues, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	jsonData := make(map[string]interface{})
	for k, v := range urlValues {
		if len(v) == 0 {
			continue
		}

		var value any

		if isJson(v[0]) {
			json.Unmarshal([]byte(v[0]), &value)
		}else {
			value = v[0]
		}

		jsonData[k] = value
	}

	return jsonData, nil
}

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

	mux.HandleFunc("POST /events/onAccountConnected", func(w http.ResponseWriter, r *http.Request) {
		body, err := DecodeHTMXValue(r)
		if err != nil {
			logger.Warn("onAccountConnected failed to decode htmx value", "err", err)
			component.Navbar("").Render(context.Background(), w)
			return
		}

		account, ok := body["account"].(string); 
		if !ok {
			logger.Warn("onAccountConnected account is not a string", "account", body["account"])
			component.Navbar("").Render(context.Background(), w)
			return
		}

		logger.Info("onAccountConnected", "account", account)
		component.Navbar(account).Render(context.Background(), w)
	})

	logger.Info("Server Started", "port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}