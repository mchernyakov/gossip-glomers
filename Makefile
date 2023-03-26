SHELL=/bin/bash -o pipefail

build-echo:
	go build -o build/bin/echo cmd/echo.go

build-unique-ids:
	go build -o build/bin/unique-ids cmd/unique_ids.go

build-broadcast:
	go build -o build/bin/broadcast cmd/broadcast.go

build-broadcast-d:
	go build -o build/bin/broadcast-d cmd/broadcast_d.go

build-broadcast-e:
	go build -o build/bin/broadcast-e cmd/broadcast_e.go

build-counter:
	go build -o build/bin/maelstrom-counter cmd/counter.go

build-kafka:
	go build -o build/bin/maelstrom-kafka cmd/kafka.go

build-kafka-b:
	go build -o build/bin/maelstrom-kafka-b cmd/kafka_b.go

build-kafka-c:
	go build -o build/bin/maelstrom-kafka-c cmd/kafka_c.go

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
