package internal

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"

func Broadcast(currNode *maelstrom.Node, body map[string]any) {
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
