package observability

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	// signalCount keeps track of how many signals were ingested.
	signalCount uint64
)

// RecordSignal increments the ingestion counter.
func RecordSignal() {
	atomic.AddUint64(&signalCount, 1)
}

// StartLogger spawns a background goroutine that logs ingestion throughput every interval.
func StartLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		var lastCount uint64

		for range ticker.C {
			currentCount := atomic.LoadUint64(&signalCount)
			delta := currentCount - lastCount
			lastCount = currentCount
			
			rate := float64(delta) / interval.Seconds()
			fmt.Printf("[Observability] Ingestion Throughput: %.2f Signals/sec\n", rate)
		}
	}()
}
