SHELL:=/bin/bash

deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run --config=.golangci.yml
	
test:
ifdef COVER_PROFILE_FILE
	$(eval TEST_ARGS := -coverprofile $(COVER_PROFILE_FILE))
endif
	go test -race -cover $(TEST_ARGS) ./...

check.coverage:
	go tool cover -func $(COVER_PROFILE_FILE) | grep total | awk '{print substr($$3, 1, length($$3)-1)}' | awk '{if ($$0 <= 90) exit 1 ; else exit 0 }'

test.complete: lint test

run:
	go run main.go

generate:
	go generate ./...