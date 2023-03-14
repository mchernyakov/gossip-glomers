package internal

import (
	"sync"
)

type KafkaStore struct {
	mu      sync.Mutex
	commits map[string]int   // topic -> committed offset
	data    map[string][]int // topic -> data (offset -- index)
}

func NewKafkaStore() *KafkaStore {
	return &KafkaStore{
		mu:      sync.Mutex{},
		commits: make(map[string]int),
		data:    make(map[string][]int),
	}
}

func (store *KafkaStore) ListCommittedOffsets(req []interface{}) map[string]int {
	store.mu.Lock()
	defer store.mu.Unlock()

	res := make(map[string]int)
	for _, topic := range req {
		key := topic.(string)
		res[key] = store.commits[key]
	}

	return res
}

func (store *KafkaStore) CommitOffsets(offsets map[string]any) {
	store.mu.Lock()
	defer store.mu.Unlock()

	for offset, value := range offsets {
		store.commits[offset] = int(value.(float64))
	}
}

func (store *KafkaStore) Send(key string, record float64) int {
	store.mu.Lock()
	defer store.mu.Unlock()

	offset := 0

	topic, ok := store.data[key]
	r := int(record)
	if ok {
		topic := append(topic, r)
		offset = len(topic) - 1
		store.data[key] = topic
	} else {
		store.data[key] = []int{r}
	}

	return offset
}

func (store *KafkaStore) Poll(req map[string]any) map[string][][2]int {
	store.mu.Lock()
	defer store.mu.Unlock()

	res := make(map[string][][2]int)
	for k, v := range req {
		commit := int(v.(float64))

		dataArr := store.data[k]
		if len(dataArr) == 0 {
			continue
		}

		msgs := dataArr[commit:]

		var final [][2]int
		counter := commit
		for _, msg := range msgs {
			r := [2]int{counter, msg}
			final = append(final, r)
			counter++
		}

		if len(final) != 0 {
			res[k] = final
		}
	}
	return res
}
