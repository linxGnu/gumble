package queue

// Type queue type
type Type int

const (
	// JDKLinkedQueueType type of JDKLinkedQueue
	JDKLinkedQueueType Type = iota
	// MutexLinkedQueueType type of MutexLinkedQueue
	MutexLinkedQueueType
)

// Queue interface
type Queue interface {
	// Offer insert the specified element into this queue if it is possible to do so immediately
	// without violating capacity restrictions. Return nil if success.
	Offer(v interface{})
	// Poll retrieve and remove the head of this queue, or returns nil if this queue is empty.
	Poll() interface{}
	// Peek retrieve, but do not remove, the head of this queue, or returns nil if this queue is empty.
	Peek() interface{}
	// Size retrieve size of current queue.
	Size() int32
	// IsEmpty return if this queue contains no elements
	IsEmpty() bool
	// Iterator return an iterator over the elements in this collection.
	Iterator() Iterator
}

// Iterator interface
type Iterator interface {
	// HasNext return true if the iteration has more elements.
	HasNext() bool
	// Next return the next element in the iteration.
	Next() interface{}
	// Remove from the underlying collection the last element returned
	// by this iterator
	Remove()
}

// NewQueue create new queue based on type
func NewQueue(t Type) Queue {
	switch t {
	case MutexLinkedQueueType:
		return NewMutexLinkedQueue()
	default:
		return NewJDKLinkedQueue()
	}
}

// DefaultQueue returns jdk concurrent, non blocking queue.
func DefaultQueue() Queue {
	return NewJDKLinkedQueue()
}
