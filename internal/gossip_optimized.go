package internal

import (
	"context"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastEntity struct {
	mu             sync.Mutex
	ticker         *time.Ticker
	broadcastStore map[float64]bool
	mainStore      *SimpleStore
	doneChan       chan bool
}

func NewBroadcast(d time.Duration, mainStore *SimpleStore) *BroadcastEntity {
	return &BroadcastEntity{
		mu:             sync.Mutex{},
		ticker:         time.NewTicker(d),
		broadcastStore: make(map[float64]bool),
		doneChan:       make(chan bool),
		mainStore:      mainStore,
	}
}

func (bc *BroadcastEntity) Add(val float64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.broadcastStore[val] = true
}

func (bc *BroadcastEntity) Start(node *maelstrom.Node) {
	go func() {
		for {
			select {
			case <-bc.doneChan:
				bc.ticker.Stop()
				return
			case <-bc.ticker.C:
				bc.doBroadcast(node)
			}
		}
	}()
}

func (bc *BroadcastEntity) doBroadcast(currNode *maelstrom.Node) {
	bcData := bc.getAll()
	if len(bcData) == 0 {
		return
	}
	data := append(bcData, bc.mainStore.ReadAll()...)
	body := &GossipMsg{
		Type:     "gossip",
		Messages: data,
	}

	nodes := currNode.NodeIDs()
	for _, neib := range nodes {
		if neib == currNode.ID() {
			continue
		}

		dst := neib
		go func() {
			for {
				_, err := currNode.SyncRPC(context.Background(), dst, body)
				if err == nil {
					break
				}
			}
		}()
	}
}

func (bc *BroadcastEntity) getAll() []float64 {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	var all []float64
	for key := range bc.broadcastStore {
		all = append(all, key)
	}
	bc.broadcastStore = make(map[float64]bool)
	return all
}

func (bc *BroadcastEntity) Stop() {
	bc.doneChan <- true
}

func SimpleBroadcast(currNode *maelstrom.Node, body map[string]any, data []float64) {
	body["messages"] = data
	nodes := currNode.NodeIDs()
	for _, node := range nodes {
		if node == currNode.ID() {
			continue
		}

		dst := node
		go func() {
			for {
				err := currNode.Send(dst, body)
				if err == nil {
					break
				}
			}
		}()
	}
}
