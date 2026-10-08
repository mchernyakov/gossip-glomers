package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gossip-glomers/internal"
	"gossip-glomers/internal/model"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type broadcastReq struct {
	Message float64 `json:"message"`
}

type readResp struct {
	Type     string    `json:"type"`
	Messages []float64 `json:"messages"`
}

func main() {
	n := maelstrom.NewNode()

	store := internal.NewSimpleStore()
	bc := internal.NewBroadcast(600*time.Millisecond, time.Second)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body broadcastReq
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if store.Add(body.Message) {
			bc.Add(body.Message)
		}
		return n.Reply(msg, &model.SimpleResp{Type: "broadcast_ok"})
	})

	n.Handle("gossip", func(msg maelstrom.Message) error {
		var body internal.GossipMsg
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		store.AddAll(body.Messages)
		return n.Reply(msg, &model.SimpleResp{Type: "gossip_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return n.Reply(msg, &readResp{Type: "read_ok", Messages: store.ReadAll()})
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		bc.Start(n)
		return n.Reply(msg, &model.SimpleResp{Type: "topology_ok"})
	})

	if err := n.Run(); err != nil {
		bc.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
