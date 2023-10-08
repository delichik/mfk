package objectpool

import (
	"container/list"
	"time"
)

type poolNode[T any] struct {
	inQueueTime uint32
	object      *T
}

type Pool[T any] struct {
	baseTime int64
	max      int

	l       list.List
	creator func() *T
}

func New[T any](max int, creator func() *T) *Pool[T] {
	return &Pool[T]{
		baseTime: time.Now().Unix(),
		max:      max,
		l:        list.List{},
		creator:  creator,
	}
}

func (p *Pool[T]) Get() *T {
	_ = p.removeFirstExpired()
	if p.l.Len() == 0 {
		if p.creator != nil {
			return p.creator()
		} else {
			return new(T)
		}
	}
	node := p.l.Remove(p.l.Back()).(poolNode[T])
	return node.object
}

func (p *Pool[T]) Put(t *T) {
	_ = p.removeFirstExpired()
	cleanObject(t)
	node := poolNode[T]{
		inQueueTime: uint32(time.Now().Unix() - p.baseTime),
		object:      t,
	}
	p.l.PushBack(node)
}

func (p *Pool[T]) removeFirstExpired() bool {
	if p.l.Len() == 0 {
		return false
	}

	node := p.l.Front()
	if node.Value.(poolNode[T]).inQueueTime > uint32(time.Now().Unix()-p.baseTime) {
		return false
	}

	p.l.Remove(node)
	return true
}
