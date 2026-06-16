.PHONY: build run tidy lint

build:
	go build -o bin/auth-service ./cmd

run:
	go run ./cmd

tidy:
	GONOSUMDB="github.com/zimnyles/*" go mod tidy

lint:
	golangci-lint run ./...
