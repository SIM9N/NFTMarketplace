package handlers

import (
	"context"
	"net/http"

	"github.com/Sim9n/nft-marketplace/view"
)

func HandleMyNFTPage(w http.ResponseWriter, r *http.Request) {
	view.MyNFT().Render(context.TODO(), w)
}