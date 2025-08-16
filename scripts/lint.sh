#!/bin/bash

set -e

modules=$(find . -name "go.mod" -exec dirname {} \;)

for mod in $modules; do
  echo "Linting $mod"
  
  pushd "$mod"
  golangci-lint run --fix --config ../.golangci.yaml
  popd
done
