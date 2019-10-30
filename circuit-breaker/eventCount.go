package cbreaker

// EventCounter event counter interface
type EventCounter interface {
	// Count return the current EventCount.
	Count() *EventCount
	// OnSuccess count success events
	OnSuccess() *EventCount
	// OnFailure count failure events
	OnFailure() *EventCount
}

// EventCount stores the count of events.
type EventCount struct {
	success int64
	failure int64
}

// NewEventCount create new event count
func NewEventCount(success, failure int64) *EventCount {
	return &EventCount{success: success, failure: failure}
}

// Success return number of success events
func (e *EventCount) Success() int64 {
	return e.success
}

// Failure return number of failure events
func (e *EventCount) Failure() int64 {
	return e.failure
}

// Total return total number of events.
func (e *EventCount) Total() int64 {
	return e.success + e.failure
}

// SuccessRate return number of success rate
func (e *EventCount) SuccessRate() float64 {
	total := e.Total()
	if total == 0 {
		return -1
	}
	return float64(e.success) / float64(total)
}

// FailureRate return number of failure rate
func (e *EventCount) FailureRate() float64 {
	total := e.Total()
	if total == 0 {
		return -1
	}
	return float64(e.failure) / float64(total)
}

// EventCountZero event count with zero
var EventCountZero = NewEventCount(0, 0)
