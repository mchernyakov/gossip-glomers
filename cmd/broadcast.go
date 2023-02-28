package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"gossip-glomers/internal"

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

		body["type"] = "broadcast_ok"
		delete(body, "message")
		return n.Reply(msg, body)
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
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "topology_ok"
		delete(body, "topology")

		// start Gossip
		gossip.Start(store, n)

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		gossip.Stop()
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}

func getRandomNodes(nodeIDs []string) []string {
	if len(nodeIDs) <= 5 {
		return nodeIDs
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	var res []string
	for i := 0; i < 5; i++ {
		id := r.Intn(len(nodeIDs))
		res = append(res, nodeIDs[id])
	}

	return res
}
