package safe

import (
	"errors"
	"sync"
)

const SignalDefaultOwner = "any"

type Signal struct {
	cond        *sync.Cond
	paused      bool
	activeOwner string
	closed      bool
}

func NewSignal() *Signal {
	return &Signal{
		cond:   sync.NewCond(&sync.Mutex{}),
		paused: true,
	}
}

func (p *Signal) IsPaused() bool {
	return p.paused
}

func (p *Signal) Close() {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	p.closed = true
	p.cond.Broadcast()
}

func (p *Signal) Pause(owner string) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	if owner != p.activeOwner {
		return
	}
	p.paused = true
	p.activeOwner = ""
}

func (p *Signal) Resume(owner string) {
	if owner == "" {
		return
	}
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	p.paused = false
	p.activeOwner = owner
	p.cond.Broadcast()
}

func (p *Signal) Wait() error {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	for p.paused && !p.closed {
		p.cond.Wait()
	}
	if p.closed {
		return errors.New("closed")
	}
	return nil
}
