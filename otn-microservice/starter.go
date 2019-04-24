package otn

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/consul"
)

// StarterEvent shows OTN infrastructure event type
type StarterEvent int

// Possible events from OTN infrastructure
const (
	StarterEventExit StarterEvent = iota
	StarterEventAPIConnected
	StarterEventAPIDisconnected
	StarterEventLocked
	StarterEventUnlocked
)

// StarterConfig specifies all necessary data for proper OTN API init
type StarterConfig struct {
	InstanceLock string `json:"instance_lock"`
	TrustedNode  string `json:"trusted_node"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// Starter deals with OTN infrastructure
type Starter struct {
	cfg        *StarterConfig
	connection api.BitsharesAPIConnection
	api        api.BitsharesAPI
	service    Microservice

	consulLock *consul.ConsulLock
	appChan    chan StarterEvent
}

// NewStarter creates properly initialized Starter object
func NewStarter(service Microservice, cfg *StarterConfig) *Starter {
	return &Starter{
		cfg:     cfg,
		service: service,
		appChan: make(chan StarterEvent, 1),
	}
}

func (a *Starter) lockStateChanged(isLocked bool) {
	log.Printf("lockStateChanged: %v", isLocked)
	if isLocked {
		a.appChan <- StarterEventLocked
	} else {
		a.appChan <- StarterEventUnlocked
	}
}

// Run interaction with OTN infrastructure (API and Consul)
func (a *Starter) Run(doneChan chan struct{}) error {
	a.connection = api.NewConnection(a.cfg.TrustedNode)
	a.api = api.New(a.connection, api.Params{Username: a.cfg.Username, Password: a.cfg.Password})
	a.api.OnLogin(func() {
		a.appChan <- StarterEventAPIConnected
	})
	a.api.OnLogout(func() {
		a.appChan <- StarterEventAPIDisconnected
	})

	consulClient, err := consul.NewClient()
	if err != nil {
		return errors.Annotate(err, "consul.NewClient")
	}

	lockOpts := &consul.ConsulLockOptions{
		Key:           a.cfg.InstanceLock,
		FakeLock:      a.cfg.InstanceLock == "",
		OnStateChange: a.lockStateChanged,
	}

	a.consulLock, err = consul.NewConsulLock(consulClient, lockOpts)
	if err != nil {
		return errors.Annotate(err, "ConsulLock")
	}

	if err := a.connection.Connect(); err != nil {
		return errors.Annotate(err, fmt.Sprintf("Connection to API node %s failed", a.cfg.TrustedNode))
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	isConnected := false

	for {
		select {
		case event := <-a.appChan:
			switch event {
			case StarterEventExit:
				a.consulLock.AsyncUnlock()
				return nil
			case StarterEventLocked:
				a.service.Start(a.api)
			case StarterEventUnlocked:
				if isConnected {
					a.service.Stop()
					// if lock failed, try to aquire it again
					a.consulLock.AsyncLock()
				}
			case StarterEventAPIConnected:
				isConnected = true
				log.Print("Waiting lock...")
				a.consulLock.AsyncLock()
			case StarterEventAPIDisconnected:
				isConnected = false
				a.service.Stop()
				a.consulLock.AsyncUnlock()
			}
		case s := <-signalChan:
			log.Printf("Got signal %d, exiting", s)
			a.service.SignalHandler(s)
			a.appChan <- StarterEventExit
		case <-doneChan:
			log.Printf("Stopping service")
			a.service.Stop()
			a.appChan <- StarterEventExit
		}
	}
}
