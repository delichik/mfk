//go:build !debug

package tracked

import (
	"sync"
)

type RWMutex = sync.RWMutex
