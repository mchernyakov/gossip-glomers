package util

import (
	"math/rand"
	"time"
)

func GetRandomNodes(nodeIDs []string) []string {
	if len(nodeIDs) <= 3 {
		return nodeIDs
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	var res []string
	for i := 0; i < 3; i++ {
		id := r.Intn(len(nodeIDs))
		res = append(res, nodeIDs[id])
	}

	return res
}
