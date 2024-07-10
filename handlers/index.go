package handlers

import (
	"context"
	"net/http"

	view "github.com/Sim9n/nft-marketplace/view"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	view.Index().Render(context.Background(), w)
}