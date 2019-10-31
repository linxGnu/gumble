package backoff

import (
	"fmt"
	"math"
)

// ExponentialBackoff wait an exponentially-increasing amount of time between attempts.
type ExponentialBackoff struct {
	initialDelayMillis int64
	maxDelayMillis     int64
	multiplier         float64
}

// NewExponentialBackoff create new ExponentialBackoff
func NewExponentialBackoff(initialDelayMillis, maxDelayMillis int64, multiplier float64) (b *ExponentialBackoff, err error) {
	if multiplier <= 1 {
		err = fmt.Errorf("multiplier: %.3f (expected: > 1.0)", multiplier)
	} else if initialDelayMillis < 0 {
		err = fmt.Errorf("initialDelayMillis: %d (expected: >= 0)", initialDelayMillis)
	} else if initialDelayMillis > maxDelayMillis {
		err = fmt.Errorf("maxDelayMillis: %d (expected: >= %d)", maxDelayMillis, initialDelayMillis)
	} else {
		b = &ExponentialBackoff{
			initialDelayMillis: initialDelayMillis,
			maxDelayMillis:     maxDelayMillis,
			multiplier:         multiplier,
		}
	}
	return
}

// NextDelayMillis return the number of milliseconds to wait for before attempting a retry.
func (f *ExponentialBackoff) NextDelayMillis(numAttemptsSoFar int) (nextDelay int64) {
	if numAttemptsSoFar == 1 {
		return f.initialDelayMillis
	}

	nextDelay = saturatedMultiply(f.initialDelayMillis, math.Pow(f.multiplier, float64(numAttemptsSoFar-1)))
	if nextDelay > f.maxDelayMillis {
		nextDelay = f.maxDelayMillis
	}
	return
}
