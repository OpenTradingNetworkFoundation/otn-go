package api

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/rpc"
)

type apiConnection struct {
	wsClient  rpc.WebsocketClient
	callbacks []func(event BitsharesAPIConnEvent)

	mutexNotify sync.Mutex // protects the following
	notifyFns   map[int]APINotifyFunc
}

func (p *apiConnection) RegisterCallback(cb func(BitsharesAPIConnEvent)) {
	p.callbacks = append(p.callbacks, cb)
}

func (p *apiConnection) Call(apiID interface{}, method string, result interface{}, args ...interface{}) error {
	if len(args) == 0 {
		args = EmptyParams
	}
	return p.wsClient.Call("call", result, apiID, method, args)
}

func (p *apiConnection) OnNotify(subscriberID int, notifyFn APINotifyFunc) {
	p.mutexNotify.Lock()
	defer p.mutexNotify.Unlock()
	p.notifyFns[subscriberID] = notifyFn
}

func (p *apiConnection) OnError(errorFn func(err error)) {
	p.wsClient.OnError(errorFn)
}

func (p *apiConnection) Connect() (err error) {
	if err := p.wsClient.Connect(); err != nil {
		return errors.Annotate(err, "websocket connect failed")
	}
	return nil
}

//Close() shuts down the API apiConnection.
func (p *apiConnection) Close() error {
	return p.wsClient.Close()
}

func (p *apiConnection) wsEventsCallback(event string, data []byte) {
	log.Printf("Got websocket event %s", event)
	t := BitsharesAPIConnEventUnknown
	switch event {
	case "ConnectionClosed":
		t = BitsharesAPIConnEventClosed
		p.resetSubscriptions()
	case "ConnectionEstablished":
		t = BitsharesAPIConnEventEstablished
	}

	for _, cb := range p.callbacks {
		go cb(t)
	}
}

func (p *apiConnection) notificationHandler(method string, params json.RawMessage) {
	var subscriber int
	var args []json.RawMessage

	data := [2]interface{}{&subscriber, &args}
	if err := json.Unmarshal(params, &data); err != nil {
		log.Printf("Failed to unmarshal parameters: %v (%s)", err, string(params))
		return
	}

	p.mutexNotify.Lock()
	fn := p.notifyFns[subscriber]
	p.mutexNotify.Unlock()

	if fn != nil {
		fn(args)
	}
}

func (p *apiConnection) resetSubscriptions() {
	p.mutexNotify.Lock()
	defer p.mutexNotify.Unlock()
	p.notifyFns = make(map[int]APINotifyFunc)
}

// NewConnection creates new websocket connection with Bitshares API.
func NewConnection(wsEndpointURL string) BitsharesAPIConnection {
	conn := &apiConnection{
		wsClient:  rpc.NewWebsocketClient(wsEndpointURL),
		notifyFns: make(map[int]APINotifyFunc),
	}
	conn.wsClient.RegisterCallback(conn.wsEventsCallback)
	conn.wsClient.OnNotify(conn.notificationHandler)
	return conn
}
