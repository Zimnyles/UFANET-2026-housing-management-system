.PHONY: build run tidy fmt lint lint-fix test

build:
	go build -o bin/api-gateway ./cmd/api

run:
	go run ./cmd/api

tidy:
	GONOSUMDB="github.com/zimnyles/*" go mod tidy

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

lint-fix:
	goimports -w .
	gofmt -s -w .
	golangci-lint run --fix ./...

test:
	go test ./...