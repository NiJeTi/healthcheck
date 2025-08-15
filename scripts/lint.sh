#!/bin/sh

set -e

GOLANGCI_LINT_IMAGE=$1

modules=$(find . -name "go.mod" -exec dirname {} \;)

for mod in $modules; do
  echo "Linting $mod"
  
  docker run -t --rm \
    -v "$mod":/src -v ./.golangci.yaml:/src/.golangci.yaml\
    -w /src "$GOLANGCI_LINT_IMAGE" golangci-lint run --fix
done
