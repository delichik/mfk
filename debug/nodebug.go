//go:build !debugable

package debug

import (
	"sync"
)

type Mutex = sync.Mutex
type RWMutex = sync.RWMutex
