package internal

import "sync"

type SimpleStore struct {
	set map[float64]bool
	mu  sync.Mutex
}

func NewSimpleStore() *SimpleStore {
	return &SimpleStore{set: make(map[float64]bool)}
}

func (store *SimpleStore) Add(key float64) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	_, ok := store.set[key]
	if ok {
		return false
	} else {
		store.set[key] = true
		return true
	}
}

func (store *SimpleStore) AddAll(keys []float64) bool {
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, key := range keys {
		store.set[key] = true
	}

	return true
}

func (store *SimpleStore) AddAll0(keys []interface{}) bool {
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, key := range keys {
		v := key.(float64)
		store.set[v] = true
	}

	return true
}

func (store *SimpleStore) ReadAll() []float64 {
	store.mu.Lock()
	defer store.mu.Unlock()

	var all []float64
	for key := range store.set {
		all = append(all, key)
	}
	return all
}
