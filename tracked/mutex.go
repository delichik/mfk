//go:build !debug

package tracked

import (
	"sync"
)

type Mutex = sync.Mutex
