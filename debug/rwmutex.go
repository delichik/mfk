//go:build debugable

package debug

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/delichik/mfk/logger"
)

type RWMutex struct {
	sync.RWMutex
	AlertTimeout time.Duration
}

func (m *RWMutex) RLock() {
	start := time.Now()
	m.RWMutex.RLock()
	if m.AlertTimeout > 0 && start.Sub(time.Now()) > m.AlertTimeout {
		logger.Warn("RLock() takes a long time to finish", zap.StackSkip("stack", 2))
	}
}

func (m *RWMutex) RUnlock() {
	start := time.Now()
	m.RWMutex.RUnlock()
	if m.AlertTimeout > 0 && start.Sub(time.Now()) > m.AlertTimeout {
		logger.Warn("RUnlock() takes a long time to finish", zap.StackSkip("stack", 2))
	}
}

func (m *RWMutex) Lock() {
	start := time.Now()
	m.RWMutex.Lock()
	if m.AlertTimeout > 0 && start.Sub(time.Now()) > m.AlertTimeout {
		logger.Warn("Lock() takes a long time to finish", zap.StackSkip("stack", 2))
	}
}

func (m *RWMutex) Unlock() {
	start := time.Now()
	m.RWMutex.Unlock()
	if m.AlertTimeout > 0 && start.Sub(time.Now()) > m.AlertTimeout {
		logger.Warn("Unlock() takes a long time to finish", zap.StackSkip("stack", 2))
	}
}
