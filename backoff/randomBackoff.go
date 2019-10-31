package backoff

import (
	"fmt"
)

// RandomBackoff compute backoff delay which is a random value between
// minDelayMillis} and maxDelayMillis.
type RandomBackoff struct {
	minDelayMillis int64
	maxDelayMillis int64
	bound          int64
}

// NewRandomBackoff create new
func NewRandomBackoff(minDelayMillis, maxDelayMillis int64) (b *RandomBackoff, err error) {
	if minDelayMillis < 0 {
		err = fmt.Errorf("minDelayMillis: %d (expected: >= 0)", minDelayMillis)
	} else if minDelayMillis > maxDelayMillis {
		err = fmt.Errorf("maxDelayMillis: %d (expected: >= %d)", maxDelayMillis, minDelayMillis)
	} else {
		b = &RandomBackoff{minDelayMillis: minDelayMillis, maxDelayMillis: maxDelayMillis, bound: maxDelayMillis - minDelayMillis}
	}
	return
}

// NextDelayMillis return the number of milliseconds to wait for before attempting a retry.
func (f *RandomBackoff) NextDelayMillis(numAttemptsSoFar int) int64 {
	if f.minDelayMillis != f.maxDelayMillis {
		return nextRandomInt64(f.bound) + f.minDelayMillis
	}
	return f.minDelayMillis
}
