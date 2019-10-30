package cbreaker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"
)

// NonBlockingCircuitBreaker a non-blocking implementation of circuit breaker pattern.
type NonBlockingCircuitBreaker struct {
	name   *Name
	ticker Ticker
	config *CircuitBreakerConfig
	s      *nonBlockingCircuitBreakerState // state
}

// NewNonBlockingCircuitBreaker create new NonBlockingCircuitBreaker
func NewNonBlockingCircuitBreaker(ticker Ticker, config *CircuitBreakerConfig) (nbc *NonBlockingCircuitBreaker, err error) {
	if ticker == nil || config == nil {
		err = fmt.Errorf("Ticker and Config is required")
		return
	}

	// validate config
	if err = config.Validate(); err != nil {
		return
	}

	nbc = &NonBlockingCircuitBreaker{
		name:   config.name,
		config: config,
		ticker: ticker,
	}
	nbc.s = nbc.newClosedState()

	nbc.logStateTransition(CircuitStateClosed, nil)
	nbc.notifyStateChanged(CircuitStateClosed)

	return
}

func (nb *NonBlockingCircuitBreaker) state() *nonBlockingCircuitBreakerState {
	return (*nonBlockingCircuitBreakerState)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nb.s))))
}

func (nb *NonBlockingCircuitBreaker) casState(old, new *nonBlockingCircuitBreakerState) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&nb.s)),
		unsafe.Pointer(old),
		unsafe.Pointer(new),
	)
}

// Name return the name of the circuit breaker.
func (nb *NonBlockingCircuitBreaker) Name() *Name {
	return nb.name
}

// Execute delegated function and report a success or a failure to the specified CircuitBreaker according to the completed value
func (nb *NonBlockingCircuitBreaker) Execute(ctx context.Context, delegatedFn Execute) (r interface{}, err error) {
	if delegatedFn == nil {
		return
	}

	if nb.CanRequest() {
		r, err = delegatedFn(ctx)
	} else {
		err = ErrFailFast
	}

	return
}

// OnSuccess report a remote invocation success.
func (nb *NonBlockingCircuitBreaker) OnSuccess() {
	currentState := nb.state()
	if currentState.isClosed() {
		// fires success event
		if updatedCount := currentState.counter.OnSuccess(); updatedCount != nil {
			// notifies the count if it has been updated
			nb.notifyCountUpdated(updatedCount)
		}
	} else if currentState.isHalfOpen() {
		// changes to CLOSED if at least one request succeeds during HALF_OPEN
		if nb.casState(currentState, nb.newClosedState()) {
			nb.logStateTransition(CircuitStateClosed, nil)
			nb.notifyStateChanged(CircuitStateClosed)
		}
	}
}

// OnFailure report a remote invocation failure.
func (nb *NonBlockingCircuitBreaker) OnFailure() {
	currentState := nb.state()
	if currentState.isClosed() {
		// fires failure event
		if updatedCount := currentState.counter.OnFailure(); updatedCount != nil {
			if nb.checkIfExceedingFailureThreshold(updatedCount) &&
				nb.casState(currentState, nb.newOpenState()) {
				nb.logStateTransition(CircuitStateOpen, updatedCount)
				nb.notifyStateChanged(CircuitStateOpen)
			} else {
				nb.notifyCountUpdated(updatedCount)
			}
		}
	} else if currentState.isHalfOpen() {
		// returns to OPEN if a request fails during HALF_OPEN
		if nb.casState(currentState, nb.newOpenState()) {
			nb.logStateTransition(CircuitStateOpen, nil)
			nb.notifyStateChanged(CircuitStateOpen)
		}
	}
}

// CanRequest decide whether a request should be sent or failed depending on the current circuit state.
func (nb *NonBlockingCircuitBreaker) CanRequest() bool {
	currentState := nb.state()
	if currentState.isClosed() {
		// all requests are allowed during CLOSED
		return true
	}

	if currentState.isHalfOpen() || currentState.isOpen() {
		if currentState.checkTimeout() &&
			nb.casState(currentState, nb.newHalfOpenState()) {
			nb.logStateTransition(CircuitStateHalfOpen, nil)
			nb.notifyStateChanged(CircuitStateHalfOpen)
			return true
		}

		// all other requests are refused
		nb.notifyRequestRejected()
		return false
	}

	return true
}

func (nb *NonBlockingCircuitBreaker) checkIfExceedingFailureThreshold(count *EventCount) bool {
	total := count.Total()
	return 0 < total && nb.config.minimumRequestThreshold <= total && nb.config.failureRateThreshold < count.FailureRate()
}

func (nb *NonBlockingCircuitBreaker) newOpenState() *nonBlockingCircuitBreakerState {
	return newNonBlockingCircuitBreakerState(nb.ticker, CircuitStateOpen, nb.config.circuitOpenWindow, noOpCounter)
}

func (nb *NonBlockingCircuitBreaker) newHalfOpenState() *nonBlockingCircuitBreakerState {
	return newNonBlockingCircuitBreakerState(nb.ticker, CircuitStateHalfOpen, nb.config.trialRequestInterval, noOpCounter)
}

func (nb *NonBlockingCircuitBreaker) newClosedState() *nonBlockingCircuitBreakerState {
	slidingWindow, _ := NewSlidingWindowCounter(nb.ticker, nb.config.counterSlidingWindow, nb.config.counterUpdateInterval)
	return newNonBlockingCircuitBreakerState(nb.ticker, CircuitStateClosed, 0, slidingWindow)
}

func (nb *NonBlockingCircuitBreaker) notifyStateChanged(circuitState CircuitState) {
	for _, listener := range nb.config.listeners {
		if listener != nil {
			if err := listener.OnStateChanged(nb, circuitState); err != nil {
				if logger != nil {
					logger.Warn("An error occurred when notifying a StateChanged event", err)
				}
			}
			if err := listener.OnEventCountUpdated(nb, EventCountZero); err != nil {
				if logger != nil {
					logger.Warn("An error occurred when notifying an EventCountUpdated event", err)
				}
			}
		}
	}
}

func (nb *NonBlockingCircuitBreaker) notifyCountUpdated(count *EventCount) {
	for _, listener := range nb.config.listeners {
		if listener != nil {
			if err := listener.OnEventCountUpdated(nb, count); err != nil {
				if logger != nil {
					logger.Warn("An error occurred when notifying an EventCountUpdated event", err)
				}
			}
		}
	}
}

func (nb *NonBlockingCircuitBreaker) notifyRequestRejected() {
	for _, listener := range nb.config.listeners {
		if listener != nil {
			if err := listener.OnRequestRejected(nb); err != nil {
				if logger != nil {
					logger.Warn("An error occurred when notifying a RequestRejected event", err)
				}
			}
		}
	}
}

func (nb *NonBlockingCircuitBreaker) logStateTransition(state CircuitState, count *EventCount) {
	if logger != nil {
		var name string
		if nb.name != nil {
			name = nb.name.Namespace + "_" + nb.name.Subsystem + "_" + nb.name.Name
		}

		var _state string
		if state == CircuitStateOpen {
			_state = "OPEN"
		} else if state == CircuitStateHalfOpen {
			_state = "HALF_OPEN"
		} else {
			_state = "CLOSED"
		}

		var fl string
		if count != nil {
			fl = fmt.Sprintf(" fail:%d total:%d", count.Failure(), count.Total())
		}
		logger.Info("name:" + name + " state:" + _state + fl)
	}
}

// nonBlockingCircuitBreakerState state inside non blocking circuit breaker
type nonBlockingCircuitBreakerState struct {
	cs                CircuitState
	counter           EventCounter
	timeout           int64
	timedOutTimeNanos time.Duration
	ticker            Ticker
}

func newNonBlockingCircuitBreakerState(ticker Ticker, cs CircuitState, timedOutTimeNanos time.Duration, counter EventCounter) *nonBlockingCircuitBreakerState {
	return &nonBlockingCircuitBreakerState{
		cs:                cs,
		counter:           counter,
		timedOutTimeNanos: timedOutTimeNanos,
		timeout:           ticker.Tick() + timedOutTimeNanos.Nanoseconds(),
		ticker:            ticker,
	}
}

func (ns *nonBlockingCircuitBreakerState) checkTimeout() bool {
	return ns.timedOutTimeNanos > 0 && ns.timeout <= ns.ticker.Tick()
}

func (ns *nonBlockingCircuitBreakerState) isOpen() bool {
	return ns.cs == CircuitStateOpen
}

func (ns *nonBlockingCircuitBreakerState) isHalfOpen() bool {
	return ns.cs == CircuitStateHalfOpen
}

func (ns *nonBlockingCircuitBreakerState) isClosed() bool {
	return ns.cs == CircuitStateClosed
}

type noopCounter struct{}

func (noop *noopCounter) Count() *EventCount {
	return EventCountZero
}

func (noop *noopCounter) OnSuccess() *EventCount {
	return nil
}

func (noop *noopCounter) OnFailure() *EventCount {
	return nil
}

var noOpCounter = &noopCounter{}
