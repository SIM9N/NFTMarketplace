package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Sim9n/nft-marketplace/services"
	"github.com/Sim9n/nft-marketplace/utils/htmx"
	view "github.com/Sim9n/nft-marketplace/view"
)

type Handler struct {
	nft721Svc *services.NFT721Service
}

func New(nft721Svc *services.NFT721Service) *Handler {
	return &Handler{
		nft721Svc: nft721Svc,
	}
}

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	view.Index().Render(context.Background(), w)
}

func (h *Handler) HandleMarketPage(w http.ResponseWriter, r *http.Request) {
	items := h.nft721Svc.ListAll()
	view.Market(items).Render(context.TODO(), w)
}

func (h *Handler) HandleMyNFTPage(w http.ResponseWriter, r *http.Request) {
	view.MyNFT().Render(context.TODO(), w)
}

func (h *Handler) HandleAccountChangedEvent(w http.ResponseWriter, r *http.Request) {
	body, err := htmx.DecodeHTMXValue(r)
	if err != nil {
		slog.Warn("onAccountConnected failed to decode htmx value", "err", err)
		view.Navbar("").Render(context.Background(), w)
		return
	}

	account, ok := body["account"].(string)
	if !ok {
		slog.Warn("onAccountConnected account is not a string", "account", body["account"])
		view.Navbar("").Render(context.Background(), w)
		return
	}

	slog.Info("onAccountConnected", "account", account)
	view.Navbar(account).Render(context.Background(), w)
}
