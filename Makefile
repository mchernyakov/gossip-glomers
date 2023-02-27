SHELL=/bin/bash -o pipefail

build-echo:
	go build -o build/bin/echo cmd/echo.go

build-unique-ids:
	go build -o build/bin/unique-ids cmd/unique_ids.go

build-broadcast:
	go build -o build/bin/broadcast cmd/broadcast.go

test-echo:
	@cd maelstrom; ./maelstrom test -w echo --bin ../build/bin/echo --node-count 1 --time-limit 10

test-unique-ids:
	@cd maelstrom; ./maelstrom test -w unique-ids --bin ../build/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

test-broadcast:
	@cd maelstrom; ./maelstrom test -w broadcast --bin ../build/bin/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition
