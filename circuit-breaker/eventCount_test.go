package cbreaker

import "testing"

func TestEventCount(t *testing.T) {
	ev := NewEventCount(0, 0)
	if ev.SuccessRate() != -1 {
		t.Errorf("Fail to catch trivial case of success rate")
	}
	if ev.FailureRate() != -1 {
		t.Errorf("Fail to catch trivial case of failure rate")
	}

	ev = NewEventCount(5, 20)
	if ev.Success() != 5 || ev.success != 5 || ev.Failure() != 20 || ev.failure != 20 || ev.Total() != 25 {
		t.Errorf("Fail to create new EventCount")
	}

	if ev.SuccessRate() != 0.2 || ev.FailureRate() != 0.8 {
		t.Errorf("Fail to return rate")
	}
}
