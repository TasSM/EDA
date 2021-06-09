.PHONY: test nats-run nats-build
.DEFAULT_GOAL := build-linux

test:
	go test -race ./...

nats-run:
	go run -race ./cmd/nats/main.go

build-osx:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/nats/main.darwin ./cmd/nats/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/nats/main.linux ./cmd/nats/main.go