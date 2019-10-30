package queue

import (
	"container/list"
	"sync"
)

// MutexLinkedQueue mutex-based concurrent linked list queue
type MutexLinkedQueue struct {
	l     *list.List
	mutex sync.RWMutex
}

// NewMutexLinkedQueue create new MutexLinkedQueue
func NewMutexLinkedQueue() *MutexLinkedQueue {
	return &MutexLinkedQueue{
		l: list.New(),
	}
}

// Offer insert the specified element into this queue if it is possible to do so immediately
// without violating capacity restrictions. Return nil if success.
func (queue *MutexLinkedQueue) Offer(v interface{}) {
	if v != nil {
		queue.mutex.Lock()
		queue.l.PushBack(v)
		queue.mutex.Unlock()
	}
}

// Poll retrieve and remove the head of this queue, or returns nil if this queue is empty.
func (queue *MutexLinkedQueue) Poll() (v interface{}) {
	queue.mutex.Lock()
	if e := queue.l.Front(); e != nil {
		v = e.Value
		queue.l.Remove(e)
	}
	queue.mutex.Unlock()
	return
}

// Peek retrieve, but do not remove, the head of this queue, or returns nil if this queue is empty.
func (queue *MutexLinkedQueue) Peek() (v interface{}) {
	queue.mutex.RLock()
	if e := queue.l.Front(); e != nil {
		v = e.Value
	}
	queue.mutex.RUnlock()
	return
}

// Size return the number of elements in this queue. If this queue
// contains more than math.MaxInt32 elements, returns math.MaxInt32.
func (queue *MutexLinkedQueue) Size() (size int32) {
	queue.mutex.RLock()
	size = int32(queue.l.Len())
	queue.mutex.RUnlock()
	return
}

// IsEmpty return if this queue contains no elements
func (queue *MutexLinkedQueue) IsEmpty() (empt bool) {
	return queue.Size() == 0
}

// Iterator not supported.
// MutexLinkedQueue not support iterator.
func (queue *MutexLinkedQueue) Iterator() Iterator {
	return nil
}
