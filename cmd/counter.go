package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"gossip-glomers/internal/model"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	key   = "key"
	round = "round"
)

type Counter struct {
	value   int
	mu      sync.Mutex
	round   map[float64]bool
	current int
}

func main() {

	n := maelstrom.NewNode()

	kv := maelstrom.NewSeqKV(n)

	counter := &Counter{value: 0,
		round: make(map[float64]bool)}

	n.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["delta"]
		v := int(m.(float64))

		counter.mu.Lock()
		for {
			val, err0 := kv.Read(context.Background(), key)
			if err0 != nil {
				val = 0
			}

			fn := val.(int) + v
			err := kv.CompareAndSwap(context.Background(), key, val, fn, true)
			if err == nil {
				counter.value = fn
				break
			}
		}

		go func() {
			body["type"] = "sync"
			body["new-val"] = counter.value
			for _, next := range n.NodeIDs() {
				dst := next
				go func() {
					for {
						_, err := n.SyncRPC(context.Background(), dst, body)
						if err == nil {
							break
						}
					}
				}()
			}
		}()
		counter.mu.Unlock()

		rsp := &model.SimpleResp{Type: "add_ok"}
		return n.Reply(msg, rsp)
	})

	n.Handle("sync", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["new-val"]
		v := int(m.(float64))

		counter.mu.Lock()
		if counter.value < v {
			counter.value = v
		}
		counter.mu.Unlock()

		rsp := &model.SimpleResp{Type: "sync_ok"}
		return n.Reply(msg, rsp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		counter.mu.Lock()

		val, err := kv.Read(context.Background(), key)
		if err != nil {
			body["value"] = 0
		} else {
			body["value"] = val

			if val.(int) > counter.value {
				body["value"] = val
			} else {
				body["value"] = counter.value
			}
		}
		counter.mu.Unlock()

		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
