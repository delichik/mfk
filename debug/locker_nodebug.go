//go:build !debuggable || !locker_track

package debug

import (
	"sync"
)

type Mutex = sync.Mutex
type RWMutex = sync.RWMutex
