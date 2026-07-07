.PHONY: build run tidy fmt lint test

build:
	go build -o bin/auth-service ./cmd

run:
	go run ./cmd

tidy:
	GONOSUMDB="github.com/zimnyles/*" go mod tidy

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

test:
	go test ./...
