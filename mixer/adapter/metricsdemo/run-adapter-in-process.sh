#!/bin/bash

set -e

cd $(dirname $0)
go run cmd/main.go 53175
