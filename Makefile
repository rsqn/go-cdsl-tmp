.PHONY: test build clean

# Default target
all: build test

# Build the project
build:
	go build -o go-cdsl -v ./cmd/example

# Run all tests
test:
	go test -v ./pkg/...

# Clean build artifacts
clean:
	rm -f go-cdsl
	go clean
