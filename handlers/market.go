package handlers

import (
	"context"
	"net/http"

	"github.com/Sim9n/nft-marketplace/view"
)

func (h *Handler) HandleMarketPage(w http.ResponseWriter, r *http.Request) {
	items := h.nft721Svc.ListAll()
	view.Market(items).Render(context.TODO(), w)
}
