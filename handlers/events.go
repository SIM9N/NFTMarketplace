package handlers

import (
	"context"
	"log/slog"
	"net/http"

	htmx "github.com/Sim9n/nft-marketplace/utils/htmx"
	view "github.com/Sim9n/nft-marketplace/view"
)

func HandleAccountChangedEvent(w http.ResponseWriter, r *http.Request) {
	body, err := htmx.DecodeHTMXValue(r)
	if err != nil {
		slog.Warn("onAccountConnected failed to decode htmx value", "err", err)
		view.Navbar("").Render(context.Background(), w)
		return
	}

	account, ok := body["account"].(string); 
	if !ok {
		slog.Warn("onAccountConnected account is not a string", "account", body["account"])
		view.Navbar("").Render(context.Background(), w)
		return
	}

	slog.Info("onAccountConnected", "account", account)
	view.Navbar(account).Render(context.Background(), w)
}