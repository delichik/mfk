package workpool

import (
	"context"
	"sync"
	"time"
)

type WorkPool struct {
	count   int
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	channel chan func()
}

func New(count int) *WorkPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkPool{
		ctx:     ctx,
		cancel:  cancel,
		count:   count,
		channel: make(chan func(), 1),
	}
}

func (p *WorkPool) Add(call func()) {
	select {
	case p.channel <- call:
	case <-p.ctx.Done():
	}
}

func (p *WorkPool) AddWithTimeout(call func(), timeout time.Duration) {
	ctx, cancel := context.WithTimeout(p.ctx, timeout)
	defer cancel()
	select {
	case p.channel <- call:
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			call()
		}
	}
}

func (p *WorkPool) Start() {
	for i := 0; i < p.count; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for p.ctx.Err() == nil {
				select {
				case call := <-p.channel:
					call()
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *WorkPool) Stop() {
	p.cancel()
	p.wg.Wait()
}
