package queue

import (
	"context"
	"sync"
	"testing"

	"github.com/scryner/lfreequeue"
)

var (
	maxNumberProducer = 8
	numberEle         = 10000
)

type ele struct {
	key   int
	value int
}

var preparedJDKQueue Queue
var preparedMutexQueue Queue
var preparedLFQueue Queue

type lfqueue struct {
	*lfreequeue.Queue
}

// Offer insert the specified element into this queue if it is possible to do so immediately
// without violating capacity restrictions. Return nil if success.
func (l *lfqueue) Offer(v interface{}) {
	l.Enqueue(v)
}

// Poll retrieve and remove the head of this queue, or returns nil if this queue is empty.
func (l *lfqueue) Poll() interface{} {
	v, _ := l.Dequeue()
	return v
}

// Peek retrieve, but do not remove, the head of this queue, or returns nil if this queue is empty.
func (l *lfqueue) Peek() interface{} {
	return nil
}

// Size retrieve size of current queue.
func (l *lfqueue) Size() int32 {
	return 0
}

// IsEmpty return if this queue contains no elements
func (l *lfqueue) IsEmpty() bool {
	return false
}

// Iterator return an iterator over the elements in this collection.
func (l *lfqueue) Iterator() Iterator {
	return nil
}

func init() {
	preparedJDKQueue = NewQueue(JDKLinkedQueueType)
	prepare(preparedJDKQueue)

	preparedMutexQueue = NewQueue(MutexLinkedQueueType)
	prepare(preparedMutexQueue)

	preparedLFQueue = &lfqueue{lfreequeue.NewQueue()}
	prepare(preparedLFQueue)
}

func prepare(q Queue) {
	for i := 0; i < numberEle; i++ {
		q.Offer(&ele{key: i, value: i})
	}
}

func Benchmark_MutexLinkedQueue_50P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(MutexLinkedQueueType), 50, 50)
	}
}

func Benchmark_JDKLinkedQueue_50P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(JDKLinkedQueueType), 50, 50)
	}
}

func Benchmark_LFQueue_50P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(&lfqueue{lfreequeue.NewQueue()}, 50, 50)
	}
}

func Benchmark_MutexLinkedQueue_50P10C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(MutexLinkedQueueType), 50, 10)
	}
}

func Benchmark_JDKLinkedQueue_50P10C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(JDKLinkedQueueType), 50, 10)
	}
}

func Benchmark_LFQueue_50P10C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(&lfqueue{lfreequeue.NewQueue()}, 50, 10)
	}
}

func Benchmark_MutexLinkedQueue_10P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(MutexLinkedQueueType), 10, 50)
	}
}

func Benchmark_JDKLinkedQueue_10P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(NewQueue(JDKLinkedQueueType), 10, 50)
	}
}

func Benchmark_LFQueue_10P50C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchQueueMix(&lfqueue{lfreequeue.NewQueue()}, 10, 50)
	}
}

func benchQueueMix(q Queue, numberProducer, numberConsumer int) {
	for i := 0; i < numberProducer; i++ {
		go func(i int) {
			for j := 0; j < numberEle; j++ {
				q.Offer(&ele{key: i, value: -j})
			}
		}(i)
	}

	ch := make(chan *ele, 100000)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < numberConsumer; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
					if item := q.Poll(); item != nil {
						ch <- item.(*ele)
					}
				}
			}
		}(ctx)
	}

	for i := 0; i < numberProducer*numberEle; i++ {
		<-ch
	}
	cancel()
	wg.Wait()
}

func Benchmark_MutexLinkedQueue_100P(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(NewQueue(MutexLinkedQueueType), 100)
	}
}
func Benchmark_JDKLinkedQueue_100P(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(NewQueue(JDKLinkedQueueType), 100)
	}
}

func Benchmark_LFQueue_100P(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(&lfqueue{lfreequeue.NewQueue()}, 100)
	}
}

func Benchmark_MutexLinkedQueue_100C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(preparedMutexQueue, 100)
	}
}
func Benchmark_JDKLinkedQueue_100C(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchProducer(preparedJDKQueue, 100)
	}
}

func benchProducer(q Queue, numberProducer int) {
	var wg sync.WaitGroup
	for i := 0; i < numberProducer; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < numberEle; j++ {
				q.Offer(&ele{value: -j})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
