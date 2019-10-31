package backoff

import "fmt"

const (
	// DefaultDelayMillis default delay millis
	DefaultDelayMillis int64 = 200
	// DefaultInitialDelayMillis default initial delay millis
	DefaultInitialDelayMillis int64 = 200
	// DefaultMinDelayMillis default min delay millis
	DefaultMinDelayMillis int64 = 0
	// DefaultMaxDelayMillis default max delay millis
	DefaultMaxDelayMillis int64 = 10000
	// DefaultMultiplier default multiplier
	DefaultMultiplier float64 = 2.0
	// DefaultMinJitterRate default min jitter rate
	DefaultMinJitterRate float64 = -0.2
	// DefaultMaxJitterRate default max jitter rate
	DefaultMaxJitterRate float64 = 0.2
)

var (
	// ErrInvalidSpecFormat invalid specification format
	ErrInvalidSpecFormat = fmt.Errorf("Invalid format of specification")
)
