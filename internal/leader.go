package internal

import (
	"context"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const leaderKey = "leader"

type State struct {
	Mu     sync.RWMutex
	Leader string
}

func NewLeaderState() *State {
	return &State{}
}

func CheckLeader(state *State, n *maelstrom.Node, kv *maelstrom.KV) error {
	if state.Leader != "" {
		return nil
	}

	ctx := context.Background()
	err := kv.CompareAndSwap(ctx, leaderKey, n.ID(), n.ID(), true)
	if err == nil {
		state.Leader = n.ID()
		return nil
	}
	if maelstrom.ErrorCode(err) != maelstrom.PreconditionFailed {
		return fmt.Errorf("leader election: %w", err)
	}

	v, err := kv.Read(ctx, leaderKey)
	if err != nil {
		return fmt.Errorf("read leader: %w", err)
	}
	leader, ok := v.(string)
	if !ok {
		return fmt.Errorf("unexpected leader value: %v", v)
	}
	state.Leader = leader
	return nil
}
