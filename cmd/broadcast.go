package main

import (
	"encoding/json"
	"fmt"
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
	// neibsState := model.NewNeibsState()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		m := body["message"]
		v := m.(float64)

		if store.Add(v) {
			neibs := getRandomNodes(n.NodeIDs())

			for _, neib := range neibs {
				err := n.Send(neib, msg.Body)
				if err != nil {
					return err
				}
			}
		}

		// for _, neib := range neibsState.GetNeibs() {
		//	err := n.Send(neib, v)
		//	if err != nil {
		//		return err
		//	}
		// }

		body["type"] = "broadcast_ok"
		delete(body, "message")
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

		// topology := body["topology"].(map[string]interface{})
		// neibs := getNodesToBroadcast(topology, n.ID())
		// neibsState.AddNeibs(neibs)

		body["type"] = "topology_ok"
		delete(body, "topology")
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
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

func getNodesToBroadcast(topology map[string]interface{}, nodeId string) []string {

	for key, value := range topology {
		fmt.Printf("%s value is %v\n", key, value)
	}

	data := topology[nodeId].([]interface{})

	var neibs []string
	for _, val := range data {
		neibs = append(neibs, val.(string))
	}
	return neibs
}
