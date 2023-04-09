package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()
	store := NewSimpleTxnStore()

	n.Handle("txn", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		txn := body["txn"].([]interface{})
		res := store.PerformTxn(txn)
		body["in_reply_to"] = body["msg_id"]
		body["txn"] = res
		body["type"] = "txn_ok"
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}

type SimpleTxnStore struct {
	store map[float64]float64
	mu    sync.Mutex
}

func NewSimpleTxnStore() *SimpleTxnStore {
	return &SimpleTxnStore{store: make(map[float64]float64)}
}

func (s *SimpleTxnStore) PerformTxn(txn []interface{}) [][]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res [][]interface{}

	for _, opr := range txn {
		op := opr.([]any)
		t := op[0].(string)

		var key float64
		var val float64

		if t == "r" {
			key = op[1].(float64)
			val = s.store[key]
		} else {
			key = op[1].(float64)
			val = op[2].(float64)
			s.store[key] = val
		}
		res = append(res, []interface{}{t, key, val})
	}
	return res
}
