package cbreaker

import (
	"fmt"
	"time"
)

const (
	defaultFailureRateThreshold    = 0.8
	defaultMinimumRequestThreshold = 10
	defaultTrialRequestInterval    = time.Duration(3 * time.Second)
	defaultCircuitOpenWindow       = time.Duration(10 * time.Second)
	defaultCounterSlidingWindow    = time.Duration(20 * time.Second)
	defaultCounterUpdateInterval   = time.Duration(1 * time.Second)
)

// CircuitState states of circuit breaker
type CircuitState byte

const (
	// CircuitStateClosed initial state. All requests are sent to the remote service.
	CircuitStateClosed CircuitState = 0
	// CircuitStateOpen the circuit is tripped. All requests fail immediately without calling the remote service.
	CircuitStateOpen CircuitState = 1
	// CircuitStateHalfOpen only one trial request is sent at a time until at least one request succeeds or fails.
	// If it doesn't complete within a certain time, another trial request will be sent again.
	// All other requests fails immediately same as OPEN.
	CircuitStateHalfOpen CircuitState = 2
)

var (
	// ErrTickerDurationInvalid ticker duration invalid
	ErrTickerDurationInvalid = fmt.Errorf("Ticker duration must be > 0")
)
