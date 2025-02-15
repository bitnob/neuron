package main

import (
	"fmt"
	"neuron/pkg/loadtest"
	"time"
)

func main() {
	results := loadtest.Run(100, 30*time.Second)
	fmt.Printf("Load Test Results:\n")
	fmt.Printf("Requests: %d\n", results.RequestCount)
	fmt.Printf("Errors: %d\n", results.ErrorCount)
	fmt.Printf("Average Latency: %v\n", results.AverageLatency)
	fmt.Printf("Max Latency: %v\n", results.MaxLatency)
	fmt.Printf("Min Latency: %v\n", results.MinLatency)
}
