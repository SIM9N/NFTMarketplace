package handlers

import (
	"context"
	"net/http"

	"github.com/Sim9n/nft-marketplace/view"
)

func HandleMarketPage(w http.ResponseWriter, r *http.Request) {
	view.Market().Render(context.TODO(), w)
}