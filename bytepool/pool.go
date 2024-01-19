package bytepool

import (
	"container/list"
	"sync"
)

type Pool struct {
	recycleList *list.List
	mu          sync.Mutex

	pageSize        int
	maxRecycleCount int
}

func New(pageSize int, maxRecycleCount int) *Pool {
	return &Pool{
		recycleList:     list.New(),
		pageSize:        pageSize,
		maxRecycleCount: maxRecycleCount,
	}
}

func (p *Pool) NewArea(size int) *Area {
	pageCount := size % p.pageSize
	if p.pageSize*pageCount < size {
		pageCount++
	}

	a := &Area{
		pool:      p,
		readIndex: 0,
		writable:  true,
		data: &DataArea{
			cap:      size,
			attached: 1,
			pages:    make([]Page, pageCount),
		},
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range a.data.pages {
		e := p.recycleList.Front()
		if e != nil {
			a.data.pages[i] = e.Value.([]byte)
			p.recycleList.Remove(e)
		} else {
			a.data.pages[i] = make(Page, p.pageSize)
		}
	}

	return a
}

func (p *Pool) Recycle(a *Area) {
	pages := a.data.pages
	a.data.pages = nil
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, page := range pages {
		if p.recycleList.Len() >= p.maxRecycleCount {
			continue
		}
		p.recycleList.PushBack(page)
	}
}
