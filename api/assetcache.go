package api

import (
	"log"
	"sync"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

type AssetCache struct {
	api          DatabaseAPI
	mx           sync.Mutex
	assetsById   map[objects.GrapheneID]*objects.Asset
	assetsByName map[string]*objects.Asset
}

func NewAssetCache(api DatabaseAPI) *AssetCache {
	return &AssetCache{
		api:          api,
		assetsById:   make(map[objects.GrapheneID]*objects.Asset),
		assetsByName: make(map[string]*objects.Asset),
	}
}

func (c *AssetCache) loadByID(id *objects.GrapheneID) *objects.Asset {
	assets, err := c.api.GetObjects(id)
	if err != nil {
		log.Println("Failed to load asset:", err)
		return nil
	}

	asset := assets[0].(objects.Asset)

	c.assetsById[*id] = &asset
	c.assetsByName[asset.Symbol] = &asset

	return &asset
}

func (c *AssetCache) loadBySymbol(name string) *objects.Asset {
	assets, err := c.api.ListAssets(name, 1)
	if err != nil {
		log.Println("Failed to load asset:", err)
		return nil
	}

	if len(assets) == 0 || assets[0].Symbol != name {
		log.Printf("Asset '%s' not found", name)
		return nil
	}

	asset := assets[0]

	c.assetsById[asset.ID] = &asset
	c.assetsByName[asset.Symbol] = &asset

	return &asset
}

func (c *AssetCache) GetByID(id *objects.GrapheneID) *objects.Asset {
	c.mx.Lock()
	defer c.mx.Unlock()

	asset, ok := c.assetsById[*id]
	if !ok {
		return c.loadByID(id)
	}

	return asset
}

func (c *AssetCache) GetBySymbol(name string) *objects.Asset {
	c.mx.Lock()
	defer c.mx.Unlock()

	asset, ok := c.assetsByName[name]
	if !ok {
		return c.loadBySymbol(name)
	}

	return asset
}

func (c *AssetCache) GetAsset(nameOrID string) *objects.Asset {
	var asset *objects.Asset

	// check by ID first
	var ID objects.GrapheneID
	if ID.FromString(nameOrID) == nil {
		asset = c.GetByID(&ID)
	}

	// try by name if not found by ID
	if asset == nil {
		asset = c.GetBySymbol(nameOrID)
	}

	return asset
}
