#!/bin/bash
set -e

echo "Checking if go.mod and go.sum are tidy..."
go mod tidy

if ! git diff --exit-code; then
  echo "ERROR: go.mod or go.sum is not tidy. Run 'go mod tidy' locally."
  exit 1
else
  echo "SUCCESS: Dependencies are tidy."
fi
