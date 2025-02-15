package loadtest

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// todo: we can import k6 as a library since it's written in Go

type Results struct {
	RequestCount   int64
	ResponseTimes  []time.Duration
	ErrorCount     int64
	AverageLatency time.Duration
	MaxLatency     time.Duration
	MinLatency     time.Duration
}

func Run(concurrency int, duration time.Duration) *Results {
	results := &Results{
		ResponseTimes: make([]time.Duration, 0, 1000),
		MinLatency:    time.Hour,
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Optimize transport settings
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
		DisableKeepAlives:   false,
		ForceAttemptHTTP2:   true,
		MaxConnsPerHost:     0, // No limit
		WriteBufferSize:     64 * 1024,
		ReadBufferSize:      64 * 1024,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Second, // Shorter timeout
	}

	// Pre-warm connections
	for i := 0; i < 10; i++ {
		resp, _ := client.Get("http://localhost:8080/")
		if resp != nil {
			resp.Body.Close()
		}
	}

	deadline := time.Now().Add(duration)

	// Launch workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localTimes := make([]time.Duration, 0, 1000)

			for time.Now().Before(deadline) {
				reqStart := time.Now()
				resp, err := client.Get("http://localhost:8080/")
				latency := time.Since(reqStart)

				atomic.AddInt64(&results.RequestCount, 1)

				if err != nil {
					atomic.AddInt64(&results.ErrorCount, 1)
					continue
				}
				resp.Body.Close()

				localTimes = append(localTimes, latency)
			}

			mu.Lock()
			results.ResponseTimes = append(results.ResponseTimes, localTimes...)
			mu.Unlock()
		}()
	}

	wg.Wait()
	transport.CloseIdleConnections()

	// Calculate statistics
	if len(results.ResponseTimes) > 0 {
		results.MinLatency = results.ResponseTimes[0]
		results.MaxLatency = results.ResponseTimes[0]
		var totalLatency time.Duration

		for _, lat := range results.ResponseTimes {
			totalLatency += lat
			if lat < results.MinLatency {
				results.MinLatency = lat
			}
			if lat > results.MaxLatency {
				results.MaxLatency = lat
			}
		}
		results.AverageLatency = totalLatency / time.Duration(len(results.ResponseTimes))
	}

	return results
}
