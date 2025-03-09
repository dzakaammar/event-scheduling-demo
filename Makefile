SHELL := /bin/bash

# Test configuration
TEST_TARGET_DIR := $$(go list ./... | grep -v "/gen/" | grep -v "/mock/" )
COVER_PROFILE_FILE ?= cover.out

.PHONY: help deps install-tools lint vuln test integration-test check-coverage test-complete run generate proto

# Default target when just running 'make'
.DEFAULT_GOAL := help

## Show this help message
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  \x1b[32;01m%-20s\x1b[0m %s\n", helpCommand, helpMessage; \
		} \
	} { lastLine = $$0 }' $(MAKEFILE_LIST)
	@printf "\n"

## Install tools
install-tools:
	@go install -modfile=go.tool.mod tool

## Run code linter
lint:
	@go tool -modfile=go.tool.mod golangci-lint run --config=.golangci.yml

## Run vulnerability check
vuln:
	@go tool -modfile=go.tool.mod govulncheck ./...

## Run unit tests with coverage
test:
	go test -race -v -count=1 -cover -coverprofile $(COVER_PROFILE_FILE) $(TEST_ARGS) $(TEST_TARGET_DIR)

## Run integration tests
integration-test:
	@go tool -modfile=go.tool.mod ginkgo --label-filter="integration" -r -p --succinct --race --trace --fail-fast

## Check if test coverage meets 90% threshold
check-coverage:
	go tool cover -func $(COVER_PROFILE_FILE) | grep total | awk '{print substr($$3, 1, length($$3)-1)}' | awk '{if ($$0 < 90) exit 1; else exit 0}'

## Run lint, unit tests, and integration tests
test-complete: lint test integration-test

## Run the application
run:
	go run main.go

## Generate code based on go:generate directives
generate:
	go generate ./...

## Generate code from protocol buffers
proto:
	@go tool -modfile=go.tool.mod buf generate
