package cbreaker

import (
	"time"
)

// CircuitBreakerBuilder instance using builder pattern.
type CircuitBreakerBuilder struct {
	name                    *Name
	ticker                  Ticker
	failureRateThreshold    float64
	minimumRequestThreshold int64
	trialRequestInterval    time.Duration
	circuitOpenWindow       time.Duration
	counterSlidingWindow    time.Duration
	counterUpdateInterval   time.Duration
	listeners               CircuitBreakerListeners
}

// NewCircuitBreakerBuilder create new circuit breaker builder
func NewCircuitBreakerBuilder() (c *CircuitBreakerBuilder) {
	c = &CircuitBreakerBuilder{
		ticker:                  SystemTicker,
		failureRateThreshold:    defaultFailureRateThreshold,
		minimumRequestThreshold: defaultMinimumRequestThreshold,
		trialRequestInterval:    defaultTrialRequestInterval,
		circuitOpenWindow:       defaultCircuitOpenWindow,
		counterSlidingWindow:    defaultCounterSlidingWindow,
		counterUpdateInterval:   defaultCounterUpdateInterval,
	}
	return
}

// Name set name for circuit breaker
func (c *CircuitBreakerBuilder) Name(name *Name) *CircuitBreakerBuilder {
	c.name = name
	return c
}

// SetTicker set ticker
func (c *CircuitBreakerBuilder) SetTicker(t Ticker) *CircuitBreakerBuilder {
	c.ticker = t
	return c
}

// SetFailureRateThreshold set the threshold of failure rate to detect a remote service fault.
func (c *CircuitBreakerBuilder) SetFailureRateThreshold(failureRateThreshold float64) *CircuitBreakerBuilder {
	c.failureRateThreshold = failureRateThreshold
	return c
}

// SetMinimumRequestThreshold set the minimum number of requests within a time window necessary to detect a remote service fault.
func (c *CircuitBreakerBuilder) SetMinimumRequestThreshold(minimumRequestThreshold int64) *CircuitBreakerBuilder {
	c.minimumRequestThreshold = minimumRequestThreshold
	return c
}

// SetTrialRequestInterval set the trial request interval in HalfOpen state.
func (c *CircuitBreakerBuilder) SetTrialRequestInterval(trialRequestInterval time.Duration) *CircuitBreakerBuilder {
	c.trialRequestInterval = trialRequestInterval
	return c
}

// SetCircuitOpenWindow set the duration of Open state.
func (c *CircuitBreakerBuilder) SetCircuitOpenWindow(circuitOpenWindow time.Duration) *CircuitBreakerBuilder {
	c.circuitOpenWindow = circuitOpenWindow
	return c
}

// SetCounterSlidingWindow set the time length of sliding window to accumulate the count of events.
func (c *CircuitBreakerBuilder) SetCounterSlidingWindow(counterSlidingWindow time.Duration) *CircuitBreakerBuilder {
	c.counterSlidingWindow = counterSlidingWindow
	return c
}

// SetCounterUpdateInterval set the interval that a circuit breaker can see the latest accumulated count of events.
func (c *CircuitBreakerBuilder) SetCounterUpdateInterval(counterUpdateInterval time.Duration) *CircuitBreakerBuilder {
	c.counterUpdateInterval = counterUpdateInterval
	return c
}

// AddListener add a CircuitBreakerListener
func (c *CircuitBreakerBuilder) AddListener(listener CircuitBreakerListener) *CircuitBreakerBuilder {
	if listener != nil {
		c.listeners = append(c.listeners, listener)
	}
	return c
}

// Build return a newly-created CircuitBreaker based on the properties of this builder.
func (c *CircuitBreakerBuilder) Build() (cb CircuitBreaker, err error) {
	cb, err = NewNonBlockingCircuitBreaker(c.ticker, &CircuitBreakerConfig{
		name:                    c.name,
		failureRateThreshold:    c.failureRateThreshold,
		minimumRequestThreshold: c.minimumRequestThreshold,
		trialRequestInterval:    c.trialRequestInterval,
		circuitOpenWindow:       c.circuitOpenWindow,
		counterSlidingWindow:    c.counterSlidingWindow,
		counterUpdateInterval:   c.counterUpdateInterval,
		listeners:               c.listeners,
	})
	return
}
