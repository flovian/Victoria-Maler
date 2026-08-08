#!/usr/bin/env bash
# Run the Victoria Maler server.
set -e
cd "$(dirname "$0")/.."
go run ./cmd/server
