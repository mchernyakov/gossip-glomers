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

[Broadcast d](cmd/broadcast_d.go) uses optimized version, 
where only first receiver propagates data to others.

[Broadcast e](cmd/broadcast_d.go) uses more optimized version: a receiver collects data,
then the received data is transmitted to other nodes according to a scheduled interval.

### Challenge 4: Grow-Only Counter

[The solution](cmd/counter.go) is based on
[CRDT](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type) 
logic.

Articles:
- https://ably.com/blog/crdts-distributed-data-consistency-challenges
- https://crdt.tech/resources

### Challenge 5: Kafka-Style Log

_COMING SOON_

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
