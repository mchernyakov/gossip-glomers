package internal

import (
	"context"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Value struct {
	val     float64
	version int
}

type TxnStoreV struct {
	store map[float64]*Value
	mu    sync.Mutex
}

func NewTxnStoreV() *TxnStoreV {
	return &TxnStoreV{store: make(map[float64]*Value)}
}

func (s *TxnStoreV) Update(key any, value any, version any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := key.(float64)
	val := value.(float64)
	ver := int(version.(float64))

	prev, ok := s.store[k]
	if ok {
		pv := prev.version
		if pv < ver {
			s.store[k] = &Value{
				val:     val,
				version: ver,
			}
		}
	} else {
		s.store[k] = &Value{
			val:     val,
			version: ver,
		}
	}
}

func (s *TxnStoreV) PerformTxn(txn []interface{}, n *maelstrom.Node) [][]interface{} {
	var snapShot map[float64]*Value
	s.mu.Lock()
	snapShot = deepCopy(s.store)
	s.mu.Unlock()

	var res [][]interface{}

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
			prev, ok := snapShot[key]
			if ok {
				nv := prev.version + 1
				snapShot[key] = &Value{
					val:     value,
					version: nv,
				}

				broadcast(n, key, value, nv)
			} else {
				snapShot[key] = &Value{
					val:     value,
					version: 0,
				}

				broadcast(n, key, value, 0)
			}

			res = append(res, []interface{}{t, key, value})
		}
	}
	return res
}

func deepCopy(store map[float64]*Value) map[float64]*Value {
	res := make(map[float64]*Value)
	for k, v := range store {
		res[k] = v
	}

	return res
}

func broadcast(n *maelstrom.Node, key float64, value float64, version int) {
	nodes := n.NodeIDs()

	body := make(map[string]any)
	body["type"] = "gossip"
	body["key"] = key
	body["value"] = value
	body["version"] = version

	for _, node := range nodes {
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
