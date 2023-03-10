package queue

import (
	"container/list"
	"sync"
	"time"
)

type Queue struct {
	element   *list.List
	cond      *sync.Cond
	closed    bool
	maxLength int
}

func New(maxLength int) *Queue {
	return &Queue{
		element:   list.New(),
		cond:      sync.NewCond(&sync.Mutex{}),
		closed:    false,
		maxLength: maxLength,
	}
}

func (q *Queue) Close() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.closed = true
	q.cond.Broadcast()
}

func (q *Queue) TryEnqueue(data interface{}) error {
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

func (q *Queue) Enqueue(data interface{}) error {
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

func (q *Queue) InsertFront(data interface{}) error {
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

func (q *Queue) Dequeue() (interface{}, error) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for q.element.Len() == 0 && !q.closed {
		// 阻塞
		q.cond.Wait()
	}

	if q.closed {
		return nil, ErrClosed
	}

	e := q.element.Front()
	q.element.Remove(e)
	return e.Value, nil
}

func (q *Queue) WaitEmpty() {
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
