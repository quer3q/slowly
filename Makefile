.PHONY: all
all: build test

.PHONY: build
build:
	go build -o ./bin/slowly ./cmd/main

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	go run ./cmd/main