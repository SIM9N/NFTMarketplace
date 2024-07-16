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
	query, err := htmx.DecodeHTMXQuery(r)
	if err != nil {
		slog.Error("HandleMyNFTPage failed to decode htmx value", "err", err)
		view.MyNFT([]*services.ItemData{}).Render(context.TODO(), w)
		return
	}

	account, ok := query["account"].(string)
	if !ok {
		slog.Error("HandleMyNFTPage account is not a string", "account", query["account"])
		view.MyNFT([]*services.ItemData{}).Render(context.TODO(), w)
		return
	}

	items := h.nft721Svc.ListByAddr(account)
	view.MyNFT(items).Render(context.TODO(), w)
}

func (h *Handler) HandleAccountChangedEvent(w http.ResponseWriter, r *http.Request) {
	body, err := htmx.DecodeHTMXBody(r)
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

	view.Navbar(account).Render(context.Background(), w)
}
