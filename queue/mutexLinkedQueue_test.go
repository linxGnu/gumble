package queue

import (
	"testing"
)

func TestMutexLinkedQueue_Producer(t *testing.T) {
	testProducer(t, NewQueue(MutexLinkedQueueType))
}

func TestMutexLinkedQueue_Mix(t *testing.T) {
	testMix(t, NewQueue(MutexLinkedQueueType), 50, 50)
}
