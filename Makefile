SHELL=/bin/bash -o pipefail

build-echo:
	go build -o build/bin/echo ./cmd/echo

build-unique-ids:
	go build -o build/bin/unique-ids ./cmd/unique-ids

build-broadcast:
	go build -o build/bin/broadcast ./cmd/broadcast

build-broadcast-d:
	go build -o build/bin/broadcast-d ./cmd/broadcast-d

build-broadcast-e:
	go build -o build/bin/broadcast-e ./cmd/broadcast-e

build-counter:
	go build -o build/bin/maelstrom-counter ./cmd/counter

build-kafka:
	go build -o build/bin/maelstrom-kafka ./cmd/kafka

build-kafka-b:
	go build -o build/bin/maelstrom-kafka-b ./cmd/kafka-b

build-kafka-c:
	go build -o build/bin/maelstrom-kafka-c ./cmd/kafka-c

build-txn:
	go build -o build/bin/maelstrom-txn ./cmd/txn

build-txn-b:
	go build -o build/bin/maelstrom-txn-b ./cmd/txn-b

build-txn-c:
	go build -o build/bin/maelstrom-txn-c ./cmd/txn-c

test-echo:
	@cd maelstrom; ./maelstrom test -w echo --bin ../build/bin/echo --node-count 1 --time-limit 10

test-unique-ids:
	@cd maelstrom; ./maelstrom test -w unique-ids --bin ../build/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

test-broadcast:
	@cd maelstrom; ./maelstrom test -w broadcast --bin ../build/bin/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

test-broadcast-d0:
	@cd maelstrom; ./maelstrom test -w broadcast --bin ../build/bin/broadcast-d --node-count 25 --time-limit 20 --rate 100 --latency 100

test-broadcast-d1:
	@cd maelstrom; ./maelstrom test -w broadcast --bin ../build/bin/broadcast-d --node-count 25 --time-limit 20 --rate 100 --latency 100 --nemesis partition

test-broadcast-e0:
	@cd maelstrom; ./maelstrom test -w broadcast --bin ../build/bin/broadcast-e --node-count 25 --time-limit 20 --rate 100 --latency 100

test-broadcast-e1:
	@cd maelstrom; ./maelstrom test -w broadcast --bin ../build/bin/broadcast-e --node-count 25 --time-limit 20 --rate 100 --latency 100 --nemesis partition

test-counter:
	@cd maelstrom; ./maelstrom test -w g-counter --bin ../build/bin/maelstrom-counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition

test-kafka:
	@cd maelstrom; ./maelstrom test -w kafka --bin ../build/bin/maelstrom-kafka --node-count 1 --concurrency 2n --time-limit 20 --rate 1000

test-kafka-b:
	@cd maelstrom; ./maelstrom test -w kafka --bin ../build/bin/maelstrom-kafka-b --node-count 2 --concurrency 2n --time-limit 20 --rate 1000

test-kafka-c:
	@cd maelstrom; ./maelstrom test -w kafka --bin ../build/bin/maelstrom-kafka-c --node-count 2 --concurrency 2n --time-limit 20 --rate 1000

test-txn:
	@cd maelstrom; ./maelstrom test -w txn-rw-register --bin ../build/bin/maelstrom-txn --node-count 1 --time-limit 20 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total

test-txn-b:
	@cd maelstrom; ./maelstrom test -w txn-rw-register --bin ../build/bin/maelstrom-txn-b --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted

test-txn-b-1:
	@cd maelstrom; ./maelstrom test -w txn-rw-register --bin ../build/bin/maelstrom-txn-b --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition

test-txn-c:
	@cd maelstrom; ./maelstrom test -w txn-rw-register --bin ../build/bin/maelstrom-txn-c --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-committed --availability total --nemesis partition
