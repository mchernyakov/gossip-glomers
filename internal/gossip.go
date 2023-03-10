package internal

import (
	"time"

	"gossip-glomers/internal/util"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Gossip struct {
	ticker   *time.Ticker
	doneChan chan bool
}

type GossipMsg struct {
	Type     string    `json:"type"`
	Messages []float64 `json:"messages"`
}

type GossipSingleMsg struct {
	Type    string  `json:"type"`
	Message float64 `json:"message"`
}

func NewGossip(d time.Duration) *Gossip {
	return &Gossip{
		ticker:   time.NewTicker(d),
		doneChan: make(chan bool),
	}
}

func (gossip Gossip) Start(store *SimpleStore, node *maelstrom.Node, maxNodes int) {
	go func() {
		for {
			select {
			case <-gossip.doneChan:
				gossip.ticker.Stop()
				return
			case <-gossip.ticker.C:
				err := doGossip(store, node, maxNodes)
				if err != nil {
					return
				}
			}
		}
	}()
}

func doGossip(store *SimpleStore, currNode *maelstrom.Node, maxNodes int) error {
	nodeIds := currNode.NodeIDs()
	neibs := util.GetRandomNodes(nodeIds, maxNodes)

	data := store.ReadAll()

	body := &GossipMsg{
		Type:     "gossip",
		Messages: data,
	}

	for _, neib := range neibs {
		if neib == currNode.ID() {
			continue
		}

		dst := neib
		go func() {
			for {
				err := currNode.Send(dst, body)
				if err == nil {
					break
				}
			}
		}()
	}
	return nil
}

func (gossip Gossip) Stop() {
	gossip.doneChan <- true
}
