// Package cbreaker (circuit-breaker) contains circuit breakers ported from line/armeria.
package cbreaker

import (
	"context"
	"fmt"
)

var (
	// ErrFailFast indicate that fail fast detected
	ErrFailFast = fmt.Errorf("FailFast detected")
)

// Name fully-qualified name
type Name struct {
	Namespace string
	Subsystem string
	Name      string
}

// Execute function
type Execute func(ctx context.Context) (result interface{}, err error)

// CircuitBreaker circuit breaker, which tracks the number of
// success/failure requests and detects a remote service failure.
type CircuitBreaker interface {
	// Name return the name of the circuit breaker.
	Name() *Name
	// OnSuccess report a remote invocation success.
	OnSuccess()
	// OnFailure report a remote invocation failure.
	OnFailure()
	// CanRequest decide whether a request should be sent or failed depending on the current circuit state.
	CanRequest() bool
	// Execute delegated function
	Execute(ctx context.Context, delegatedFn Execute) (r interface{}, err error)
}

// CircuitBreakerListener listener interface for receiving events
type CircuitBreakerListener interface {
	// OnStateChanged invoked when the circuit state is changed.
	OnStateChanged(cb CircuitBreaker, state CircuitState) (err error)
	// OnEventCountUpdated invoked when the circuit breaker's internal EventCount is updated.
	OnEventCountUpdated(cb CircuitBreaker, eventCount *EventCount) (err error)
	// OnRequestRejected invoked when the circuit breaker rejects a request.
	OnRequestRejected(cb CircuitBreaker) (err error)
	// Stop notify listener to stop
	Stop()
}

// CircuitBreakerListeners collection of CircuitBreakerListener
type CircuitBreakerListeners []CircuitBreakerListener
