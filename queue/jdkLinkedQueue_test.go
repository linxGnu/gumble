package queue

import (
	"testing"
)

func TestJDKLinkedQueue_Producer(t *testing.T) {
	testProducer(t, NewQueue(JDKLinkedQueueType))
}

func TestJDKLinkedQueue_Mix(t *testing.T) {
	testMix(t, NewQueue(JDKLinkedQueueType), 50, 50)
}

type anonymousStruct struct {
	v int
}

func TestJDKLinkedQueue_Iterator(t *testing.T) {
	q := NewJDKLinkedQueue()

	as := make([]*anonymousStruct, 100)
	for i := range as {
		as[i] = &anonymousStruct{v: i}
		q.Offer(as[i])
	}

	iter := q.Iterator()
	for iter.HasNext() {
		v := iter.Next().(*anonymousStruct)
		if v.v < 50 {
			iter.Remove()
		}
	}

	iter = q.Iterator()
	for iter.HasNext() {
		v := iter.Next().(*anonymousStruct)
		if v.v < 50 {
			t.Fatal()
		}
	}
}
