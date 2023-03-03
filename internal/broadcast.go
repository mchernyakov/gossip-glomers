package internal

import (
	"context"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func Broadcast(currNode *maelstrom.Node, body map[string]any) {
	nodes := currNode.NodeIDs()
	for _, node := range nodes {
		if node == currNode.ID() {
			continue
		}

		dst := node
		go func() {
			for {
				msg, err := currNode.SyncRPC(context.Background(), dst, body)
				println(msg.Body)
				if err == nil {
					break
				}
			}
		}()
	}
}
