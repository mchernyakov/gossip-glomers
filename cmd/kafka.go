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

	kafka := internal.NewKafkaStore()

	n.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		key := body["key"].(string)
		value := body["msg"].(float64)
		offset := kafka.Send(key, value)

		resp := make(map[string]any)
		resp["type"] = "send_ok"
		resp["offset"] = offset
		return n.Reply(msg, resp)
	})

	n.Handle("poll", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		data := body["offsets"].(map[string]any)
		res := kafka.Poll(data)

		resp := make(map[string]any)
		resp["type"] = "poll_ok"
		resp["msgs"] = res
		return n.Reply(msg, resp)
	})

	n.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		data := body["keys"].([]interface{})
		res := kafka.ListCommittedOffsets(data)

		resp := make(map[string]any)
		resp["type"] = "list_committed_offsets_ok"
		resp["offsets"] = res
		return n.Reply(msg, resp)
	})

	n.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		data := body["offsets"].(map[string]any)
		kafka.CommitOffsets(data)

		rsp := &model.SimpleResp{Type: "commit_offsets_ok"}
		return n.Reply(msg, rsp)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
