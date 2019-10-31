package backoff

import (
	"fmt"
)

// FixedBackoff wait a fixed delay between attempts.
type FixedBackoff struct {
	delayMillis int64
}

// NewFixedBackoff create new fixed backoff
func NewFixedBackoff(delayMillis int64) (b *FixedBackoff, err error) {
	if delayMillis >= 0 {
		b = &FixedBackoff{delayMillis: delayMillis}
	} else {
		err = fmt.Errorf("delayMillis: %d (expected: >= 0)", delayMillis)
	}
	return
}

// NextDelayMillis return the number of milliseconds to wait for before attempting a retry.
func (f *FixedBackoff) NextDelayMillis(numAttemptsSoFar int) int64 {
	return f.delayMillis
}

// NoDelayBackoff return a Backoff that will never wait between attempts.
// In most cases, using Backoff without delay is very dangerous.
var NoDelayBackoff Backoff = &FixedBackoff{delayMillis: 0}

// NoRetry return a Backoff indicates that no retry
var NoRetry Backoff = &FixedBackoff{delayMillis: -1}
