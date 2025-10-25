package tracker

import (
	"sync"
	"time"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
}

// IPTracker tracks unique IP addresses with thread-safe operations
type IPTracker struct {
	mu        sync.RWMutex
	uniqueIPs map[string]time.Time
	logger    Logger
}

// NewIPTracker creates a new IP tracker
func NewIPTracker(logger Logger) *IPTracker {
	return &IPTracker{
		uniqueIPs: make(map[string]time.Time),
		logger:    logger,
	}
}

// Add adds an IP address to the tracker with timestamp
func (t *IPTracker) Add(ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.uniqueIPs[ip]; !exists {
		t.logger.Debug("new unique IP tracked", "ip", ip)
	}

	t.uniqueIPs[ip] = time.Now()
}

// Count returns the number of unique IP addresses
func (t *IPTracker) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.uniqueIPs)
}

// GetIPs returns a copy of all tracked IPs
func (t *IPTracker) GetIPs() map[string]time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ips := make(map[string]time.Time, len(t.uniqueIPs))
	for ip, ts := range t.uniqueIPs {
		ips[ip] = ts
	}
	return ips
}

// Clear removes all tracked IPs (useful for testing)
func (t *IPTracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.uniqueIPs = make(map[string]time.Time)
	t.logger.Info("tracker cleared")
}
