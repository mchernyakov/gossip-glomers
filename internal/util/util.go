package util

import (
	"math/rand"
	"time"
)

const numNodes = 5

func GetRandomNodes(nodeIDs []string, num int) []string {
	if len(nodeIDs) <= num {
		return nodeIDs
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	var res []string
	for i := 0; i < num; i++ {
		id := r.Intn(len(nodeIDs))
		res = append(res, nodeIDs[id])
	}

	return res
}

func GetRandomNodesDefault(nodeIDs []string) []string {
	return GetRandomNodes(nodeIDs, numNodes)
}
