package wrapper

import (
	"sync"

	"github.com/google/go-cmp/cmp"
)

type SyncMap[K cmp.MapIndex, V any] struct {
	m sync.Map
}

func (m *SyncMap[K, V]) Load(key K) (value V, loaded bool) {
	v, loaded := m.m.Load(key)
	if loaded {
		return v.(V), true
	} else {
		return *(new(V)), false
	}
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual any, loaded bool) {
	return m.m.LoadOrStore(key, value)
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.m.LoadAndDelete(key)
	if loaded {
		return v.(V), true
	} else {
		return *(new(V)), false
	}
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *SyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	p, loaded := m.m.Swap(key, value)
	if loaded {
		return p.(V), true
	} else {
		return *(new(V)), false
	}
}

func (m *SyncMap[K, V]) CompareAndSwap(key K, old V, new V) bool {
	return m.m.CompareAndSwap(key, old, new)
}

func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
