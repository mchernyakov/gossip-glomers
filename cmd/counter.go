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

const key = "key"

var counter int

var mu sync.Mutex

func main() {

	n := maelstrom.NewNode()

	kv := maelstrom.NewSeqKV(n)

	counter = 0

	n.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["delta"]
		v := m.(float64)

		mu.Lock()
		counter = counter + int(v)
		go func() {
			err := kv.Write(context.Background(), key, counter)
			if err != nil {
				return
			}
		}()
		mu.Unlock()

		for _, next := range n.NodeIDs() {
			dst := next
			go func() {
				for {
					err := n.Send(dst, body)
					if err == nil {
						break
					}
				}
			}()
		}

		rsp := &model.SimpleResp{Type: "add_ok"}
		return n.Reply(msg, rsp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		val, err := kv.Read(context.Background(), key)
		if err != nil {
			body["value"] = 0
		} else {
			body["value"] = val
		}

		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
