package internal

import (
	"context"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TxnStoreVc struct {
	store map[float64]*Value
	mu    sync.Mutex
}

func NewTxnStoreVc() *TxnStoreVc {
	return &TxnStoreVc{store: make(map[float64]*Value)}
}

func (s *TxnStoreVc) Update(snapShot map[float64]any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range snapShot {
		val := v.(Value)
		s.writeToStore(k, &val)
	}
}

func (s *TxnStoreVc) updateCurrent(snapShot map[float64]Value) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range snapShot {
		s.writeToStore(k, &v)
	}
}

func (s *TxnStoreVc) writeToStore(key float64, value *Value) {
	ver := value.version
	prev, ok := s.store[key]
	if ok {
		pv := prev.version
		if pv < ver {
			s.store[key] = &Value{
				val:     value.val,
				version: ver,
			}
		}
	} else {
		s.store[key] = &Value{
			val:     value.val,
			version: ver,
		}
	}
}

func (s *TxnStoreVc) PerformTxn(txn []interface{}, n *maelstrom.Node) [][]interface{} {
	var snapShot map[float64]Value
	s.mu.Lock()
	snapShot = deepCopy(s.store)
	s.mu.Unlock()

	var res [][]interface{}
	broadcastData := make(map[float64]Value)

	for _, opr := range txn {
		op := opr.([]any)
		t := op[0].(string)

		var key float64
		var value float64

		if t == "r" {
			key = op[1].(float64)
			v, ok := snapShot[key]
			if !ok {
				res = append(res, []interface{}{t, key, nil})
			} else {
				value = v.val
				res = append(res, []interface{}{t, key, value})
			}
		} else {
			// write
			key = op[1].(float64)
			value = op[2].(float64)

			nv := 0

			prev, ok := snapShot[key]
			if ok {
				nv = prev.version + 1
			}

			record := Value{
				val:     value,
				version: nv,
			}

			snapShot[key] = record

			res = append(res, []interface{}{t, key, value})

			broadcastData[key] = record
		}
	}

	s.updateCurrent(snapShot)

	broadcastStore(n, broadcastData)

	return res
}

func broadcastStore(n *maelstrom.Node, snapShot map[float64]Value) {
	nodes := n.NodeIDs()

	body := make(map[string]any)
	body["type"] = "gossip"
	body["snapshot"] = snapShot

	for _, node := range nodes {
		if n.ID() == node {
			continue
		}

		dst := node
		go func() {
			for {
				_, err := n.SyncRPC(context.Background(), dst, body)
				if err == nil {
					break
				}
			}
		}()

	}
}
