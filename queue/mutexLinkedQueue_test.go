package queue

import (
	"testing"
)

func TestMutexLinkedQueue_Producer(t *testing.T) {
	testProducer(t, NewQueue(MutexLinkedQueueType))
}
