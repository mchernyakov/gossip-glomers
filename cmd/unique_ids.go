package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Counter struct {
	seq      atomic.Int64
	nodeId   uint32
	lastTime int64
	mu       sync.Mutex
}

func main() {
	s := rand.NewSource(time.Now().UnixNano())
	randId := rand.New(s).Intn(100-1) + 1

	counter := &Counter{
		seq:      atomic.Int64{},
		nodeId:   uint32(randId),
		lastTime: 0,
		mu:       sync.Mutex{},
	}

	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["id"] = generateId(counter)
		body["type"] = "generate_ok"

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}

func generateId(counter *Counter) int64 {
	ts := getTime(counter)
	var id = ts << 22
	id = id | int64(counter.nodeId)<<15
	next := counter.seq.Add(1)
	id = id | (next % 2048)
	return id
}

func getTime(counter *Counter) int64 {
	counter.mu.Lock()
	defer counter.mu.Unlock()

	currentTime := time.Now().UnixMilli()
	lastTime := counter.lastTime

	for currentTime == lastTime {
		time.Sleep(time.Millisecond)
		currentTime = time.Now().UnixMilli()
	}
	counter.lastTime = currentTime
	return currentTime
}
