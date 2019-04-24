package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/juju/errors"
)

type NotifyFunc func(method string, params json.RawMessage)
type ErrorFunc func(error)

var (
	ErrShutdown    = errors.New("connection is shut down")
	ErrCallTimeout = errors.New("rpc request timed out")
)

type WebsocketClient interface {
	OnError(fn ErrorFunc)
	OnNotify(fn NotifyFunc)
	Call(method string, result interface{}, args ...interface{}) error
	Close() error
	Connect() error
	RegisterCallback(fn func(event string, data []byte))
}

type RPCCall struct {
	Method  string
	Request rpcRequest
	Result  json.RawMessage
	Done    chan error
}

func (call *RPCCall) done(err error) {
	select {
	case call.Done <- err:
		// ok
	default:
		log.Println("rpc: discarding Call reply due to insufficient Done chan capacity")
	}
}

func (call *RPCCall) wait(timeout time.Duration) error {
	select {
	case err := <-call.Done:
		return err
	case <-time.After(timeout):
		return ErrCallTimeout
	}
}

type rpcRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	ID     uint64      `json:"id"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("JSONRPC error: %s (%d)", e.Message, e.Code)
}

type rpcFrame struct {
	// rpcResponse
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *rpcError       `json:"error,omitempty"`

	// rpcNotify
	Method *string         `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
}

func (f *rpcFrame) IsResponse() bool {
	return f.ID != 0 && f.Method == nil
}

func (f *rpcFrame) IsNotify() bool {
	return f.ID == 0 && f.Method != nil
}
