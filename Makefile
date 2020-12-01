.PHONY: build
build:
	go build -o bin/vin cmd/vin/main.go

.PHONY: format
format:
	goimports -w .

.PHONY: lint
lint:
	golangci-lint run --tests ./...

.PHONY: test
test:
	go test -v ./...
