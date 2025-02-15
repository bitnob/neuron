.PHONY: build test run clean

# Build the framework
build:
	go build -v ./...

# Run tests
test:
	go test -v ./...

# Run example application
run:
	go run example/basic/main.go

# Clean build artifacts
clean:
	go clean
	rm -f neuron

# Install dependencies
deps:
	go mod download

# Run linter
lint:
	golangci-lint run

# Generate documentation
docs:
	godoc -http=:6060

# Build and run all tests
all: deps build test 