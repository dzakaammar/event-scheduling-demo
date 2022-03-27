SHELL:=/bin/bash

ROOT_DIRECTORY=.
DIRS := ${sort ${dir ${wildcard ${ROOT_DIRECTORY}/*/*/*/}}}

deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest \

lint:
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --enable=goimports  --enable=unconvert --enable=unparam --enable=gosec --timeout=2m

cc:
	for dir in $(DIRS) ; do \
		gocyclo -over 10 -ignore "pb|vendor/|mock/" $${dir} ; \
	done
	
test:
	go test -race -cover ./...

test.complete: lint cc test

run:
	go run main.go

generate:
	go generate ./...