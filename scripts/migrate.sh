#!/usr/bin/env bash
# Apply database migrations.
set -e
cd "$(dirname "$0")/.."
go run ./cmd/migrate
