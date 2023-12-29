//go:build debug

package tracked

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"vap/pkg/logger"
)

type RWMutex struct {
	sync.RWMutex
	AlertTimeout time.Duration
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
