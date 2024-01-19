package bytepool

import (
	"fmt"
	"sync"
)

type DataArea struct {
	cap int

	recycleLocker sync.Mutex
	attached      int

	dataLocker sync.RWMutex
	writeIndex int
	pages      []Page
}

func (a *DataArea) read(p []byte, startIndex int) (int, error) {
	a.dataLocker.RLock()
	defer a.dataLocker.RUnlock()

	if a.pages == nil {
		panic(ErrClosed)
	}

	if startIndex >= a.writeIndex {
		return 0, nil
	}

	if len(a.pages) == 0 {
		return 0, nil
	}

	pageSize := len(a.pages[0])
	startPage := startIndex / pageSize
	startIndexInFirstPage := startIndex - startPage*pageSize

	pNowIndex := 0
	for i := startPage; i < len(a.pages); i++ {
		if i == startPage {
			pNowIndex += copy(p[pNowIndex:], a.pages[i][startIndexInFirstPage:])
		} else {
			pNowIndex += copy(p[pNowIndex:], a.pages[i])
		}
		if pNowIndex == len(p) {
			break
		}
	}

	return pNowIndex, nil
}

func (a *DataArea) write(p []byte) (int, error) {
	a.dataLocker.Lock()
	defer a.dataLocker.Unlock()

	if a.pages == nil {
		panic(ErrClosed)
	}

	if a.writeIndex >= a.cap {
		return 0, nil
	}

	if len(a.pages) == 0 {
		return 0, nil
	}

	pageSize := len(a.pages[0])
	startPage := a.writeIndex / pageSize
	startIndexInFirstPage := a.writeIndex - startPage*pageSize

	pNowIndex := 0
	for i := startPage; i < len(a.pages); i++ {
		n := 0
		if i == startPage {
			n = copy(a.pages[i][startIndexInFirstPage:], p[pNowIndex:])
		} else {
			n = copy(a.pages[i], p[pNowIndex:])
		}
		pNowIndex += n
		a.writeIndex += n
		if pNowIndex == len(p) {
			break
		}
	}

	return pNowIndex, nil
}

func (a *DataArea) detach() (recycle bool) {
	a.recycleLocker.Lock()
	defer a.recycleLocker.Unlock()
	a.attached--
	if a.attached <= 0 {
		recycle = true
	}
	return
}

func (a *DataArea) attach() {
	a.recycleLocker.Lock()
	defer a.recycleLocker.Unlock()
	a.attached++
}

func (a *DataArea) copyDataFrom(newA *DataArea) {
	a.dataLocker.Lock()
	defer a.dataLocker.Unlock()
	if a.pages == nil {
		panic(ErrClosed)
	}

	if len(a.pages) != len(newA.pages) {
		panic(fmt.Errorf("len of pages: %w", ErrOutOfRange))
	}

	if len(a.pages) == 0 {
		return
	}

	if len(a.pages[0]) != len(newA.pages[0]) {
		panic(fmt.Errorf("len of pages: %w", ErrOutOfRange))
	}
	for i, page := range a.pages {
		copy(newA.pages[i], page)
	}
}

func (a *DataArea) dumpData() []Page {
	a.dataLocker.Lock()
	defer a.dataLocker.Unlock()
	pages := a.pages
	a.pages = nil
	return pages
}
