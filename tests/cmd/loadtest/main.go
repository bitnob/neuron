package main

// Before running high-load tests:
// 1. Tune system limits (as root):
//    sysctl -w net.ipv4.tcp_fin_timeout=30
//    sysctl -w net.ipv4.tcp_tw_reuse=1
//    sysctl -w net.core.somaxconn=65535
//    ulimit -n 65535
//
// 2. Run server with optimized settings:
//    GOMAXPROCS=8 GOGC=800 go run cmd/server/main.go
//
// 3. Gradually increase load:
//    - Start with current values
//    - If stable, try concurrency=2000, requests=500000
//    - Finally try concurrency=5000, requests=1000000

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	concurrency = 100   // Increased concurrency
	requests    = 10000 // Increased request count
	url         = "http://localhost:8080/"
)

func main() {
	// Run ab command
	cmd := exec.Command("ab",
		"-k",       // Use HTTP KeepAlive
		"-r",       // Don't exit on socket receive errors
		"-l",       // Don't report errors for response length differences
		"-s", "30", // Timeout after 30 seconds
		"-g", "results.tsv", // Save detailed results
		"-n", strconv.Itoa(requests),
		"-c", strconv.Itoa(concurrency),
		"-H", "Accept-Encoding: gzip, deflate",
		"-H", "Connection: keep-alive",
		url,
	)

	// Get output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Warning: ab completed with error: %v\n", err)
	}

	// Parse and display results even if there were some errors
	results := string(output)
	fmt.Println("\nLoad Test Results:")
	fmt.Println("==================")
	fmt.Printf("Test Configuration:\n")
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Number of requests: %d\n\n", requests)

	// Extract key metrics
	metrics := map[string]string{
		"Complete requests":      extractMetric(results, "Complete requests:"),
		"Failed requests":        extractMetric(results, "Failed requests:"),
		"Requests per second":    extractMetric(results, "Requests per second:"),
		"Time per request":       extractMetric(results, "Time per request:"),
		"Transfer rate":          extractMetric(results, "Transfer rate:"),
		"Percentage of requests": extractPercentiles(results),
	}

	// Print metrics
	for k, v := range metrics {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func extractMetric(output, metric string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, metric) {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return "N/A"
}

func extractPercentiles(output string) string {
	var percentiles strings.Builder
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "50%") ||
			strings.Contains(line, "90%") ||
			strings.Contains(line, "99%") {
			percentiles.WriteString("\n    " + strings.TrimSpace(line))
		}
	}
	return percentiles.String()
}
