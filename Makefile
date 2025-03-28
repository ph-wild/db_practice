BINARY=app

.DEFAULT_GOAL := build

build:
	go build -o $(BINARY) ./cmd

run: build
	./$(BINARY)

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

fmt:
	go fmt ./...

tidy:
	go mod tidy

all: fmt tidy lint test build

help:
	@echo "Makefile for Go project:"
	@echo "  make build   - Build the binary"
	@echo "  make run     - Run the application"
	@echo "  make test    - Run tests"
	@echo "  make lint    - Run linter"
	@echo "  make clean   - Remove built files"
	@echo "  make fmt     - Format the code"
	@echo "  make tidy    - Update dependencies"
	@echo "  make all     - Run all steps (fmt, tidy, lint, test, build)"
