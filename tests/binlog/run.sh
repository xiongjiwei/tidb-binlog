#!/bin/sh

set -e

cd "$(dirname "$0")"

run_drainer &

go build -o out

./out -config ./config.toml > ${OUT_DIR-/tmp}/$TEST_NAME.out 2>&1

killall drainer