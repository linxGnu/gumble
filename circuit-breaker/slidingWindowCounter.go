package cbreaker

import (
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"

	ga "github.com/linxGnu/gumble/adder"
	queue "github.com/linxGnu/gumble/queue"
)

// bucket hold the count of events within {@code updateInterval}.
type bucket struct {
	timestamp int64
	s         ga.LongAdder // number of success
	f         ga.LongAdder // number of failure
}

// newBucket create new bucket
func newBucket(timestamp int64) *bucket {
	return &bucket{
		timestamp: timestamp,
		s:         ga.NewLongAdder(ga.JDKAdderType),
		f:         ga.NewLongAdder(ga.JDKAdderType),
	}
}

func (b *bucket) add(succ bool) {
	if succ {
		b.s.Add(1)
	} else {
		b.f.Add(1)
	}
}

// Success return number of success operation
func (b *bucket) success() int64 {
	return b.s.Sum()
}

// Failure return number of failure operation
func (b *bucket) failure() int64 {
	return b.f.Sum()
}

// SlidingWindowCounter that accumulates the count of events within a time window.
type SlidingWindowCounter struct {
	ticker              Ticker
	cur                 *bucket
	slidingWindowNanos  int64
	updateIntervalNanos int64
	snapshot            atomic.Value
	reservoir           queue.Queue
}

// NewSlidingWindowCounter create new SlidingWindowCounter
func NewSlidingWindowCounter(ticker Ticker, slidingWindowNanos, updateIntervalNanos time.Duration) (s *SlidingWindowCounter, e error) {
	if ticker == nil {
		e = fmt.Errorf("Ticker is required")
		return
	}

	s = &SlidingWindowCounter{
		ticker:              ticker,
		slidingWindowNanos:  int64(slidingWindowNanos),
		updateIntervalNanos: int64(updateIntervalNanos),
		cur:                 newBucket(ticker.Tick()),
		reservoir:           queue.DefaultQueue(),
	}
	s.snapshot.Store(EventCountZero)
	return
}

func (s *SlidingWindowCounter) current() *bucket {
	return (*bucket)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.cur))))
}

func (s *SlidingWindowCounter) casCurrent(old, new *bucket) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&s.cur)),
		unsafe.Pointer(old),
		unsafe.Pointer(new),
	)
}

// Count return the current EventCount
func (s *SlidingWindowCounter) Count() *EventCount {
	return s.snapshot.Load().(*EventCount)
}

// OnSuccess count success events.
func (s *SlidingWindowCounter) OnSuccess() *EventCount {
	return s.onEvent(true)
}

// OnFailure count failure events.
func (s *SlidingWindowCounter) OnFailure() *EventCount {
	return s.onEvent(false)
}

func (s *SlidingWindowCounter) onEvent(succ bool) (e *EventCount) {
	tickerNanos, currentBucket := s.ticker.Tick(), s.current()

	if tickerNanos < currentBucket.timestamp {
		// if current timestamp is older than bucket's timestamp (maybe race or GC pause?),
		// then creates an instant bucket and puts it to the reservoir not to lose event.
		bucket := newBucket(tickerNanos)
		bucket.add(succ)
		s.reservoir.Offer(bucket)
		return
	}

	if tickerNanos < currentBucket.timestamp+s.updateIntervalNanos {
		currentBucket.add(succ)
		return
	}

	nextBucket := newBucket(tickerNanos)
	nextBucket.add(succ)

	// replaces the bucket
	if s.casCurrent(currentBucket, nextBucket) {
		// puts old one to the reservoir
		s.reservoir.Offer(currentBucket)

		// and then updates count
		e = s.trimAndSum(tickerNanos)
		s.snapshot.Store(e)
	} else {
		// the bucket has been replaced already
		// puts new one as an instant bucket to the reservoir not to lose event
		s.reservoir.Offer(nextBucket)
	}
	return
}

func (s *SlidingWindowCounter) trimAndSum(t int64) *EventCount {
	oldLimit, iterator := t-s.slidingWindowNanos, s.reservoir.Iterator()

	var nxt interface{}
	var bck *bucket
	var success, failure int64

	for iterator.HasNext() {
		if nxt = iterator.Next(); nxt != nil {
			if bck = nxt.(*bucket); bck.timestamp < oldLimit {
				// removes old bucket
				iterator.Remove()
			} else {
				success += bck.success()
				failure += bck.failure()
			}
		}
	}

	return NewEventCount(success, failure)
}
