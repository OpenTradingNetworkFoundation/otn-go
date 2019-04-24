package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

type Monitor struct {
	conn      api.BitsharesAPIConnection
	rpc       api.BitsharesAPI
	blockChan chan *objects.BlockID
}

func NewMonitor(nodeAddress string) *Monitor {
	m := &Monitor{
		blockChan: make(chan *objects.BlockID),
	}
	m.conn, m.rpc = api.NewBuilder().Node(nodeAddress).BlockHandler(m.onBlockApplied).Build()
	return m
}

func (m *Monitor) Start() error {
	go m.worker()
	return m.conn.Connect()
}

func (m *Monitor) onBlockApplied(blockID objects.BlockID) {
	log.Printf("New block applied: %d (id=%s)", blockID.BlockNumber(), blockID)
	m.blockChan <- &blockID
}

func (m *Monitor) worker() {
	dbAPI, err := m.rpc.DatabaseAPI()
	if err != nil {
		log.Fatalf("Failed to get database API: %v", err)
	}
	for b := range m.blockChan {
		block, err := dbAPI.GetBlock(uint64(b.BlockNumber()))
		if err != nil {
			log.Printf("Failed to get block: %s", err)
		} else {
			log.Printf("#%08d Timestamp: %s, transactions: %d",
				b.BlockNumber(), block.Timestamp, len(block.Transactions))
		}
	}
}

var (
	nodeAddress string
)

func main() {
	flag.StringVar(&nodeAddress, "node", "ws://127.0.0.1:8090", "Node websocket address")
	flag.Parse()

	monitor := NewMonitor(nodeAddress)

	if err := monitor.Start(); err != nil {
		log.Fatalf("Failed to connect %s", nodeAddress)
	}

	// wait for signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChan:
	}
}
