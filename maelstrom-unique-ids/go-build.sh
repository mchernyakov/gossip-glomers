#!/usr/bin/env bash

PATH_TO_BUILD="build/bin/"
ENTRY_POINT="cmd/main.go"
OUTPUT_PREFIX=${PATH_TO_BUILD}"maelstrom-unique-ids"

go build -o ${OUTPUT_PREFIX} ${ENTRY_POINT}
