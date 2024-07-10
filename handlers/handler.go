package handlers

import (
	"github.com/Sim9n/nft-marketplace/services"
)

type Handler struct {
	nft721Svc *services.NFT721Service
}

func New(nft721Svc *services.NFT721Service) *Handler {
	return &Handler{
		nft721Svc: nft721Svc,
	}
}