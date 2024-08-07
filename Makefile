SHELL:=/bin/bash


TEST_TARGET_DIR = $$(go list ./... | grep -v "/gen/" | grep -v "/mock/" )
COVER_PROFILE_FILE ?= cover.out

## to install all dependencies
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
	go install golang.org/x/vuln/cmd/govulncheck@latest && \
	go install github.com/golang/mock/mockgen@v1.6.0 && \
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

## to install tools from tools.go
install-tools:
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

## to run linter
lint:
	golangci-lint run --config=.golangci.yml

## to run vulnarability check
vuln:
	govulncheck ./...
	
## to run all unit test with coverage
test:
	go test -race -v -count=1 -cover -coverprofile $(COVER_PROFILE_FILE) $(TEST_ARGS) $(TEST_TARGET_DIR)

## to run integration test with docker container
integration.test:
	go run github.com/onsi/ginkgo/v2/ginkgo --label-filter="integration" -r -p --succinct --race --trace --fail-fast
	
## to check the total of test coverage 
check.coverage:
	go tool cover -func $(COVER_PROFILE_FILE) | grep total | awk '{print substr($$3, 1, length($$3)-1)}' | awk '{if ($$0 <= 90) exit 1 ; else exit 0 }'

## to run lint, unit test, and integration test sequentially
test.complete: lint test integration.test

## to run the application
run:
	go run main.go

## to generate defined task
generate:
	go generate ./...

## to generate protobuf
proto:
	@buf generate

help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
		helpCommand = $$1; \
		helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
	gsub("\\\\", "", helpCommand); \
	gsub(":+$$", "", helpCommand); \
		printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"