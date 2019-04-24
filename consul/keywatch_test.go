package consul_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/opentradingnetworkfoundation/otn-go/consul"
	"github.com/stretchr/testify/assert"
)

func TestObserver(t *testing.T) {
	testkey := "test/observer/testkey"
	client, _ := consul.NewClient()
	kv := client.KV()
	kw := consul.NewKeyWatch(client)
	keyCreatedReceived := make(chan bool)

	kv.Delete(testkey, nil)
	kw.Add(testkey)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		kv.Put(&api.KVPair{Key: testkey, Value: []byte("new")}, nil)
		<-keyCreatedReceived

		for i := 0; i < 10; i++ {
			value := &api.KVPair{
				Key:   testkey,
				Value: []byte(fmt.Sprintf("v%d", i)),
			}
			kv.Put(value, nil)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	lastidx := -1
	var idx int

	for v := range kw.C {
		keyValue := string(v.Value)
		if keyValue == "new" {
			keyCreatedReceived <- true
		} else {
			fmt.Sscanf(keyValue, "v%d", &idx)
			// check ordering
			assert.True(t, idx > lastidx)
			lastidx = idx
			if idx == 9 {
				break
			}
		}
	}

	assert.Equal(t, 9, idx)

	wg.Wait()
	kw.Shutdown()

	kv.Delete(testkey, nil)
}
