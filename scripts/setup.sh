#!/bin/bash

# Create necessary directories
mkdir -p tmp/logs
mkdir -p tmp/cache

# Install development tools
go install golang.org/x/tools/cmd/godoc@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Download dependencies
go mod download

# Build the framework
go build -v ./...

echo "Development environment setup complete!" 