package util

import (
	"math/rand"
	"time"
)

const numNodes = 3

func GetRandomNodes(nodeIDs []string) []string {
	if len(nodeIDs) <= numNodes {
		return nodeIDs
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	var res []string
	for i := 0; i < numNodes; i++ {
		id := r.Intn(len(nodeIDs))
		res = append(res, nodeIDs[id])
	}

	return res
}
