#!/bin/bash
set -e

if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

go run main.go 