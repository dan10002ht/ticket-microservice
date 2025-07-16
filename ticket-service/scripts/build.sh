#!/bin/bash
set -e

go mod tidy
go build -o ticket-service main.go 