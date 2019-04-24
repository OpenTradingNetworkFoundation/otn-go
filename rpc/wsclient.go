package rpc

import (
	"encoding/json"
	"net"
	"sync/atomic"

	"log"
	"sync"
	"time"

	"github.com/juju/errors"
)

type event struct {
	name string
	data []byte
}

const (
	writeTimeout = 5 * time.Second
	replyTimeout = 10 * time.Second
)

type wsClient struct {
	conn          *RecConn
	url           string
	onError       ErrorFunc
	errors        chan error
	closing       bool
	shutdown      bool
	eventChan     chan event
	writeChan     chan *RPCCall
	currentID     uint64
	wg            sync.WaitGroup
	mutex         sync.Mutex // protects the following
	pending       map[uint64]*RPCCall
	notifyHandler NotifyFunc
}

func NewWebsocketClient(endpointURL string) WebsocketClient {
	cli := wsClient{
		conn:      NewRecConn(),
		pending:   make(map[uint64]*RPCCall),
		errors:    make(chan error, 10),
		eventChan: make(chan event, 10),
		writeChan: make(chan *RPCCall, 1),
		currentID: 1,
		url:       endpointURL,
	}

	return &cli
}

func (p *wsClient) wsEventHandler(e string, data []byte) {
	p.eventChan <- event{
		name: e,
		data: data,
	}
}

func (p *wsClient) Connect() error {
	p.conn.RegisterCallback(p.wsEventHandler)
	if err := p.conn.Dial(p.url); err != nil {
		return err
	}

	p.wg.Add(3)
	go p.monitor()
	go p.receive()
	go p.writer()

	return nil
}

func (p *wsClient) Close() error {
	if p.conn != nil {
		p.closing = true
		p.conn.Close()
		p.wg.Wait()
		p.conn = nil
	}

	return nil
}

func (p *wsClient) monitor() {
	defer p.wg.Done()

	for err := range p.errors {
		if err != nil {
			if p.onError != nil {
				p.onError(err)
			} else {
				log.Println("wsclient error: ", err)
			}
		}
	}
}

func (p *wsClient) writer() {
	defer p.wg.Done()

	for call := range p.writeChan {
		if err := p.conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
			call.done(errors.Annotate(err, "wsclient writer"))
			continue
		}

		if err := p.conn.WriteJSON(call.Request); err != nil {
			call.done(errors.Annotate(err, "JSON request send failed"))
			continue
		}
	}
}

func (p *wsClient) receive() {
	defer func() {
		close(p.errors)
		close(p.writeChan)
		p.wg.Done()
	}()

	for !p.closing {
		if !p.conn.IsConnected() {
			select {
			case e := <-p.eventChan:
				log.Printf("Reconnected, receiving event '%s'\n", e.name)
			case <-time.After(10 * time.Second):
				log.Printf("Retrying connection status")
				continue
			}
		}

		var data rpcFrame
		if err := p.conn.ReadJSON(&data); err != nil {
			if e, ok := err.(*net.OpError); ok {
				if e.Err.Error() == "use of closed network connection" {
					// cancel pending calls
					p.cancelPendingRequests(false)
					// continue loop without notification
					continue
				}
			}

			p.errors <- errors.Annotate(err, "decode in")
			continue
		}

		if data.IsResponse() {
			p.mutex.Lock()
			call, ok := p.pending[data.ID]
			p.mutex.Unlock()

			if ok {
				call.Result = data.Result
				if data.Error != nil {
					call.done(data.Error)
				} else {
					call.done(nil)
				}
			} else {
				p.errors <- errors.Errorf("no corresponding call found for incoming rpc data %#v", data)
				continue
			}
		} else if data.IsNotify() {
			if p.notifyHandler != nil {
				p.notifyHandler(*data.Method, data.Params)
			}
		} else {
			log.Printf("Received unknown websocket message: %#v", data)
		}
	}

	// Terminate pending calls
	p.cancelPendingRequests(true)
}

func (p *wsClient) cancelPendingRequests(shutdown bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.shutdown = shutdown
	for _, call := range p.pending {
		call.done(ErrShutdown)
	}
}

func (p *wsClient) OnNotify(fn NotifyFunc) {
	p.notifyHandler = fn
}

func (p *wsClient) OnError(fn ErrorFunc) {
	p.onError = fn
}

func (p *wsClient) nextID() uint64 {
	return atomic.AddUint64(&p.currentID, 1)
}

func (p *wsClient) Call(method string, result interface{}, args ...interface{}) error {
	if p.shutdown || p.closing {
		return ErrShutdown
	}

	if !p.conn.IsConnected() {
		return ErrNotConnected
	}

	call := &RPCCall{
		Request: rpcRequest{
			Method: method,
			Params: args,
			ID:     p.nextID(),
		},
		Done: make(chan error, 1),
	}

	p.mutex.Lock()
	p.pending[call.Request.ID] = call
	p.mutex.Unlock()

	p.writeChan <- call

	err := call.wait(replyTimeout)

	p.mutex.Lock()
	delete(p.pending, call.Request.ID)
	p.mutex.Unlock()

	if err != nil {
		return err
	}

	if result == nil {
		return nil
	}

	return json.Unmarshal(call.Result, result)
}

func (p *wsClient) RegisterCallback(fn func(event string, data []byte)) {
	p.conn.RegisterCallback(fn)
}
