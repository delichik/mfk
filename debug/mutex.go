//go:build debuggable && locker_track

package debug

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/delichik/mfk/logger"
)

type Mutex struct {
	sync.Mutex
	AlertTimeout time.Duration
}

func (m *Mutex) Lock() {
	start := time.Now()
	m.Mutex.Lock()
	if m.AlertTimeout > 0 && start.Sub(time.Now()) > m.AlertTimeout {
		logger.Warn("Lock() takes a long time to finish", zap.StackSkip("stack", 2))
	}
}

func (m *Mutex) Unlock() {
	start := time.Now()
	m.Mutex.Unlock()
	if m.AlertTimeout > 0 && start.Sub(time.Now()) > m.AlertTimeout {
		logger.Warn("Unlock() takes a long time to finish", zap.StackSkip("stack", 2))
	}
}
