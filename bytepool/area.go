package bytepool

type Page []byte

type Area struct {
	pool *Pool

	readIndex int
	detached  bool
	writable  bool
	data      *DataArea
}

func (a *Area) Read(p []byte) (int, error) {
	if a.detached {
		return 0, ErrClosed
	}
	n, err := a.data.read(p, a.readIndex)
	if err != nil {
		return n, err
	}
	a.readIndex += n
	return n, nil
}

func (a *Area) Write(p []byte) (int, error) {
	if a.detached {
		return 0, ErrClosed
	}
	return a.data.write(p)
}

func (a *Area) Detach() {
	if a.detached {
		return
	}
	a.detached = true
	if a.data.detach() {
		a.pool.Recycle(a)
	}
}

func (a *Area) ReadOnlyCopy() *Area {
	if a.detached {
		panic(ErrClosed)
	}
	a.data.attach()
	return &Area{
		pool:     a.pool,
		writable: false,
		data:     a.data,
	}
}

func (a *Area) WritableCopy() *Area {
	if a.detached {
		panic(ErrClosed)
	}
	newA := a.pool.NewArea(a.data.cap)
	newA.data.copyDataFrom(a.data)
	return newA
}
