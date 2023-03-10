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

func main() {

	n := maelstrom.NewNode()

	store := internal.NewSimpleStore()
	gossip := internal.NewGossip(100 * time.Millisecond)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["message"]
		v := m.(float64)
		store.Add(v)

		rsp := &model.SimpleResp{Type: "broadcast_ok"}
		return n.Reply(msg, rsp)
	})

	n.Handle("gossip", func(msg maelstrom.Message) error {
		var body internal.GossipMsg
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		store.AddAll(body.Messages)

		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["messages"] = store.ReadAll()
		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		rsp := &model.SimpleResp{Type: "topology_ok"}
		gossip.Start(store, n, 5)
		return n.Reply(msg, rsp)
	})

	if err := n.Run(); err != nil {
		gossip.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
