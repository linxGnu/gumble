package queue

import (
	"sync/atomic"
	"unsafe"
)

type linkedListNode struct {
	_v interface{}    // real value
	_i unsafe.Pointer // wrapper over value
	_n unsafe.Pointer // next
}

func newLinkedListNode(v interface{}) *linkedListNode {
	return &linkedListNode{
		_v: v,
		_i: unsafe.Pointer(&v),
	}
}

func (n *linkedListNode) value() interface{} {
	return n._v
}

func (n *linkedListNode) next() unsafe.Pointer {
	return atomic.LoadPointer(&n._n)
}

func (n *linkedListNode) item() unsafe.Pointer {
	return atomic.LoadPointer(&n._i)
}

func (n *linkedListNode) setItemNil() {
	atomic.StorePointer(&n._i, nil)
}

func (n *linkedListNode) casItemNil(old unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&n._i, old, nil)
}

func (n *linkedListNode) casNext(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&n._n, old, new)
}

func (n *linkedListNode) setNext(new unsafe.Pointer) {
	atomic.StorePointer(&n._n, new)
}
