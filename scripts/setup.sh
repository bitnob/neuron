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

# Kernel parameters for high performance
cat >> /etc/sysctl.conf << EOF
# Network performance tuning
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 65535
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 10
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_max_tw_buckets = 2000000
net.ipv4.tcp_fastopen = 3
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 87380 16777216

# File descriptor limits
fs.file-max = 2097152
fs.nr_open = 2097152

# VM settings
vm.swappiness = 10
vm.dirty_ratio = 60
vm.dirty_background_ratio = 2

# Performance settings
kernel.sched_min_granularity_ns = 10000000
kernel.sched_wakeup_granularity_ns = 15000000
EOF

# Apply changes
sysctl -p

# Increase file descriptor limits
cat >> /etc/security/limits.conf << EOF
* soft nofile 1048576
* hard nofile 1048576
EOF

# Set Go GC parameters
export GOGC=100
export GOMEMLIMIT=4GiB 