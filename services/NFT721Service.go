package services

import (
	"log/slog"
	"math/big"

	contract "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ItemData struct {
	Owner string
	IsListing bool
	Price uint64
	Url string
}

type NFT721Service struct {
	etherClient *ethclient.Client
	abi string
	nft721 *contract.NFT721
}

func NewNFT721Service(client *ethclient.Client, abi, address string) *NFT721Service {
	nft721, err := contract.NewNFT721(common.HexToAddress(address), client)
	if err != nil {
		panic("Failed to construct nft721 instance")
	}
	return &NFT721Service{
		client,
		abi,
		nft721,
	}
}

func (svc *NFT721Service) TokenCount() (int64, error) {
	count, err := svc.nft721.NFT721Caller.TokenCount(nil)
	return count.Int64(), err
}

func (svc *NFT721Service) GetItemData(tokenId int64) ItemData {
	ownerAddr, err := svc.nft721.NFT721Caller.OwnerOf(nil, big.NewInt(tokenId))
	if err != nil {
		slog.Error("GetItemData OwnerOf", "error", err)
	}

	isListing, err := svc.nft721.NFT721Caller.IsListing(nil, big.NewInt(tokenId))
	if err != nil {
		slog.Error("GetItemData IsListing", "error", err)
	}


	price, err := svc.nft721.NFT721Caller.Prices(nil, big.NewInt(tokenId))
	if err != nil {
		slog.Error("GetItemData Prices", "error", err)
	}

	tokenUrl, err := svc.nft721.NFT721Caller.TokenURI(nil, big.NewInt(tokenId))
	if err != nil {
		slog.Error("GetItemData TokenUrl", "error", err)
	}

	return ItemData{
		Owner: ownerAddr.String(),
		IsListing: isListing,
		Price: price.Uint64(),
		Url: tokenUrl,
	}
}
