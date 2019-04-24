package consul

import (
	"context"
	"crypto/sha256"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

type KVPair = api.KVPair
type KeyChangeHandler func(kv *KVPair)

type KeyWatch struct {
	api    *api.Client
	cancel context.CancelFunc
	ctx    context.Context
	wg     sync.WaitGroup

	changes chan *KVPair
	C       <-chan *KVPair
}

func NewKeyWatch(api *api.Client) *KeyWatch {
	c := make(chan *KVPair)
	o := &KeyWatch{api: api, changes: c, C: c}
	o.ctx, o.cancel = context.WithCancel(context.Background())
	return o
}

func (o *KeyWatch) AddHandler(key string, handler KeyChangeHandler) {
	o.wg.Add(1)
	go o.poller(key, handler)
}

func (o *KeyWatch) Add(key string) {
	o.AddHandler(key, func(kv *KVPair) {
		o.changes <- kv
	})
}

func (o *KeyWatch) poller(key string, handler KeyChangeHandler) {
	kv := o.api.KV()
	q := api.QueryOptions{}
	var wait time.Duration
	var lastHash [sha256.Size]byte
	defer o.wg.Done()

	for {
		select {
		case <-o.ctx.Done():
			return
		case <-time.After(wait):
		}
		pair, _, err := kv.Get(key, q.WithContext(o.ctx))
		if err != nil || pair == nil {
			wait = 3 * time.Second
			continue
		}
		wait = 0
		// call handler only if content changed
		hash := sha256.Sum256(pair.Value)
		if hash != lastHash {
			lastHash = hash
			handler(pair)
		}
		// update wait index to current value
		q.WaitIndex = pair.ModifyIndex
	}
}

func (o *KeyWatch) Shutdown() {
	go func() {
		// drain channel
		for range o.changes {
		}
	}()
	o.cancel()
	o.wg.Wait()
	close(o.changes)
}
