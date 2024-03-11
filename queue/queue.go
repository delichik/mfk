package queue

import (
	"sync"
	"time"

	"github.com/delichik/daf/wrapper"
)

type Queue[T any] struct {
	element   *wrapper.List[T]
	cond      *sync.Cond
	closed    bool
	maxLength int
}

func New[T any](maxLength int) *Queue[T] {
	return &Queue[T]{
		element:   wrapper.New[T](),
		cond:      sync.NewCond(&sync.Mutex{}),
		closed:    false,
		maxLength: maxLength,
	}
}

func (q *Queue[T]) Close() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.closed = true
	q.cond.Broadcast()
}

func (q *Queue[T]) TryEnqueue(data T) error {
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()
	}()

	if q.closed {
		return ErrClosed
	}

	if q.maxLength != 0 && q.element.Len() >= q.maxLength {
		return ErrFulled
	}

	q.element.PushBack(data)
	q.cond.Signal()

	return nil
}

func (q *Queue[T]) Enqueue(data T) error {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for !q.closed && q.maxLength != 0 && q.element.Len() >= q.maxLength {
		q.cond.Wait()
	}

	if q.closed {
		return ErrClosed
	}

	q.element.PushBack(data)
	q.cond.Signal()

	return nil
}

func (q *Queue[T]) InsertFront(data T) error {
	if data == nil {
		return nil
	}

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.closed {
		return ErrClosed
	}

	q.element.PushFront(data)
	q.cond.Signal()
	return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for q.element.Len() == 0 && !q.closed {
		// 阻塞
		q.cond.Wait()
	}

	if q.closed {
		return *(new(T)), ErrClosed
	}

	e := q.element.Front()
	q.element.Remove(e)
	return e.Value, nil
}

func (q *Queue[T]) WaitEmpty() {
	for {
		q.cond.L.Lock()
		if q.element.Len() == 0 || q.closed {
			q.cond.L.Unlock()
			return
		}
		q.cond.L.Unlock()
		time.Sleep(time.Millisecond)
	}
}
