package api

import (
	"encoding/json"
	"time"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

const (
	InvalidAPIID         = -1
	AssetsListAll        = -1
	AssetsMaxBatchSize   = 100
	GetCallOrdersLimit   = 100
	GetLimitOrdersLimit  = 100
	GetSettleOrdersLimit = 100
	GetTradeHistoryLimit = 100
)

var (
	EmptyParams = []interface{}{}
)

type BitsharesAPIConnEvent int

const (
	BitsharesAPIConnEventUnknown BitsharesAPIConnEvent = iota
	BitsharesAPIConnEventEstablished
	BitsharesAPIConnEventClosed
)

type APINotifyFunc func(params []json.RawMessage) error

// Params specifies some API connection options
type Params struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BitsharesAPIConnection interface {
	Connect() error
	Close() error
	OnError(func(error))
	OnNotify(subscriberID int, handler APINotifyFunc)
	Call(apiID interface{}, method string, result interface{}, args ...interface{}) error
	RegisterCallback(fn func(BitsharesAPIConnEvent))
}

type BlockAppliedHandler func(objects.BlockID)

type DatabaseAPI interface {
	// Objects
	GetObjects(objectIDs ...objects.GrapheneObject) ([]interface{}, error)
	ListAssets(lowerBoundSymbol string, limit int) ([]objects.Asset, error)
	GetTicker(base, quote string) (*objects.MarketTicker, error)

	// Subscriptions
	SubscribeBlockApplied(BlockAppliedHandler) error
	CancelAllSubscriptions() error

	// Blocks & transactions
	GetBlock(number uint64) (*objects.Block, error)
	GetRecentTransactionByID(tx objects.TxID) (*objects.Transaction, error)

	// Globals
	GetChainID() objects.ChainID
	GetDynamicGlobalProperties() (*objects.DynamicGlobalProperties, error)

	// Accounts
	GetAccountBalances(account objects.GrapheneObject, assets ...objects.GrapheneObject) ([]objects.AssetAmount, error)
	GetAccountByName(name string) (*objects.Account, error)
	GetAccounts(accountIDs ...objects.GrapheneObject) ([]objects.Account, error)
	GetAccountCount() (uint64, error)
	GetFullAccounts(subscribe bool, accountIDs ...objects.GrapheneObject) ([]*objects.FullAccountResult, error)

	// Markets
	CancelOrder(orderID objects.GrapheneObject, broadcast bool) (*objects.Transaction, error)
	GetMarginPositions(accountID objects.GrapheneObject) ([]objects.CallOrder, error)
	GetCallOrders(assetID objects.GrapheneObject, limit int) ([]objects.CallOrder, error)
	GetLimitOrders(base, quote objects.GrapheneObject, limit int) (objects.LimitOrders, error)
	GetSettleOrders(assetID objects.GrapheneObject, limit int) ([]objects.SettleOrder, error)
	GetTradeHistory(base, quote objects.GrapheneObject, toTime, fromTime time.Time, limit int) ([]objects.MarketTrade, error)

	// Autority/Validation
	GetRequiredFees(ops objects.Operations, feeAsset objects.GrapheneObject) ([]objects.AssetAmount, error)
	GetPotentialSignatures(tx *objects.Transaction) ([]objects.PublicKey, error)
	GetRequiredSignatures(tx *objects.Transaction, keys []objects.PublicKey) ([]objects.PublicKey, error)
}

// BroadcastAPI implements
type BroadcastAPI interface {
	BroadcastTransaction(tx *objects.Transaction) error
}

// HistoryAPI implements
type HistoryAPI interface {
	GetAccountHistory(account objects.GrapheneObject, stop objects.GrapheneObject, start objects.GrapheneObject, limit int) ([]objects.OperationHistory, error)
	GetChronoRelativeHistory(account objects.GrapheneObject, start uint32, limit int) ([]objects.OperationHistory, error)
}

// NetworkNodeAPI interface provides access to internal OTN blockchain networking information
type NetworkNodeAPI interface {
	GetNetworkNodeInfo() (objects.NetworkNodeInfo, error)
	GetConnectedPeers() ([]objects.PeerStatus, error)
	GetAdvancedNodeParameters() (objects.AdvancedNodeParameters, error)
	GetPotentialPeers() ([]objects.PotentialPeerRecord, error)
}

// CryptoAPI interface is intended to produce some low level OTN cryptography operations
type CryptoAPI interface {
}

type BitsharesAPI interface {
	OnLogin(fn func())
	OnLogout(fn func())

	DatabaseAPI() (DatabaseAPI, error)
	BroadcastAPI() (BroadcastAPI, error)
	HistoryAPI() (HistoryAPI, error)
	NetworkNodeAPI() (NetworkNodeAPI, error)
	CryptoAPI() (CryptoAPI, error)
}
