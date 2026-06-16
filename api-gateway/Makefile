lint:
	golangci-lint run ./...

lint-fix:
	goimports -w .
	gofmt -s -w .
	golangci-lint run --fix ./...