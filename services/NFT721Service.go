package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"sync"

	contract "github.com/Sim9n/nft-marketplace/contracts/gen"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ItemDataMetadata struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ItemData struct {
	TokenId   uint64
	Owner     string
	IsListing bool
	Price     uint64
	Url       string
	MetaData  *ItemDataMetadata
}

type NFT721Service struct {
	etherClient *ethclient.Client
	abi         string
	nft721      *contract.NFT721
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

func (svc *NFT721Service) ListAll() []*ItemData {
	count, err := svc.TokenCount()
	if err != nil {
		return []*ItemData{}
	}

	items := make([]*ItemData, count)
	wg := sync.WaitGroup{}
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			item, err := svc.GetItemData(uint64(i))
			if err == nil {
				items[item.TokenId] = item
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return items
}

func (svc *NFT721Service) ListByAddr(addr string) []*ItemData {
	count, err := svc.TokenCount()
	if err != nil {
		return []*ItemData{}
	}

	items := []*ItemData{}
	wg := sync.WaitGroup{}
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			item, err := svc.GetItemData(uint64(i))
			if err == nil && strings.ToLower(item.Owner) == addr {
				items = append(items, item)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return items
}

func (svc *NFT721Service) TokenCount() (uint64, error) {
	count, err := svc.nft721.NFT721Caller.TokenCount(nil)
	return count.Uint64(), err
}

func (svc *NFT721Service) GetItemData(tokenId uint64) (*ItemData, error) {
	ownerAddr, err := svc.nft721.NFT721Caller.OwnerOf(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		slog.Error("GetItemData failed to fetch owner address", "tokenId", tokenId, "error", err)
		return nil, err
	}

	isListing, err := svc.nft721.NFT721Caller.IsListing(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		slog.Error("GetItemData failed to fetch IsListing", "tokenId", tokenId, "error", err)
		return nil, err
	}

	price, err := svc.nft721.NFT721Caller.Prices(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		slog.Error("GetItemData failed to fetch Prices", "tokenId", tokenId, "error", err)
		return nil, err
	}

	tokenUrl, err := svc.nft721.NFT721Caller.TokenURI(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		slog.Error("GetItemData failed to fetch TokenURI", "tokenId", tokenId, "error", err)
		return nil, err
	}

	metaData, err := svc.fetchItemMetadata(tokenUrl)
	if err != nil {
		slog.Error("GetItemData failed to fetch item metadata", "tokenId", tokenId, "error", err)
		return nil, err
	}

	return &ItemData{
		TokenId:   tokenId,
		Owner:     ownerAddr.String(),
		IsListing: isListing,
		Price:     price.Uint64(),
		Url:       tokenUrl,
		MetaData:  metaData,
	}, nil
}

func (svc *NFT721Service) fetchItemMetadata(url string) (*ItemDataMetadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch item metadata at %s, non-200 status code", url)
	}

	metadata := &ItemDataMetadata{}
	err = json.NewDecoder(resp.Body).Decode(metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
