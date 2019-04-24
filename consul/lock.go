package consul

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	DefaultConsulLockRetryInterval = 5 * time.Second
)

type ConsulLockOptions struct {
	Key           string
	FakeLock      bool
	OnStateChange func(bool)
	RetryInterval time.Duration
}

type ConsulLock struct {
	lock *api.Lock
	opts *ConsulLockOptions

	// Lock management
	mutex      sync.Mutex
	inProgress bool
	isLocked   bool
	stopChan   chan struct{}
}

func lockValue() string {
	return fmt.Sprintf("Lock:PID=%d", os.Getpid())
}

func NewConsulLock(client *api.Client, opts *ConsulLockOptions) (*ConsulLock, error) {
	lockOpts := &api.LockOptions{
		Key:         opts.Key,
		SessionName: lockValue(),
	}
	if opts.FakeLock {
		lockOpts.Key = "fake/lock"
	}
	lock, err := client.LockOpts(lockOpts)
	if err != nil {
		return nil, err
	}

	if opts.RetryInterval == 0 {
		opts.RetryInterval = DefaultConsulLockRetryInterval
	}

	consulLock := &ConsulLock{
		lock: lock,
		opts: opts,
	}

	return consulLock, nil
}

func (p *ConsulLock) watchLock(leaderCh <-chan struct{}, stopChan <-chan struct{}) {
	select {
	case <-leaderCh:
		// Lock failure, update state
		p.isLocked = false
		// Reset lock state
		p.lock.Unlock()
		// Notify subscribers
		p.opts.OnStateChange(false)
	case <-stopChan:
	}
}

func (p *ConsulLock) aquireLock(stopChan chan struct{}) {
	defer func() {
		p.mutex.Lock()
		p.inProgress = false
		p.mutex.Unlock()
	}()

	for {
		leaderCh, err := p.lock.Lock(stopChan)
		if err != nil {
			log.Print("Lock failed: ", err)
		} else {
			p.isLocked = true
			p.opts.OnStateChange(true)
			p.watchLock(leaderCh, stopChan)
			break
		}

		select {
		case <-stopChan:
			return
		case <-time.After(p.opts.RetryInterval):
		}
	}
}

func (p *ConsulLock) AsyncLock() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.opts.FakeLock {
		if !p.isLocked {
			p.isLocked = true
			p.opts.OnStateChange(true)
		}
		return
	}

	if p.inProgress {
		return
	}

	p.stopChan = make(chan struct{})
	p.inProgress = true

	go p.aquireLock(p.stopChan)
}

func (p *ConsulLock) AsyncUnlock() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	wasLocked := p.isLocked

	if p.opts.FakeLock {
		if p.isLocked {
			p.isLocked = false
			p.opts.OnStateChange(false)
		}
		return wasLocked
	}

	if p.inProgress {
		close(p.stopChan)
		p.inProgress = false
	}

	if p.isLocked {
		p.lock.Unlock()
		p.isLocked = false
	}

	return wasLocked
}

func (p *ConsulLock) CancelPendingLock() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.inProgress {
		close(p.stopChan)
		p.inProgress = false
	}
}
