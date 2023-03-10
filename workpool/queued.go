package workpool

import (
	"context"
	"hash/fnv"
	"runtime"
	"sync"
	"time"
)

type QueuedWorkPool struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	chans     []chan func()
	chanLen   int
	chanCount int
}

func NewQueuedWorkPool(chanLen int) *QueuedWorkPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &QueuedWorkPool{
		ctx:       ctx,
		cancel:    cancel,
		chanLen:   chanLen,
		chanCount: runtime.NumCPU() * 2,
		chans:     make([]chan func(), runtime.NumCPU()*2),
	}
}

func (p *QueuedWorkPool) Add(tag string, call func()) {
	if tag == "" {
		call()
		return
	}
	v := 0
	h := fnv.New32a()
	_, err := h.Write([]byte(tag))
	if err == nil {
		v = int(h.Sum32()) % p.chanCount
	}

	select {
	case p.chans[v] <- call:
	case <-p.ctx.Done():
	}
}

func (p *QueuedWorkPool) AddWithTimeout(tag string, call func(), timeout time.Duration) {
	v := 0
	h := fnv.New32a()
	_, err := h.Write([]byte(tag))
	if err == nil {
		v = int(h.Sum32()) % p.chanCount
	}

	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()
	select {
	case p.chans[v] <- call:
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			call()
		}
	}
}

func (p *QueuedWorkPool) Start() {
	for i := 0; i < p.chanCount; i++ {
		p.chans[i] = make(chan func(), p.chanLen)
		p.wg.Add(1)
		go func(i int) {
			defer p.wg.Done()
			for p.ctx.Err() == nil {
				select {
				case call := <-p.chans[i]:
					call()
				case <-p.ctx.Done():
					return
				}
			}
		}(i)
	}
}

func (p *QueuedWorkPool) Stop() {
	p.cancel()
	p.wg.Wait()
}
