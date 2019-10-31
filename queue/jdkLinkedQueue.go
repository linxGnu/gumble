package queue

import (
	"math"
	"sync/atomic"
	"unsafe"
)

// JDKLinkedQueue jdk-based concurrent non blocking linked-list queue
type JDKLinkedQueue struct {
	// The padding members 1 to 3 below are here to ensure each item is on a separate cache line.
	// This prevents false sharing and hence improves performance.
	_ [8]uint64
	h unsafe.Pointer // head
	_ [8]uint64
	t unsafe.Pointer // tail
	_ [8]uint64
}

// NewJDKLinkedQueue create new JDKLinkedQueue
func NewJDKLinkedQueue() *JDKLinkedQueue {
	q := &JDKLinkedQueue{
		t: unsafe.Pointer(&linkedListNode{}),
	}
	q.h = q.t
	return q
}

func (queue *JDKLinkedQueue) head() unsafe.Pointer {
	return atomic.LoadPointer(&queue.h)
}

func (queue *JDKLinkedQueue) tail() unsafe.Pointer {
	return atomic.LoadPointer(&queue.t)
}

// Offer inserts the specified element at the tail of this queue.
func (queue *JDKLinkedQueue) Offer(v interface{}) {
	if v != nil {
		newNode := unsafe.Pointer(newLinkedListNode(v))

		var (
			t       = queue.tail()
			p       = t
			q, oldT unsafe.Pointer
			_p      *linkedListNode
		)

		for {
			_p = (*linkedListNode)(p)
			if q = _p.next(); q == nil {
				// p is last node
				if _p.casNext(nil, newNode) {
					// Successful CAS is the linearization point
					// for e to become an element of this queue,
					// and for newNode to become "live".
					if p != t { // hop two nodes at a time
						queue.casTail(t, newNode) // Failure is OK.
					}
					return
				}
				// Lost CAS race to another thread; re-read next
			} else if p == q {
				// We have fallen off list.  If tail is unchanged, it
				// will also be off-list, in which case we need to
				// jump to head, from which all live nodes are always
				// reachable.  Else the new tail is a better bet.
				if oldT, t = t, queue.tail(); oldT != t { // t != (t = tail)?
					p = t
				} else {
					p = queue.head()
				}
			} else if p != t { // Check for tail updates after two hops.
				if oldT, t = t, queue.tail(); oldT != t {
					p = t
				} else {
					p = q
				}
			} else {
				p = q
			}
		}
	}
}

// Poll head element
func (queue *JDKLinkedQueue) Poll() (v interface{}) {
	var (
		h, p, q   unsafe.Pointer
		item, tmp unsafe.Pointer
		_p        *linkedListNode
	)

	for {
		h = queue.head()
		p = h

	loopCheck:
		for {
			_p = (*linkedListNode)(p)
			if item = _p.item(); item != nil && _p.casItemNil(item) {
				v = _p.value()
				// Successful CAS is the linearization point
				// for item to be removed from this queue.
				if p != h { // hop two nodes at a time
					if q = _p.next(); q != nil {
						tmp = q
					} else {
						tmp = p
					}
					queue.updateHead(h, tmp)
				}
				return
			} else if q = _p.next(); q == nil {
				queue.updateHead(h, p)
				return
			} else if p == q {
				break loopCheck
			} else {
				p = q
			}
		}
	}
}

// Peek return head element
func (queue *JDKLinkedQueue) Peek() (v interface{}) {
	var (
		h, p, q, item unsafe.Pointer
		_p            *linkedListNode
	)
	for {
		h = queue.head()
		p = h

	loopCheck:
		for {
			_p = (*linkedListNode)(p)
			if item = _p.item(); item != nil {
				v = _p.value()
				queue.updateHead(h, p)
				return
			} else if q = _p.next(); q == nil {
				queue.updateHead(h, p)
				return
			} else if p == q {
				break loopCheck // restart loop
			} else {
				p = q
			}
		}
	}
}

// Returns the first live (non-deleted) node on list, or null if none.
// This is yet another variant of poll/peek; here returning the
// first node, not element.  We could make peek() a wrapper around
// first(), but that would cost an extra volatile read of item,
// and the need to add a retry loop to deal with the possibility
// of losing a race to a concurrent poll().
func (queue *JDKLinkedQueue) first() unsafe.Pointer {
	var (
		h, p, q unsafe.Pointer
		hasItem bool
		_p      *linkedListNode
	)
	for {
		h = queue.head()
		p = h

	loopCheck:
		for {
			_p = (*linkedListNode)(p)
			if hasItem = _p.item() != nil; hasItem {
				queue.updateHead(h, p)
				if hasItem {
					return p
				}
				return nil
			} else if q = _p.next(); q == nil {
				queue.updateHead(h, p)
				if hasItem {
					return p
				}
				return nil
			} else if p == q {
				break loopCheck // restart loop
			} else {
				p = q
			}
		}
	}
}

// IsEmpty return if this queue contains no elements
func (queue *JDKLinkedQueue) IsEmpty() bool {
	return queue.first() == nil
}

// Size return the number of elements in this queue. If this queue
// contains more than math.MaxInt32 elements, returns math.MaxInt32.
// Beware that, unlike in most collections, this method is
// NOT a constant-time operation. Because of the
// asynchronous nature of these queues, determining the current
// number of elements requires an O(n) traversal.
// Additionally, if elements are added or removed during execution
// of this method, the returned result may be inaccurate. Thus,
// this method is typically not very useful in concurrent applications.
func (queue *JDKLinkedQueue) Size() (count int32) {
	var _p *linkedListNode
	for p := queue.first(); p != nil; p = queue.succ(p) {
		if _p = (*linkedListNode)(p); _p.item() != nil {
			if count++; count == math.MaxInt32 {
				return
			}
		}
	}
	return
}

// Remove a single instance of the specified element from this queue,
// if it is present.  More formally, removes an element e such
// that item.value is e, if this queue contains one or more such
// elements.
// Returns true if this queue contained the specified element
// (or equivalently, if this queue changed as a result of the call).
func (queue *JDKLinkedQueue) Remove(v interface{}) bool {
	if v != nil {
		var (
			next, pred, p unsafe.Pointer
			item          unsafe.Pointer
			_p            *linkedListNode
		)

		for p = queue.first(); p != nil; p = queue.succ(p) {
			_p = (*linkedListNode)(p)

			if item = _p.item(); item != nil && v == _p.value() && _p.casItemNil(item) {
				next = queue.succ(p)
				if pred != nil && next != nil {
					(*linkedListNode)(pred).casNext(p, next)
				}
				return true
			}

			pred = p
		}
	}
	return false
}

// Iterator return iterator of underlying elements.
func (queue *JDKLinkedQueue) Iterator() Iterator {
	return newJdkLinkedQueueIter(queue)
}

func (queue *JDKLinkedQueue) casTail(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&queue.t, old, new)
}

func (queue *JDKLinkedQueue) casHead(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&queue.h, old, new)
}

func (queue *JDKLinkedQueue) updateHead(h, p unsafe.Pointer) {
	if h != p && queue.casHead(h, p) {
		(*linkedListNode)(h).setNext(h)
	}
}

func (queue *JDKLinkedQueue) succ(node unsafe.Pointer) unsafe.Pointer {
	if next := (*linkedListNode)(node).next(); next != node {
		return next
	}
	return queue.head()
}

type jdkLinkedQueueIter struct {
	q        *JDKLinkedQueue
	nextNode unsafe.Pointer
	nextItem unsafe.Pointer
	lastRet  unsafe.Pointer
}

func newJdkLinkedQueueIter(q *JDKLinkedQueue) (iter *jdkLinkedQueueIter) {
	iter = &jdkLinkedQueueIter{
		q: q,
	}
	iter.advance()
	return
}

func (i *jdkLinkedQueueIter) advance() interface{} {
	i.lastRet = i.nextNode

	var (
		x                   = i.nextItem
		pred, p, next, item unsafe.Pointer
	)

	if i.nextNode == nil {
		p = i.q.first()
		pred = nil
	} else {
		pred = i.nextNode
		p = i.q.succ(i.nextNode)
	}

	for {
		if p == nil {
			i.nextNode = nil
			i.nextItem = nil
			if x != nil {
				return (*linkedListNode)(x).value()
			}
			return nil
		}

		if item = (*linkedListNode)(p).item(); item != nil {
			i.nextNode = p
			i.nextItem = item
			if x != nil {
				return (*linkedListNode)(x).value()
			}
			return nil
		}

		// skip over nils
		if next = i.q.succ(p); pred != nil && next != nil {
			(*linkedListNode)(pred).casNext(p, next)
		}

		p = next
	}
}

// HasNext return true if has next
func (i *jdkLinkedQueueIter) HasNext() bool {
	return i.nextNode != nil
}

// Next return next elements. There is no guarantee that hasNext and next are atomically due to data racy.
func (i *jdkLinkedQueueIter) Next() interface{} {
	if i.nextNode == nil {
		return nil
	}
	return i.advance()
}

// Remove from the underlying collection the last element returned
// by this iterator
func (i *jdkLinkedQueueIter) Remove() {
	l := i.lastRet
	if l == nil {
		return
	}

	// rely on a future traversal to relink.
	_l := (*linkedListNode)(l)
	_l.setItemNil()
	i.lastRet = nil
}
