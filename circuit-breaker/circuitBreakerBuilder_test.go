package cbreaker

import (
	"testing"
	"time"
)

func TestNewCircuitBreakerBuilder(t *testing.T) {
	if builder := NewCircuitBreakerBuilder(); builder == nil {
		t.Errorf("Fail to create new default circuit breaker")
	}
}

func TestBuilderSetFailureRate(t *testing.T) {
	builder := NewCircuitBreakerBuilder()
	builder.Name(&Name{Name: "dummy-breaker"})

	// try to set FailureRateThreshold
	validFailureRateThreshold := 0.5
	if _, err := builder.SetFailureRateThreshold(-1).Build(); err == nil {
		t.Errorf("Fail to set FailureRateThreshold")
	} else if _, err = builder.SetFailureRateThreshold(0).Build(); err == nil {
		t.Errorf("Fail to set FailureRateThreshold")
	} else if _, err = builder.SetFailureRateThreshold(1.1).Build(); err == nil {
		t.Errorf("Fail to set FailureRateThreshold")
	} else if _, err = builder.SetFailureRateThreshold(validFailureRateThreshold).Build(); err != nil || builder.failureRateThreshold != validFailureRateThreshold {
		t.Errorf("Fail to set FailureRateThreshold")
	}
}

func TestBuilderSetMinimumRequestThreshold(t *testing.T) {
	builder := NewCircuitBreakerBuilder()

	// try to set MinimumRequestThreshold
	var validMinimumRequestThreshold int64 = 123
	builder.SetMinimumRequestThreshold(validMinimumRequestThreshold)
	if builder.minimumRequestThreshold != validMinimumRequestThreshold {
		t.Errorf("Fail to set MinimumRequestThreshold")
	}
}

func TestBuilderSetTrialRequestInterval(t *testing.T) {
	builder := NewCircuitBreakerBuilder()

	// try to set TrialRequestInterval
	validTrialRequestInterval := time.Duration(23 * time.Second)
	if _, err := builder.SetTrialRequestInterval(-1).Build(); err == nil {
		t.Errorf("Fail to set TrialRequestInterval")
	} else if _, err = builder.SetTrialRequestInterval(0).Build(); err == nil {
		t.Errorf("Fail to set TrialRequestInterval")
	} else if _, err = builder.SetTrialRequestInterval(validTrialRequestInterval).Build(); err != nil || builder.trialRequestInterval != validTrialRequestInterval {
		t.Errorf("Fail to set TrialRequestInterval")
	}
}

func TestBuilderSetCircuitOpenWindow(t *testing.T) {
	builder := NewCircuitBreakerBuilder()

	// try to set CircuitOpenWindow
	validCircuitOpenWindow := time.Duration(10 * time.Millisecond)
	if _, err := builder.SetCircuitOpenWindow(-1).Build(); err == nil {
		t.Errorf("Fail to set CircuitOpenWindow")
	} else if _, err = builder.SetCircuitOpenWindow(0).Build(); err == nil {
		t.Errorf("Fail to set CircuitOpenWindow")
	} else if _, err = builder.SetCircuitOpenWindow(validCircuitOpenWindow).Build(); err != nil || builder.circuitOpenWindow != validCircuitOpenWindow {
		t.Errorf("Fail to set CircuitOpenWindow")
	}
}

func TestBuilderSetCounterSlidingWindow(t *testing.T) {
	builder := NewCircuitBreakerBuilder()

	// try to set CounterSlidingWindow
	validCounterSlidingWindow := time.Duration(11 * time.Second)
	if _, err := builder.SetCounterSlidingWindow(-1).Build(); err == nil {
		t.Errorf("Fail to set CounterSlidingWindow")
	} else if _, err = builder.SetCounterSlidingWindow(0).Build(); err == nil {
		t.Errorf("Fail to set CounterSlidingWindow")
	} else if _, err = builder.SetCounterSlidingWindow(validCounterSlidingWindow).Build(); err != nil || builder.counterSlidingWindow != validCounterSlidingWindow {
		t.Errorf("Fail to set CounterSlidingWindow")
	}
}

func TestBuilderSetCounterUpdateInterval(t *testing.T) {
	builder := NewCircuitBreakerBuilder()

	// try to set CounterUpdateInterval
	validCounterUpdateInterval := time.Duration(9 * time.Second)
	if _, err := builder.SetCounterUpdateInterval(-1).Build(); err == nil {
		t.Errorf("Fail to set CounterUpdateInterval")
	} else if _, err = builder.SetCounterUpdateInterval(0).Build(); err == nil {
		t.Errorf("Fail to set CounterUpdateInterval")
	} else if _, err = builder.SetCounterUpdateInterval(validCounterUpdateInterval).Build(); err != nil || builder.counterUpdateInterval != validCounterUpdateInterval {
		t.Errorf("Fail to set CounterUpdateInterval")
	}
}

func TestBuilderAddListener(t *testing.T) {
	builder := NewCircuitBreakerBuilder()

	// try to add 4 listeners
	builder.AddListener(nil)
	builder.AddListener(&dummyCircuitBreakerListener{})
	builder.AddListener(&dummyCircuitBreakerListener{})
	builder.AddListener(&dummyCircuitBreakerListener{})
	builder.AddListener(&dummyCircuitBreakerListener{})
	builder.AddListener(nil)
}

type dummyCircuitBreakerListener struct{}

// OnStateChanged invoked when the circuit state is changed.
func (d *dummyCircuitBreakerListener) OnStateChanged(cb CircuitBreaker, state CircuitState) (err error) {
	return
}

// OnEventCountUpdated invoked when the circuit breaker's internal EventCount is updated.
func (d *dummyCircuitBreakerListener) OnEventCountUpdated(cb CircuitBreaker, eventCount *EventCount) (err error) {
	return
}

// OnRequestRejected invoked when the circuit breaker rejects a request.
func (d *dummyCircuitBreakerListener) OnRequestRejected(cb CircuitBreaker) (err error) {
	return
}

// Stop listening
func (d *dummyCircuitBreakerListener) Stop() {
}
