package main

import (
	"encoding/json"
	"log"
	"os"

	"gossip-glomers/internal"
	"gossip-glomers/internal/model"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()

	store := internal.NewSimpleStore()
	gossip := internal.NewGossip()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["message"]
		v := m.(float64)
		store.Add(v)
		_, ok := body["gossip"]
		if !ok {
			body["gossip"] = true
			internal.Broadcast(n, body)
		}

		rsp := &model.SimpleResp{Type: "broadcast_ok"}
		return n.Reply(msg, rsp)
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
		return n.Reply(msg, rsp)
	})

	if err := n.Run(); err != nil {
		gossip.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
