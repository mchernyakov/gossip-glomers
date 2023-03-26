package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	"gossip-glomers/internal"
	"gossip-glomers/internal/model"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()

	state := internal.NewLeaderState()
	kafka := internal.NewKafkaStore()
	kv := maelstrom.NewLinKV(n)

	n.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		state.Mu.Lock()
		defer state.Mu.Unlock()

		internal.CheckLeader(state, n, kv)

		if state.Leader == n.ID() {
			key := body["key"].(string)
			value := body["msg"].(float64)
			offset := kafka.Add(key, value)

			body["type"] = "send_update"
			// update followers
			for _, f := range n.NodeIDs() {
				if f == n.ID() {
					continue
				}

				respL, err := n.SyncRPC(context.Background(), f, body)
				if err != nil {
					return err
				}

				var bd map[string]any
				if err := json.Unmarshal(respL.Body, &bd); err != nil {
					return err
				}
				fo := bd["offset"]
				fv := int(fo.(float64))
				if fv != offset {
					return errors.New("consistency problem : send. leader:" + strconv.Itoa(offset) + " , follower " + strconv.Itoa(fv))
				}
			}

			resp := make(map[string]any)
			resp["type"] = "send_ok"
			resp["offset"] = offset
			return n.Reply(msg, resp)
		} else {
			respL, err := n.SyncRPC(context.Background(), state.Leader, body)
			if err != nil {
				return err
			}
			return n.Reply(msg, respL.Body)
		}
	})

	n.Handle("poll", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		state.Mu.RLock()
		defer state.Mu.RUnlock()

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

		state.Mu.RLock()
		defer state.Mu.RUnlock()

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

		state.Mu.Lock()
		defer state.Mu.Unlock()

		internal.CheckLeader(state, n, kv)

		if state.Leader == n.ID() {
			data := body["offsets"].(map[string]any)
			kafka.CommitOffsets(data)

			body["type"] = "commit_offsets_update"
			// update followers
			for _, f := range n.NodeIDs() {
				if f == n.ID() {
					continue
				}

				_, err := n.SyncRPC(context.Background(), f, body)
				if err != nil {
					return err
				}
			}
		} else {
			_, err := n.SyncRPC(context.Background(), state.Leader, body)
			if err != nil {
				return err
			}
		}

		rsp := &model.SimpleResp{Type: "commit_offsets_ok"}
		return n.Reply(msg, rsp)
	})

	n.Handle("send_update", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		key := body["key"].(string)
		value := body["msg"].(float64)
		offset := kafka.Add(key, value)

		resp := make(map[string]any)
		resp["type"] = "send_update_ok"
		resp["offset"] = offset
		return n.Reply(msg, resp)
	})

	n.Handle("commit_offsets_update", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		data := body["offsets"].(map[string]any)
		kafka.CommitOffsets(data)

		rsp := &model.SimpleResp{Type: "commit_offsets_update_ok"}
		return n.Reply(msg, rsp)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
