package internal

import "sync"

type TxnStore struct {
	store map[float64]float64
	mu    sync.Mutex
}

func NewTxnStore() *TxnStore {
	return &TxnStore{store: make(map[float64]float64)}
}

func (s *TxnStore) PerformTxn(txn []interface{}) [][]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res [][]interface{}

	for _, opr := range txn {
		op := opr.([]any)
		t := op[0].(string)

		var key float64
		var val float64

		if t == "r" {
			key = op[1].(float64)
			val = s.store[key]
		} else {
			key = op[1].(float64)
			val = op[2].(float64)
			s.store[key] = val
		}
		res = append(res, []interface{}{t, key, val})
	}
	return res
}
