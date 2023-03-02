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

func NewGossip() *Gossip {
	return &Gossip{
		ticker:   time.NewTicker(50 * time.Millisecond),
		doneChan: make(chan bool),
	}
}

func (gossip Gossip) Start(store *SimpleStore, node *maelstrom.Node) {
	go func() {
		for {
			select {
			case <-gossip.doneChan:
				gossip.ticker.Stop()
				return
			case <-gossip.ticker.C:
				err := doGossip(store, node)
				if err != nil {
					return
				}
			}
		}
	}()
}

func doGossip(store *SimpleStore, node *maelstrom.Node) error {
	nodeIds := node.NodeIDs()
	neibs := util.GetRandomNodes(nodeIds)

	data := store.ReadAll()

	gmsg := &GossipMsg{
		Type:     "gossip",
		Messages: data,
	}

	for _, neib := range neibs {
		err := node.Send(neib, gmsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gossip Gossip) Stop() {
	gossip.doneChan <- true
}
