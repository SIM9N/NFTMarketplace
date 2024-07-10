package services

import (
	"math/big"

	contract "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ItemData struct {
	TokenId uint64
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

func (svc *NFT721Service) ListAll() []ItemData {
	count, err := svc.TokenCount()
	if err != nil {
		return []ItemData{}
	}
	items := make([]ItemData, count)
	for i := 0; i < int(count); i++ {
		item, err := svc.GetItemData(uint64(i))
		if err == nil {
			items[item.TokenId] = item
		}
	}

	return items
}

func (svc *NFT721Service) TokenCount() (uint64, error) {
	count, err := svc.nft721.NFT721Caller.TokenCount(nil)
	return count.Uint64(), err
}

func (svc *NFT721Service) GetItemData(tokenId uint64) (ItemData, error) {
	ownerAddr, err := svc.nft721.NFT721Caller.OwnerOf(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		return ItemData{}, err
	}

	isListing, err := svc.nft721.NFT721Caller.IsListing(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		return ItemData{}, err
	}


	price, err := svc.nft721.NFT721Caller.Prices(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		return ItemData{}, err
	}

	tokenUrl, err := svc.nft721.NFT721Caller.TokenURI(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		return ItemData{}, err
	}

	return ItemData{
		TokenId: tokenId,
		Owner: ownerAddr.String(),
		IsListing: isListing,
		Price: price.Uint64(),
		Url: tokenUrl,
	}, nil
}
