package debouncer

import (
	"sync"
	"time"
)

// Debouncer uses an in-memory map to debounce signals per component_id.
type Debouncer struct {
	mu     sync.RWMutex
	store  map[string]time.Time
	window time.Duration
}

func New(window time.Duration) *Debouncer {
	d := &Debouncer{
		store:  make(map[string]time.Time),
		window: window,
	}
	
	// Start a cleanup goroutine to prevent memory leaks from the map
	go d.cleanup()
	return d
}

// Allow returns true if the signal for the component should be processed (it's the first in the window).
// It returns false if it should be debounced.
func (d *Debouncer) Allow(componentID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	lastSeen, exists := d.store[componentID]
	now := time.Now()

	if !exists || now.Sub(lastSeen) > d.window {
		d.store[componentID] = now
		return true
	}

	return false
}

// cleanup periodically removes expired keys to free up memory.
func (d *Debouncer) cleanup() {
	ticker := time.NewTicker(d.window * 2)
	for range ticker.C {
		now := time.Now()
		d.mu.Lock()
		for k, v := range d.store {
			if now.Sub(v) > d.window {
				delete(d.store, k)
			}
		}
		d.mu.Unlock()
	}
}
