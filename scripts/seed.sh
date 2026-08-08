#!/usr/bin/env bash
# Seed the database with demo data.
set -e
cd "$(dirname "$0")/.."
go run ./cmd/seed
