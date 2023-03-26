package internal

import (
	"context"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type State struct {
	Mu     sync.RWMutex
	Leader string
}

func NewLeaderState() *State {
	return &State{
		Mu:     sync.RWMutex{},
		Leader: "",
	}
}

func CheckLeader(state *State, n *maelstrom.Node, kv *maelstrom.KV) {
	if state.Leader == "" {
		err := kv.CompareAndSwap(context.Background(), "leader", n.ID(), n.ID(), true)
		if err != nil {
			leader, err0 := kv.Read(context.Background(), "leader")
			if err0 != nil {
				_ = fmt.Errorf(err0.Error())
			} else {
				state.Leader = leader.(string)
			}
		} else {
			state.Leader = n.ID()
		}
	}
}
