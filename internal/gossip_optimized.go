package internal

import (
	"context"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastEntity struct {
	mu        sync.Mutex
	interval  time.Duration
	timeout   time.Duration
	fresh     map[float64]struct{}
	pending   map[string]map[float64]struct{}
	inFlight  map[string]bool
	done      chan struct{}
	startOnce sync.Once
	stopOnce  sync.Once
}

func NewBroadcast(interval, timeout time.Duration) *BroadcastEntity {
	return &BroadcastEntity{
		interval: interval,
		timeout:  timeout,
		fresh:    make(map[float64]struct{}),
		pending:  make(map[string]map[float64]struct{}),
		inFlight: make(map[string]bool),
		done:     make(chan struct{}),
	}
}

func (bc *BroadcastEntity) Add(val float64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.fresh[val] = struct{}{}
}

func (bc *BroadcastEntity) Start(node *maelstrom.Node) {
	bc.startOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(bc.interval)
			defer ticker.Stop()
			for {
				select {
				case <-bc.done:
					return
				case <-ticker.C:
					bc.flush(node)
				}
			}
		}()
	})
}

func (bc *BroadcastEntity) Stop() {
	bc.stopOnce.Do(func() { close(bc.done) })
}

func (bc *BroadcastEntity) flush(node *maelstrom.Node) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for _, peer := range node.NodeIDs() {
		if peer == node.ID() {
			continue
		}
		if len(bc.fresh) > 0 {
			p, ok := bc.pending[peer]
			if !ok {
				p = make(map[float64]struct{})
				bc.pending[peer] = p
			}
			for v := range bc.fresh {
				p[v] = struct{}{}
			}
		}
		if bc.inFlight[peer] || len(bc.pending[peer]) == 0 {
			continue
		}
		batch := make([]float64, 0, len(bc.pending[peer]))
		for v := range bc.pending[peer] {
			batch = append(batch, v)
		}
		bc.inFlight[peer] = true
		go bc.send(node, peer, batch)
	}
	bc.fresh = make(map[float64]struct{})
}

func (bc *BroadcastEntity) send(node *maelstrom.Node, peer string, batch []float64) {
	ctx, cancel := context.WithTimeout(context.Background(), bc.timeout)
	defer cancel()

	_, err := node.SyncRPC(ctx, peer, &GossipMsg{Type: "gossip", Messages: batch})

	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.inFlight[peer] = false
	if err != nil {
		return
	}
	for _, v := range batch {
		delete(bc.pending[peer], v)
	}
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
