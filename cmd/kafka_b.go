package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"gossip-glomers/internal"
	"gossip-glomers/internal/model"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type State struct {
	mu     sync.Mutex
	leader string
}

func main() {

	n := maelstrom.NewNode()

	state := &State{
		mu:     sync.Mutex{},
		leader: "",
	}
	kafka := internal.NewKafkaStore()
	kv := maelstrom.NewLinKV(n)

	n.Handle("send", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		state.mu.Lock()
		defer state.mu.Unlock()

		checkLeader(state, n, kv)

		if state.leader == n.ID() {
			key := body["key"].(string)
			value := body["msg"].(float64)
			offset := kafka.Add(key, value)

			resp := make(map[string]any)
			resp["type"] = "send_ok"
			resp["offset"] = offset
			return n.Reply(msg, resp)
		} else {
			respL, err := n.SyncRPC(context.Background(), state.leader, body)
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

		state.mu.Lock()
		defer state.mu.Unlock()

		checkLeader(state, n, kv)

		if state.leader == n.ID() {
			data := body["offsets"].(map[string]any)
			res := kafka.Poll(data)

			resp := make(map[string]any)
			resp["type"] = "poll_ok"
			resp["msgs"] = res
			return n.Reply(msg, resp)
		} else {
			respL, err := n.SyncRPC(context.Background(), state.leader, body)
			if err != nil {
				return err
			}
			return n.Reply(msg, respL.Body)
		}
	})

	n.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		state.mu.Lock()
		defer state.mu.Unlock()

		checkLeader(state, n, kv)
		resp := make(map[string]any)

		if state.leader == n.ID() {
			data := body["keys"].([]interface{})
			res := kafka.ListCommittedOffsets(data)
			resp["offsets"] = res
		} else {
			respL, err := n.SyncRPC(context.Background(), state.leader, body)
			if err != nil {
				return err
			}
			var bd map[string]any
			if err := json.Unmarshal(respL.Body, &bd); err != nil {
				return err
			}
			resp["offsets"] = bd["offsets"]
		}

		resp["type"] = "list_committed_offsets_ok"
		return n.Reply(msg, resp)

	})

	n.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		state.mu.Lock()
		defer state.mu.Unlock()

		checkLeader(state, n, kv)

		if state.leader == n.ID() {
			data := body["offsets"].(map[string]any)
			kafka.CommitOffsets(data)
		} else {
			_, err := n.SyncRPC(context.Background(), state.leader, body)
			if err != nil {
				return err
			}
		}

		rsp := &model.SimpleResp{Type: "commit_offsets_ok"}
		return n.Reply(msg, rsp)

	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}

func checkLeader(state *State, n *maelstrom.Node, kv *maelstrom.KV) {
	if state.leader == "" {
		err := kv.CompareAndSwap(context.Background(), "leader", n.ID(), n.ID(), true)
		if err != nil {
			leader, err0 := kv.Read(context.Background(), "leader")
			if err0 != nil {
				_ = fmt.Errorf(err0.Error())
			} else {
				state.leader = leader.(string)
			}
		} else {
			state.leader = n.ID()
		}
	}
}
