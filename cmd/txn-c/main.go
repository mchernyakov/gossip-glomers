package main

import (
	"encoding/json"
	"log"
	"os"

	"gossip-glomers/internal"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()
	store := internal.NewTxnStoreVc()

	n.Handle("txn", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		txn := body["txn"].([]interface{})
		res := store.PerformTxn(txn, n)
		body["in_reply_to"] = body["msg_id"]
		body["txn"] = res
		body["type"] = "txn_ok"
		return n.Reply(msg, body)
	})

	n.Handle("gossip", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		snapShot := body["snapshot"].(map[float64]any)

		store.Update(snapShot)

		body["type"] = "gossip_ok"
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
