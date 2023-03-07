package internal

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// var roundCounter atomic.Int64

// var roundMap map[uint64]bool

// var mu sync.Mutex

func Broadcast(currNode *maelstrom.Node, body map[string]any, data []float64) {
	// c := roundCounter.Add(1)
	// body["round"] = c
	body["messages"] = data
	nodes := currNode.NodeIDs()
	for _, node := range nodes {
		if node == currNode.ID() {
			continue
		}

		dst := node
		go func() {
			for {
				err := currNode.Send(dst, body)
				if err == nil {
					break
				}
			}
		}()
	}
}
