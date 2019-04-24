package api

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

type bitsharesAPI struct {
	connection       BitsharesAPIConnection
	chainID          objects.ChainID
	username         string
	password         string
	databaseAPIID    int
	historyAPIID     int
	cryptoAPIID      int
	networkNodeAPIID int
	broadcastAPIID   int
	cbLogin          []func()
	cbLogout         []func()
	// Incremental id
	notifyID uint32
}

func (p *bitsharesAPI) OnLogin(fn func()) {
	p.cbLogin = append(p.cbLogin, fn)
}

func (p *bitsharesAPI) OnLogout(fn func()) {
	p.cbLogout = append(p.cbLogout, fn)
}

func (p *bitsharesAPI) connEventHandler(ev BitsharesAPIConnEvent) {
	log.Printf("Got bitshares api connection event %v", ev)
	switch ev {
	case BitsharesAPIConnEventEstablished:
		if err := p.init(); err != nil {
			log.Printf("BitsharesAPI failed to init: %v", err)
			return
		}
		for _, cb := range p.cbLogin {
			go cb()
		}

	case BitsharesAPIConnEventClosed:
		for _, cb := range p.cbLogout {
			go cb()
		}
	}
}

func (p *bitsharesAPI) getAPIID(id *int, identifier string) error {
	if *id != InvalidAPIID {
		return nil
	}
	return p.call(1, identifier, id)
}

func (p *bitsharesAPI) initHistoryAPI() error {
	return p.getAPIID(&p.historyAPIID, "history")
}

func (p *bitsharesAPI) initBroadcastAPI() error {
	return p.getAPIID(&p.broadcastAPIID, "network_broadcast")
}

func (p *bitsharesAPI) initCryptoAPI() error {
	return p.getAPIID(&p.cryptoAPIID, "crypto")
}

func (p *bitsharesAPI) initDatabaseAPI() error {
	return p.getAPIID(&p.databaseAPIID, "database")
}

func (p *bitsharesAPI) initNetworkNodeAPI() (err error) {
	return p.getAPIID(&p.networkNodeAPIID, "network_node")
}

func (p *bitsharesAPI) init() (err error) {
	if ok, err := p.login(); err != nil || !ok {
		if err != nil {
			return errors.Annotate(err, "API login failed")
		}
		return errors.New("login not successful")
	}

	if p.initDatabaseAPI() == nil {
		// cache chainID if db api is availiable
		p.chainID, err = p.getChainID()
	}

	return err
}

func (p *bitsharesAPI) nextSubsciptionID() int {
	return int(atomic.AddUint32(&p.notifyID, 1))
}

func proxyBlockApplied(params []json.RawMessage, handler BlockAppliedHandler) error {
	if len(params) < 1 {
		return errors.New("invalid number of parameters")
	}
	var data objects.BlockID
	if err := json.Unmarshal(params[0], &data); err != nil {
		return errors.Annotate(err, "decode block notification")
	}

	handler(data)
	return nil
}

func (p *bitsharesAPI) SubscribeBlockApplied(cb BlockAppliedHandler) error {
	notifyID := p.nextSubsciptionID()
	if err := p.call(p.databaseAPIID, "set_block_applied_callback", nil, notifyID); err != nil {
		return err
	}

	p.connection.OnNotify(notifyID,
		func(params []json.RawMessage) error {
			return proxyBlockApplied(params, cb)
		})

	return nil
}

func (p *bitsharesAPI) CancelAllSubscriptions() error {
	return p.call(p.databaseAPIID, "cancel_all_subscriptions", nil)
}

//Broadcast a transaction to the network.
//The transaction will be checked for validity in the local database prior to broadcasting.
//If it fails to apply locally, an error will be thrown and the transaction will not be broadcast.
func (p *bitsharesAPI) BroadcastTransaction(tx *objects.Transaction) error {
	return p.call(p.broadcastAPIID, "broadcast_transaction", nil, tx)
}

//GetPotentialSignatures will return the set of all public keys that could possibly sign for a given transaction.
//This call can be used by wallets to filter their set of public keys to just the relevant subset prior to calling
//GetRequiredSignatures to get the minimum subset.
func (p *bitsharesAPI) GetPotentialSignatures(tx *objects.Transaction) ([]objects.PublicKey, error) {
	var result []objects.PublicKey
	err := p.call(p.databaseAPIID, "get_potential_signatures", &result, tx)
	return result, err
}

//GetRequiredSignatures returns the minimum subset of keys required to sign the given transaction
func (p *bitsharesAPI) GetRequiredSignatures(tx *objects.Transaction, keys []objects.PublicKey) ([]objects.PublicKey, error) {
	var result []objects.PublicKey
	err := p.call(p.databaseAPIID, "get_required_signatures", &result, tx, keys)
	return result, err
}

//GetBlock returns a Block by block number.
func (p *bitsharesAPI) GetBlock(number uint64) (*objects.Block, error) {
	var result *objects.Block
	err := p.call(p.databaseAPIID, "get_block", &result, number)
	return result, err
}

//GetTicker returns ticker for specified asset pair
func (p *bitsharesAPI) GetTicker(base, quote string) (*objects.MarketTicker, error) {
	var result *objects.MarketTicker
	err := p.call(p.databaseAPIID, "get_ticker", &result, base, quote, false)
	return result, err
}

//GetAccountByName returns a Account object by username.ListAccountBalances(account objects.GrapheneObject) ([]objects.AssetAmount, error)
func (p *bitsharesAPI) GetAccountByName(name string) (*objects.Account, error) {
	var result *objects.Account
	err := p.call(p.databaseAPIID, "get_account_by_name", &result, name)
	if result == nil {
		return nil, errors.NotFoundf("account '%s'", name)
	}
	return result, err
}

//GetAccounts returns a list of accounts by ID.
func (p *bitsharesAPI) GetAccounts(accounts ...objects.GrapheneObject) ([]objects.Account, error) {
	var result []objects.Account
	err := p.call(p.databaseAPIID, "get_accounts", &result, accounts)
	return result, err
}

//GetAccountCount returns the total number of accounts registered with the blockchain
func (p *bitsharesAPI) GetAccountCount() (uint64, error) {
	var count objects.UInt64
	err := p.call(p.databaseAPIID, "get_account_count", &count)
	return uint64(count), err
}

//GetFullAccounts fetch all objects relevant to the specified accounts, optionally subscribe
func (p *bitsharesAPI) GetFullAccounts(subscribe bool, accounts ...objects.GrapheneObject) ([]*objects.FullAccountResult, error) {
	var result []*objects.FullAccountResult
	err := p.call(p.databaseAPIID, "get_full_accounts", &result, accounts, subscribe)
	return result, err
}

//GetDynamicGlobalProperties
func (p *bitsharesAPI) GetDynamicGlobalProperties() (*objects.DynamicGlobalProperties, error) {
	var result *objects.DynamicGlobalProperties
	err := p.call(p.databaseAPIID, "get_dynamic_global_properties", &result)
	return result, err
}

//GetAccountBalances retrieves AssetAmount objects by given AccountID
func (p *bitsharesAPI) GetAccountBalances(account objects.GrapheneObject, assets ...objects.GrapheneObject) ([]objects.AssetAmount, error) {
	var result []objects.AssetAmount
	ids := objects.GrapheneObjects(assets).ToObjectIDs()
	err := p.call(p.databaseAPIID, "get_account_balances", &result, account.Id(), ids)
	return result, err
}

//ListAssets retrieves assets
//@param lowerBoundSymbol: Lower bound of symbol names to retrieve
//@param limit: Maximum number of assets to fetch, if the constant AssetsListAll is passed, all existing assets will be retrieved.
func (p *bitsharesAPI) ListAssets(lowerBoundSymbol string, limit int) ([]objects.Asset, error) {
	lim := limit
	if limit > AssetsMaxBatchSize || limit == AssetsListAll {
		lim = AssetsMaxBatchSize
	}

	var result []objects.Asset
	err := p.call(p.databaseAPIID, "list_assets", &result, lowerBoundSymbol, lim)
	return result, err
}

//GetRequiredFees calculates the required fee for each operation in the specified asset type.
//If the asset type does not have a valid core_exchange_rate
func (p *bitsharesAPI) GetRequiredFees(ops objects.Operations, feeAsset objects.GrapheneObject) ([]objects.AssetAmount, error) {
	var result []objects.AssetAmount
	err := p.call(p.databaseAPIID, "get_required_fees", &result, ops, feeAsset.Id())
	return result, err
}

//GetLimitOrders returns a slice of LimitOrder objects.
func (p *bitsharesAPI) GetLimitOrders(base, quote objects.GrapheneObject, limit int) (objects.LimitOrders, error) {
	if limit > GetLimitOrdersLimit {
		limit = GetLimitOrdersLimit
	}

	var result objects.LimitOrders
	err := p.call(p.databaseAPIID, "get_limit_orders", &result, base.Id(), quote.Id(), limit)
	return result, err
}

//GetSettleOrders returns a slice of SettleOrder objects.
func (p *bitsharesAPI) GetSettleOrders(assetID objects.GrapheneObject, limit int) ([]objects.SettleOrder, error) {
	if limit > GetSettleOrdersLimit {
		limit = GetSettleOrdersLimit
	}

	var result []objects.SettleOrder
	err := p.call(p.databaseAPIID, "get_settle_orders", &result, assetID.Id(), limit)
	return result, err
}

//GetCallOrders returns a slice of CallOrder objects.
func (p *bitsharesAPI) GetCallOrders(assetID objects.GrapheneObject, limit int) ([]objects.CallOrder, error) {
	if limit > GetCallOrdersLimit {
		limit = GetCallOrdersLimit
	}

	var result []objects.CallOrder
	err := p.call(p.databaseAPIID, "get_call_orders", &result, assetID.Id(), limit)
	return result, err
}

//GetMarginPositions returns a slice of CallOrder objects for the specified account.
func (p *bitsharesAPI) GetMarginPositions(accountID objects.GrapheneObject) ([]objects.CallOrder, error) {
	var result []objects.CallOrder
	err := p.call(p.databaseAPIID, "get_margin_positions", &result, accountID.Id())
	return result, err
}

//GetTradeHistory returns MarketTrade object.
func (p *bitsharesAPI) GetTradeHistory(base, quote objects.GrapheneObject, toTime, fromTime time.Time, limit int) ([]objects.MarketTrade, error) {
	if limit > GetTradeHistoryLimit {
		limit = GetTradeHistoryLimit
	}

	var result []objects.MarketTrade
	err := p.call(p.databaseAPIID, "get_trade_history", &result, base.Id(), quote.Id(), toTime, fromTime, limit)
	return result, err
}

//GetChainID returns the ID of the chain we are connected to.
func (p *bitsharesAPI) GetChainID() objects.ChainID {
	return p.chainID // return cached value
}

func (p *bitsharesAPI) getChainID() (objects.ChainID, error) {
	var result objects.Binary
	err := p.call(p.databaseAPIID, "get_chain_id", &result)
	return objects.ChainID(result), err
}

//GetObjects returns a list of Graphene Objects by ID.
func (p *bitsharesAPI) GetObjects(ids ...objects.GrapheneObject) ([]interface{}, error) {
	params := objects.GrapheneObjects(ids).ToObjectIDs()

	var data []json.RawMessage
	err := p.call(p.databaseAPIID, "get_objects", &data, params)
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, len(data))

	for idx, obj := range data {
		if obj == nil {
			continue
		}

		parsed, err := objects.UnmarshalObject(obj)
		if err != nil {
			return nil, err
		}

		ret[idx] = parsed
	}

	return ret, nil
}

func (p *bitsharesAPI) CancelOrder(orderID objects.GrapheneObject, broadcast bool) (*objects.Transaction, error) {
	var result *objects.Transaction
	err := p.call(p.databaseAPIID, "cancel_order", &result, orderID.Id(), broadcast)
	return result, err
}

func (p *bitsharesAPI) GetRecentTransactionByID(txID objects.TxID) (*objects.Transaction, error) {
	var result *objects.Transaction
	err := p.call(p.databaseAPIID, "get_recent_transaction_by_id", &result, hex.EncodeToString(txID))
	return result, err
}

func (p *bitsharesAPI) GetAccountHistory(account objects.GrapheneObject,
	stop objects.GrapheneObject, start objects.GrapheneObject, limit int) ([]objects.OperationHistory, error) {

	var result []objects.OperationHistory
	err := p.call(p.historyAPIID, "get_account_history", &result, account, stop, limit, start)
	return result, err
}

func (p *bitsharesAPI) GetChronoRelativeHistory(account objects.GrapheneObject,
	start uint32, limit int) ([]objects.OperationHistory, error) {

	var result []objects.OperationHistory
	err := p.call(p.historyAPIID, "get_chrono_relative_history", &result, account, start, limit)
	return result, err
}

func (p *bitsharesAPI) GetNetworkNodeInfo() (objects.NetworkNodeInfo, error) {
	var result objects.NetworkNodeInfo
	err := p.call(p.networkNodeAPIID, "get_info", &result)
	return result, err
}

func (p *bitsharesAPI) GetConnectedPeers() ([]objects.PeerStatus, error) {
	var result []objects.PeerStatus
	err := p.call(p.networkNodeAPIID, "get_connected_peers", &result)
	return result, err
}

func (p *bitsharesAPI) GetAdvancedNodeParameters() (objects.AdvancedNodeParameters, error) {
	var result objects.AdvancedNodeParameters
	err := p.call(p.networkNodeAPIID, "get_advanced_node_parameters", &result)
	return result, err
}

func (p *bitsharesAPI) GetPotentialPeers() ([]objects.PotentialPeerRecord, error) {
	var result []objects.PotentialPeerRecord
	err := p.call(p.networkNodeAPIID, "get_potential_peers", &result)
	return result, err
}

func (p *bitsharesAPI) call(apiID int, method string, result interface{}, args ...interface{}) error {
	return p.connection.Call(apiID, method, result, args...)
}

func (p *bitsharesAPI) login() (bool, error) {
	var result bool
	err := p.call(1, "login", &result, p.username, p.password)
	return result, err
}

func (p *bitsharesAPI) DatabaseAPI() (DatabaseAPI, error) {
	return p, p.initDatabaseAPI()
}
func (p *bitsharesAPI) BroadcastAPI() (BroadcastAPI, error) {
	return p, p.initBroadcastAPI()
}
func (p *bitsharesAPI) HistoryAPI() (HistoryAPI, error) {
	return p, p.initHistoryAPI()
}
func (p *bitsharesAPI) NetworkNodeAPI() (NetworkNodeAPI, error) {
	return p, p.initNetworkNodeAPI()
}
func (p *bitsharesAPI) CryptoAPI() (CryptoAPI, error) {
	return p, p.initCryptoAPI()
}

//New creates a new BitsharesAPI interface.
func New(wsConnection BitsharesAPIConnection, params Params) BitsharesAPI {
	api := &bitsharesAPI{
		username:         params.Username,
		password:         params.Password,
		connection:       wsConnection,
		databaseAPIID:    InvalidAPIID,
		historyAPIID:     InvalidAPIID,
		broadcastAPIID:   InvalidAPIID,
		cryptoAPIID:      InvalidAPIID,
		networkNodeAPIID: InvalidAPIID,
	}
	api.connection.RegisterCallback(api.connEventHandler)
	return api
}
