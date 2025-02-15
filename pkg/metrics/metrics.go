package metrics

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	RequestCount     uint64
	ResponseTime     uint64
	ErrorCount       uint64
	ActiveGoroutines uint64
	MemoryUsage      uint64
	StartupTime      time.Duration
}

func (m *Metrics) TrackRequest(duration time.Duration) {
	// todo: atomic will lose state on restarts, great for syncing across
	// goroutines, prometheus >>>>>>>>
	atomic.AddUint64(&m.RequestCount, 1)
	atomic.AddUint64(&m.ResponseTime, uint64(duration))
}

func (m *Metrics) TrackError() {
	atomic.AddUint64(&m.ErrorCount, 1)
}
