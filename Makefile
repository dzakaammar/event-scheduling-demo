SHELL:=/bin/bash

lint:
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --enable=goimports  --enable=unconvert --enable=unparam --enable=gosec --timeout=2m

test:
	go test -race -cover ./...

run:
	go run main.go

generate:
	go generate ./...