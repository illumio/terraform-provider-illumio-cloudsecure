#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "[pre-commit] Missing required tool: $1" >&2
    exit 1
  fi
}

need go
need golangci-lint

echo "[pre-commit] Cleaning generated files..."
make clean

echo "[pre-commit] Generating code..."
make generate

echo "[pre-commit] Running golangci-lint..."
# Use local config
golangci-lint run

echo "[pre-commit] Running attrcheck..."
go run ./tools/attrcheck ./api/schema

echo "[pre-commit] SUCCESS"

