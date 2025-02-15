#!/bin/bash

# Clean up
rm -f go.sum
go clean -modcache

# Get dependencies
go get -u golang.org/x/time/rate
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/go-redis/redis/v8
go get -u gopkg.in/yaml.v3
go get -u github.com/golang-jwt/jwt/v4

# Tidy up modules
go mod tidy

echo "Dependencies updated successfully!" 