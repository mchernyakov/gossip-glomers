# Gossip Glomers
 
An awesome challenge by [fly.io](https://fly.io/blog/gossip-glomers/)
and [Kyle Kingsbury](https://aphyr.com/).

## Solutions

### Challenge 1: Echo

[Echo](cmd/echo.go) -- a simple echo, nothing fancy.

### Challenge 2: Unique Id Generation

[The solution](cmd/unique_ids.go) generates unique int64 id based on timestamp, 
machine-id and an internal counter.

### Challenge 3: Broadcast

All broadcast solutions are a simplified version of gossip protocol,
[a good article from Martin Fowler](https://martinfowler.com/articles/patterns-of-distributed-systems/gossip-dissemination.html).

[Broadcast a,b,c ](cmd/broadcast.go) -- everyone sends to everyone.

[Broadcast d](cmd/broadcast_d.go) uses an optimized version, 
where only the first receiver propagates data to others.

[Broadcast e](cmd/broadcast_d.go) uses a more optimized version: a receiver collects data,
then the received data is transmitted to other nodes according to a scheduled interval.

### Challenge 4: Grow-Only Counter

[The solution](cmd/counter.go) is based on
[CRDT](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type) 
logic.

Articles:
- https://ably.com/blog/crdts-distributed-data-consistency-challenges
- https://crdt.tech/resources

### Challenge 5: Kafka-Style Log

Solutions are based on single writer/multiple-readers approach:
- using linearizable-KV we can set up our leader (partition-leader), which serves all writes (readers redirect writes to the leader).
- reads are served by every node.

Articles:
- https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-1-storage-mechanics/
- https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-2-data-replication/

### Challenge 6: Totally-Available Transactions

[Txn b](cmd/txn_b.go): based on snapshots and _last-write-wins_ model. 
Also, it broadcasts all writes to other nodes.

[Txn c](cmd/txn_c.go): initially, I thought about something more sophisticated like proper MVCC, 
but then it turned out (and proved by [@elh's solution](https://github.com/elh/gossip-glomers)) 
that snapshots (again) + broadcasting the whole store works.

Articles:
- https://jepsen.io/consistency
- https://vldb.org/pvldb/vol5/p298_per-akelarson_vldb2012.pdf
- https://ebrary.net/64772/computer_science/implementing_read_committed

## Project structure

The project tree

```bash
├── LICENSE
├── Makefile  // build and run commands 
├── README.md
├── build
├── cmd       // <- the solutions are here
├── go.mod
├── go.sum
├── internal  // set of common/reusable go-files
└── maelstrom // maelstrom is in the root, but under gitignore
```

I use Makefile for building and running solutions.

Note that custom names of bin-files are used.
