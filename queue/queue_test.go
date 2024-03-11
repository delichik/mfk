package queue

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkDequeueQueueSingle(b *testing.B) {
	q := New[int](0)
	for i := 0; i < b.N; i++ {
		_ = q.Enqueue(1)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = q.Dequeue()
	}
}

func TestQueue(t *testing.T) {
	if runtime.NumCPU() < 3 {
		t.Error("Then count of CPU cores is less than 3")
		t.Fail()
	}
	q := New[int](100)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	total := int64(0)
	dropped := int64(0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runtime.LockOSThread()
		for ctx.Err() == nil {
			if err := q.TryEnqueue(1); err != nil {
				atomic.AddInt64(&dropped, 1)
			}
			atomic.AddInt64(&total, 1)
		}
	}()

	runtime.LockOSThread()
	for r := 0; r < 10; r++ {
		for i := 0; i < 1000000; i++ {
			_, _ = q.Dequeue()
		}
	}
	cancel()
	wg.Wait()
	t.Logf("Total: %d, Dropped: %d, DropRate: %f", total, dropped, float64(dropped)/float64(total))
}
